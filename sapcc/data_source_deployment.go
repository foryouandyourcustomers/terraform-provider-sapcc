package sapcc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"net/http"
	"time"
)

type dataSourceDeploymentType struct{}

func (r dataSourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
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
				Description: "Status of the Deployment.",
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

func (r dataSourceDeploymentType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceDeployment{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceDeployment struct {
	p provider
}

func (r dataSourceDeployment) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !r.p.configured {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider not configured",
			Detail:   "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		})
		return
	}

	// Declare struct that this function will set to this data source's state
	var deployment Deployment

	for _, d := range req.Config.Get(ctx, &deployment) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	fmt.Fprintf(stderr, "[DEBUG] deployment %s\n", deployment.Code.Value)

	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/deployments/%s", r.p.SubscriptionBaseUrl, deployment.Code.Value)
	authToken := r.p.AuthToken
	deploymentCode := deployment.Code.Value
	fmt.Fprintf(stderr, "[DEBUG] %s deployment url : %s\n", deploymentCode, url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating http client",
		})
		return
	}
	request.Header = http.Header{
		"Authorization": []string{authToken},
		"Content-Type":  []string{"application/json"},
	}
	res, err := client.Do(request)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Error retrieving build %s", err),
		})
		return
	}
	defer res.Body.Close()
	st := res.StatusCode
	switch st {
	case 404:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Build '%s' not found", deploymentCode),
		})
		return
	case 401:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Unauthorized, credentials invalid for build '%s', please verify your 'auth_token' and 'subscription_id' ", deploymentCode),
		})
		return
	case 403:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Forbidden, can not access build '%s'", deploymentCode),
		})
		return
	case 200:
		break
	default:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Unexpected http status %d for build '%s' from upstream api; won't continue. expected 200 ", st, deploymentCode),
		})
		return
	}

	deploymentResponse := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&deploymentResponse)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Error decoding build response %s", err),
		})
		return
	}

	for k, v := range deploymentResponse {
		var val string
		// we get many nil elements
		if v != nil {
			stringVal, ok := v.(string)
			if ok {
				val = stringVal
			} else {
				fmt.Fprintf(stderr, "\n[DEBUG] Nonstring type received key:%s value:%s", k, v)
			}
		}

		switch k {
		case "createdBy":
			deployment.CreatedBy = types.String{Value: val}
		case "buildCode":
			deployment.BuildCode = types.String{Value: val}
		case "createdTimestamp":
			deployment.CreatedTimestamp = types.String{Value: val}
		case "subscriptionCode":
			deployment.SubscriptionCode = types.String{Value: val}
		case "environmentCode":
			deployment.EnvironmentCode = types.String{Value: val}
		case "code":
			deployment.Code = types.String{Value: val}
		case "databaseUpdateMode":
			deployment.DatabaseUpdateMode = types.String{Value: val}
		case "strategy":
			deployment.Strategy = types.String{Value: val}
		case "scheduledTimestamp":
			deployment.ScheduledTimestamp = types.String{Value: val}
		case "deployedTimestamp":
			deployment.DeployedTimestamp = types.String{Value: val}
		case "undeployedTimestamp":
			deployment.UndeployedTimestamp = types.String{Value: val}
		case "failedTimestamp":
			deployment.FailedTimestamp = types.String{Value: val}
		case "status":
			deployment.Status = types.String{Value: val}
		case "cancelation":
			fmt.Fprintf(stderr, "\n[DEBUG]-cancelation:%s", v)
			var cancelation []DeployCancellation
			if v != nil {
				v := v.(map[string]interface{})
				cancelation = append(cancelation, DeployCancellation{
					CancelledBy:      types.String{Value: v["canceledBy"].(string)},
					StartTimestamp:   types.String{Value: v["startTimestamp"].(string)},
					FinishTimestamp:  types.String{Value: v["finishedTimestamp"].(string)},
					Failed:           types.Bool{Value: v["failed"].(bool)},
					RollbackDatabase: types.Bool{Value: v["rollbackDatabase"].(bool)},
				})
			}
			deployment.Cancelation = cancelation
		default:
			fmt.Fprintf(stderr, "\n[DEBUG] dataSourceDeployment %s Unhandled key:%s value:%s, ignoring", deploymentCode, k, v)
		}

	}
	fmt.Fprintf(stderr, "\n[DEBUG]-Resource State deployment:%+v", deployment)

	// Set state
	for _, d := range resp.State.Set(ctx, &deployment) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}
}
