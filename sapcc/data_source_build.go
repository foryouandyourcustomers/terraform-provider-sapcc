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

type dataSourceBuildType struct{}

func (r dataSourceBuildType) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
	return tfsdk.Schema{
		Description: "Fetches the current build details for the provided `code`",
		Attributes: map[string]tfsdk.Attribute{
			"created_by": {
				Description: "The S-user Id of the user who created this build.",
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
			"application_code": {
				Description: "The application code for the build.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"application_definition_version": {
				Description: "The application definition version.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"branch": {
				Description: "The name of the source branch used for creating the build.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "The name of the build used when it was created.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"code": {
				Description: "The build code for this build.",
				Type:        types.StringType,
				Computed:    false,
				Required:    true,
			},
			"build_start_timestamp": {
				Description: "The timestamp when the build was started.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"build_end_timestamp": {
				Description: "The timestamp when the build was ended.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"build_version": {
				Description: "The full build version.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"status": {
				Description: "The final status of this build.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"properties": {
				Description: "List of properties associated with this build.",
				Computed:    true,
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						Description: "Key of the property.",
						Type:        types.StringType,
						Computed:    true,
						Required:    true,
					},
					"val": {
						Description: "Value associated with this property.",
						Type:        types.StringType,
						Computed:    true,
						Required:    true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (r dataSourceBuildType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceBuild{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceBuild struct {
	p provider
}

func (r dataSourceBuild) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !r.p.configured {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider not configured",
			Detail:   "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		})
		return
	}

	// Declare struct that this function will set to this data source's state
	var build Build

	for _, d := range req.Config.Get(ctx, &build) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/builds/%s", r.p.SubscriptionBaseUrl, build.Code.Value)
	authToken := r.p.AuthToken
	buildCode := build.Code.Value
	fmt.Fprintf(stderr, "[DEBUG] %s url : %s\n", buildCode, url)

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
			Summary:  fmt.Sprintf("Build '%s' not found", buildCode),
		})
		return
	case 401:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Unauthorized, credentials invalid for build '%s', please verify your 'auth_token' and 'subscription_id' ", buildCode),
		})
		return
	case 403:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Forbidden, can not access build '%s'", buildCode),
		})
		return
	case 200:
		break
	default:
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Unexpected http status %d for build '%s' from upstream api; won't continue. expected 200 ", st, buildCode),
		})
		return
	}

	buildResponse := make(map[string]interface{}, 0)
	err = json.NewDecoder(res.Body).Decode(&buildResponse)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Error decoding build response %s", err),
		})
		return
	}

	for k, v := range buildResponse {
		switch k {
		case "createdBy":
			build.CreatedBy = types.String{Value: v.(string)}
		case "branch":
			build.Branch = types.String{Value: v.(string)}
		case "applicationDefinitionVersion":
			build.ApplicationDefinitionVersion = types.String{Value: v.(string)}
		case "subscriptionCode":
			build.SubscriptionCode = types.String{Value: v.(string)}
		case "applicationCode":
			build.ApplicationCode = types.String{Value: v.(string)}
		case "code":
			build.Code = types.String{Value: v.(string)}
		case "buildVersion":
			build.BuildVersion = types.String{Value: v.(string)}
		case "name":
			build.Name = types.String{Value: v.(string)}
		case "buildStartTimestamp":
			build.BuildStartTimestamp = types.String{Value: v.(string)}
		case "buildEndTimestamp":
			build.BuildEndTimestamp = types.String{Value: v.(string)}
		case "status":
			build.Status = types.String{Value: v.(string)}
		case "properties":
			var properties []BuildProperty
			for _, v := range v.([]interface{}) {
				v := v.(map[string]interface{})
				properties = append(properties, BuildProperty{
					Key: types.String{Value: v["key"].(string)},
					Val: types.String{Value: v["value"].(string)},
				})
			}
			build.Properties = properties
		default:
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityWarning,
				Summary:  fmt.Sprintf("Unexpected data recevied from build response k=%s v=%s, ignoring", k, v),
			})
		}
	}

	fmt.Fprintf(stderr, "\n[DEBUG]-Resource State Build:%+v", build)

	// Set state
	for _, d := range resp.State.Set(ctx, &build) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}
}
