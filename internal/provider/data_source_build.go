package provider

import (
	"context"
	"fmt"
	"terraform-provider-sapcc/internal/client"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceBuildType struct{}

func (r dataSourceBuildType) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
	return tfsdk.Schema{
		Description: "Fetches the  Commerce Cloud build details for the provided build `code`. More information on the configuration parameters at [getBuild api](https://help.sap.com/viewer/452dcbb0e00f47e88a69cdaeb87a925d/v1905/en-US/9041daaf93c144acb4726f0c86e58337.html)",
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

	var buildRequest models.Build
	// TODO: try using GetAttribute instead?
	for _, d := range req.Config.Get(ctx, &buildRequest) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	buildCode := buildRequest.Code.Value
	httpClient, err := client.NewClient(r.p.SubscriptionBaseURL, r.p.AuthToken)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating http client",
		})

		return
	}

	buildResponse, st, err := httpClient.GetBuild(buildCode)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  fmt.Sprintf("Error fetching build %s", err),
		})

		return
	}

	if buildResponse == nil {
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
	}

	// Set state
	for _, d := range resp.State.Set(ctx, &buildResponse) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}
}
