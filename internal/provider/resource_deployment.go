package provider

import (
	"context"
	"fmt"
	"terraform-provider-sapcc/internal/client"
	"terraform-provider-sapcc/internal/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// track the deployment progress every xxx mins
const trackProgressTimeSecs = 30

type resourceDeploymentType struct{}

// GetSchema resourceDeploymentType Resource schema
func (r resourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Creates & triggers a deployment for SAP Commerce Cloud. More information on the configuration parameters at [createDeployment api](https://help.sap.com/viewer/452dcbb0e00f47e88a69cdaeb87a925d/v1905/en-US/d80fd1dbefff4b8bbbbac66822d4a038.html)",
		Attributes: map[string]tfsdk.Attribute{
			"created_by": {
				Description: "The User Id of the user who created this build.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"subscription_code": {
				Description: "The subscription id associated to the build.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"code": {
				Description: "The deployment code assigned to this deployment.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"build_code": {
				Description: "The build code associated with this deployment.",
				Type:        types.StringType,
				Required:    true,
				Computed:    false,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplaceIf(func(ctx context.Context, state, config attr.Value, path *tftypes.AttributePath) (bool, diag.Diagnostics) {
						if state != nil && config != nil {
							stateVal := state.(types.String)
							configVal := config.(types.String)

							if !stateVal.Unknown && !stateVal.Null && !configVal.Unknown && !configVal.Null {
								if configVal.Value != stateVal.Value {
									return true, nil
								}
							}
						}
						return false, nil
					}, "Changes in build code will trigger new deployment", ""),
				},
			},
			"strategy": {
				Description: "The strategy used for this deployment.",
				Type:        types.StringType,
				Required:    true,
				Computed:    false,
			},
			"environment_code": {
				Description: "The environment code of the environment of the deployment.",
				Type:        types.StringType,
				Required:    true,
				Computed:    false,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplaceIf(func(ctx context.Context, state, config attr.Value, path *tftypes.AttributePath) (bool, diag.Diagnostics) {
						if state != nil && config != nil {
							stateVal := state.(types.String)
							configVal := config.(types.String)

							if !stateVal.Unknown && !stateVal.Null && !configVal.Unknown && !configVal.Null {
								if configVal.Value != stateVal.Value {
									return true, nil
								}
							}
						}
						return false, nil
					}, "Changes in environment code will trigger new deployment", ""),
				},
			},
			"created_timestamp": {
				Description: "Build start timestamp in UTC.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"deployed_timestamp": {
				Description: "Deploy start timestamp in UTC.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"database_update_mode": {
				Description: "The database update mode for the deployment.",
				Type:        types.StringType,
				Required:    true,
				Computed:    false,
			},
			"scheduled_timestamp": {
				Description: "Timestamp when the deployment was initially scheduled.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"failed_timestamp": {
				Description: "If the deployment fails, the failed timestamp.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"undeployed_timestamp": {
				Description: "If the deployment was rolledback, the rollback timestamp.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"status": {
				Description: "Status of the Deployment.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"deploy_progress_percentage": {
				Description: "Overall deployment progress percentage.",
				Type:        types.NumberType,
				Computed:    true,
				Optional:    true,
			},
			"cancelation": {
				Description: "If the deployment was cancelled, the cancellation details.",
				Computed:    true,
				Optional:    true,
				//FIXME: This is a possible bug in the framework:
				// we expect here schema.SingleNestedAttributes but we use a List as workaround
				// https://github.com/hashicorp/terraform-plugin-framework/issues/112
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"canceled_by": {
						Description: "The UserId of the user who cancelled the deployment.",
						Type:        types.StringType,
						Computed:    true,
						Required:    true,
					},
					"start_timestamp": {
						Description: "The deployment start of the timestamp.",
						Type:        types.StringType,
						Computed:    true,
						Required:    true,
					},
					"finished_timestamp": {
						Description: "The deployment finished timestamp.",
						Type:        types.StringType,
						Computed:    true,
						Required:    true,
					},
					"failed": {
						Description: "If the deployment failed.",
						Type:        types.BoolType,
						Computed:    true,
						Required:    true,
					},
					"rollback_database": {
						Description: "Id the database was rollback.",
						Type:        types.BoolType,
						Computed:    true,
						Required:    true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

// New resource instance
func (r resourceDeploymentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDeployment{
		provider: *(p.(*provider)),
	}, nil
}

type resourceDeployment struct {
	provider provider
}

// Create a new resource
func (rs resourceDeployment) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !rs.provider.configured {
		resp.Diagnostics.Append(
			diag.NewErrorDiagnostic(
				"Provider not configured",
				"The provider hasn't been configured before apply,"+
					" likely because it depends on an unknown value from another resource."+
					" This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
			))

		return
	}

	// Retrieve values from plan
	var plan models.Deployment
	for _, d := range req.Config.Get(ctx, &plan) {
		resp.Diagnostics.Append(d)
		return
	}

	err, deployResponse := createNewDeployment(rs.provider.client, &resp.Diagnostics, &plan)

	if err {
		return
	}

	for _, d := range resp.State.Set(ctx, deployResponse) {
		resp.Diagnostics.Append(d)
		return
	}
}

// Read resource information
func (rs resourceDeployment) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from plan
	var state models.Deployment
	for _, d := range req.State.Get(ctx, &state) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	deployCode := state.Code

	if deployCode.Unknown || deployCode.Null {
		// this means the resource hasn't yet to be created - silently return

	} else {

		err, state := fetchDeployment(state.Code.Value, rs.provider.client, &resp.Diagnostics)

		if err {
			return
		}

		progress, status, pErr := rs.provider.client.GetDeploymentProgress(deployCode.Value)

		if pErr != nil {
			resp.Diagnostics.Append(
				diag.NewErrorDiagnostic(fmt.Sprintf("Error fetching deployment progress %s", pErr),
					"",
				))

			return
		}

		dErr := handleDeploymentDiags(state.BuildCode.Value, status, &resp.Diagnostics)

		if dErr {
			return
		}

		state.ProgressPercentage = progress.ProgressPercentage
		state.Status = progress.DeployStatus

		for _, d := range resp.State.Set(ctx, &state) {
			resp.Diagnostics.Append(d)
			return
		}
	}
}

// Update resource
func (rs resourceDeployment) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// In general, there isn't anything to update on the already existing deployment
	// we will not change in the existing state, a change that will trigger deployments is already handled in "Schema" definition
	// TODO: should we re-create/update deployment on minor updates?
	var plan, state models.Deployment

	for _, d := range req.Plan.Get(ctx, &plan) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	for _, d := range req.State.Get(ctx, &state) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}
	// Setting the planned updates to state! the drifted attributes will be updated during a refresh
	for _, d := range resp.State.Set(ctx, plan) {
		resp.Diagnostics.Append(d)
		return
	}
}

