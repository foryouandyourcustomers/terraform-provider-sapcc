package provider

import (
	"context"
	"fmt"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type resourceDeploymentType struct{}

// resourceDeploymentType Resource schema
func (r resourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
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
				Description: "Status of the models.Deployment.",
				Type:        types.StringType,
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
func (r resourceDeploymentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
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
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider not configured",
			Detail:   "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		})

		return
	}

	// Retrieve values from plan
	var plan models.Deployment
	for _, d := range req.Config.Get(ctx, &plan) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	deployResponse, st, err := rs.provider.client.CreateDeployment(&plan)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Error creating deployment %s", err),
		})

		return
	}

	switch st {
	case 404:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Build '%s' not found", plan.BuildCode.Value),
		})

		return
	case 401:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unauthorized, credentials invalid for build, please verify your 'auth_token' and 'subscription_id'",
		})

		return
	case 403:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Forbidden, can not access build",
		})

		return
	case 200:
		if deployResponse == nil {
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "No deployment code received from api, check logs for errors. ",
			})
			return
		}

	default:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Unexpected http status %d  from upstream api; won't continue. expected 200 ", st),
		})

		return
	}

	for _, d := range resp.State.Set(ctx, deployResponse) {
		resp.Diagnostics = append(resp.Diagnostics, d)
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

	if state.Code.Unknown || state.Code.Null {
		// this means the resource hasn't yet to be created - silently return

	} else {
		err, diags, state := fetchDeployment(state.Code.Value, rs.provider.client, resp.Diagnostics)

		if err {
			resp.Diagnostics = diags
			return
		}

		for _, d := range resp.State.Set(ctx, &state) {
			resp.Diagnostics = append(diags, d)
			return
		}
	}
}

// Update resource
func (rs resourceDeployment) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
		Severity: tfprotov6.DiagnosticSeverityWarning,
		Summary:  "Not implemented",
		Detail:   "Update of the deployment is not supported yet.",
	})
}

// Delete resource
func (rs resourceDeployment) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
		Severity: tfprotov6.DiagnosticSeverityWarning,
		Summary:  "Not implemented",
		Detail:   "Deleting/Rollback of the deployment is not supported yet.",
	})
}
