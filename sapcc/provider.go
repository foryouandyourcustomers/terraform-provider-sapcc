package sapcc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
)

var stderr = os.Stderr

const defaultApiBaseUrl = "https://portalrotapi.hana.ondemand.com/v2/subscriptions"

func New() tfsdk.Provider {
	return &provider{}
}

// provider describes the data is passed along the context and is available to the resources
type provider struct {
	configured          bool
	SubscriptionBaseUrl string
	AuthToken           string
}

// Provider schema struct
type providerData struct {
	ApiBaseUrl     types.String `tfsdk:"api_baseurl"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	AuthToken      types.String `tfsdk:"auth_token"`
}

// GetSchema returns the schema for the provider
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"auth_token": {
				Description: "The authentication token for accessing SAP Commerce Cloud API. Can be specified with the `SAPCC_AUTH_TOKEN` environment variable.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Computed:    false,
			},
			"subscription_id": {
				Description: "The subscription Id associated with the SAP Commerce Cloud. Can be specified with the `SAPCC_SUBSCRIPTION_ID` environment variable.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    false,
			},
			"api_baseurl": {
				Description: "The base api url Subscription url for SAP Commerce Cloud API. Default: `https://portalrotapi.hana.ondemand.com/v2/subscriptions",
				Type:        types.StringType,
				Optional:    true,
				Computed:    false,
			},
		},
	}, nil
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData

	for _, d := range req.Config.Get(ctx, &config) {
		resp.Diagnostics = append(resp.Diagnostics, d)
		return
	}

	var apiBaseUrl string
	if config.ApiBaseUrl.Unknown || config.ApiBaseUrl.Null || config.ApiBaseUrl.Value == "" {
		apiBaseUrl = defaultApiBaseUrl
	} else {
		apiBaseUrl = config.ApiBaseUrl.Value
	}

	var subscriptionId string
	if config.SubscriptionId.Unknown {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Can not create the provider.",
			Detail:   "Cannot use unknown value as for 'subscription_id'",
		})
	}

	if config.SubscriptionId.Null {
		subscriptionId = os.Getenv("SAPCC_SUBSCRIPTION_ID")
	} else {
		subscriptionId = config.SubscriptionId.Value
	}

	if subscriptionId == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Can not create the provider.",
			Detail:   "Cannot use empty value for 'subscription_id'",
		})
	}

	var authToken string
	if config.AuthToken.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Can not create the provider.",
			Detail:   "Cannot use unknown value as 'auth_token'",
		})
	}

	if config.AuthToken.Null {
		authToken = os.Getenv("SAPCC_AUTH_TOKEN")
	} else {
		authToken = config.AuthToken.Value
	}

	if authToken == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Can not create the provider.",
			Detail:   "Cannot use empty value for 'auth_token'",
		})
	}

	p.SubscriptionBaseUrl = fmt.Sprintf("%s/%s", apiBaseUrl, subscriptionId)
	p.AuthToken = authToken
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{
		"sapcc_deployment": resourceDeploymentType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{
		"sapcc_build":      dataSourceBuildType{},
		"sapcc_deployment": dataSourceDeploymentType{},
	}, nil
}
