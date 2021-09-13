package provider

import (
	"context"
	"fmt"
	"terraform-provider-sapcc/internal/client"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceDeploymentType struct{}

func (r dataSourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the  Commerce Cloud deployment details for the provided deployment `code`. More information on the configuration parameters at [getDeployment api](https://help.sap.com/viewer/452dcbb0e00f47e88a69cdaeb87a925d/v1905/en-US/d86d3539bd284410bc83817297a117ac.html)",
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
				Computed:    false,
				Required:    true,
			},
			"build_code": {
				Description: "The build code associated with this deployment.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"strategy": {
				Description: "The strategy used for this deployment.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"environment_code": {
				Description: "The environment code of the environment of the deployment.",
				Type:        types.StringType,
				Computed:    false,
				Optional:    true,
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
				Computed:    true,
				Optional:    true,
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
						Optional:    true,
					},
					"start_timestamp": {
						Description: "The deployment start of the timestamp.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
					"finished_timestamp": {
						Description: "The deployment finished timestamp.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
					"failed": {
						Description: "If the deployment failed.",
						Type:        types.BoolType,
						Computed:    true,
						Optional:    true,
					},
					"rollback_database": {
						Description: "Id the database was rollback.",
						Type:        types.BoolType,
						Computed:    true,
						Optional:    true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (r dataSourceDeploymentType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceDeployment{
		provider: *(p.(*provider)),
	}, nil
}

type dataSourceDeployment struct {
	provider provider
}

func (ds dataSourceDeployment) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !ds.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)

		return
	}

	// Declare struct that this function will set to this data source's state
	var deploymentRequest models.Deployment

	resp.Diagnostics.Append(req.Config.Get(ctx, &deploymentRequest)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deployment := fetchDeployment(deploymentRequest.Code.Value, ds.provider.client, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &deployment)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func fetchDeployment(deployCode string, client *client.Client, diags *diag.Diagnostics) *models.Deployment {
	deployResponse, st, err := client.GetDeployment(deployCode)

	if err != nil {
		diags.AddError(fmt.Sprintf("Error fetching deployment %s", err), "")
	}

	handleDeploymentDiags(deployCode, st, diags)

	if !diags.HasError() {
		progress, st, err := client.GetDeploymentProgress(deployCode)

		if err != nil {
			diags.AddError(fmt.Sprintf("Error fetching deployment progress %s", err), "")
			return deployResponse

		}

		handleDeploymentDiags(deployCode, st, diags)

		if !diags.HasError() {
			deployResponse.ProgressPercentage = progress.ProgressPercentage
			deployResponse.Status = progress.DeployStatus
		}
	}

	return deployResponse
}

func handleDeploymentDiags(deployCode string, st int, diags *diag.Diagnostics) {
	switch st {
	case 400:
		diags.AddError(
			fmt.Sprintf("Bad Request got 400 for code '%s'; Is the environment busy with another deployment in-progress?", deployCode),
			"",
		)

	case 404:
		diags.AddError(
			fmt.Sprintf("Deployment or progress not found; code '%s'; Check logs or report it provider developer", deployCode),
			"",
		)

	case 401:
		diags.AddError(
			fmt.Sprintf("Unauthorized, credentials invalid for code '%s', please verify your 'auth_token' and 'subscription_id'", deployCode),
			"",
		)

	case 403:
		diags.AddError(
			fmt.Sprintf("Forbidden, can not access deployment with code '%s'", deployCode),
			"",
		)

	case 200:
		return

	case 201:
		return

	default:
		diags.AddError(
			fmt.Sprintf("Unexpected http status %d for code '%s' from upstream api; won't continue. expected 200 ", st, deployCode),
			"",
		)
	}
}