// Delete resource
func (rs resourceDeployment) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// There isn't anything to clean up on the server side, the old deployment is already created and doesn't need any cleanup
	// In-future, there might be some cleanups here
	resp.State.RemoveResource(ctx)
}

func createNewDeployment(client *client.Client, diags *diag.Diagnostics, plan *models.Deployment) (bool, *models.Deployment) {
	deployResponse, st, err := client.CreateDeployment(plan)

	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(fmt.Sprintf("Error creating deployment %s", err),
				"",
			))

		return true, deployResponse
	}

	er := handleDeploymentDiags(plan.BuildCode.Value, st, diags)

	if er {
		return true, deployResponse
	}

	deployStatus := deployResponse.Status.Value
	deployCode := deployResponse.Code.Value

	for {
		if deployStatus == "DEPLOYED" || deployStatus == "FAIL" {
			break
		}

		time.Sleep(trackProgressTimeSecs * time.Second)

		progress, status, err := client.GetDeploymentProgress(deployCode)

		if err != nil {
			diags.Append(
				diag.NewErrorDiagnostic(fmt.Sprintf("Error fetching deployment progress %s", err),
					"",
				))

			return true, deployResponse
		}

		er = handleDeploymentDiags(plan.BuildCode.Value, status, diags)

		if er {
			return true, deployResponse
		}

		deployResponse.ProgressPercentage = progress.ProgressPercentage
		deployStatus = progress.DeployStatus.Value

		logger.Info("Deploying buildcode#", hclog.Fmt("%s %s (%f)", deployCode, progress.DeployStatus.Value, progress.ProgressPercentage.Value))
	}

	if deployStatus != "DEPLOYED" {
		diags.Append(
			diag.NewErrorDiagnostic(fmt.Sprintf("Buiild wasn't successfully deployed; status is %s", deployStatus),
				"",
			))

		return true, deployResponse
	}

	// deployment is now created
	// let's fill in the rest of the data
	er, deployment := fetchDeployment(deployResponse.Code.Value, client, diags)

	if er {
		return true, deployResponse
	}

	return false, deployment
}
