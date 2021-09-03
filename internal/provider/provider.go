package provider

import (
	"context"
	"fmt"
	"os"
	"terraform-provider-sapcc/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var stderr = os.Stderr

const defaultAPIBaseURL = "https://portalrotapi.hana.ondemand.com/v2/subscriptions"

func New() tfsdk.Provider {
	return &provider{}
}

// provider describes the data is passed along the context and is available to the resources
type provider struct {
	configured bool
	client     *client.Client
}

// Provider schema struct
type providerData struct {
	APIBaseURL     types.String `tfsdk:"api_baseurl"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
	AuthToken      types.String `tfsdk:"auth_token"`
}

// GetSchema returns the schema for the provider
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, []*tfprotov6.Diagnostic) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"auth_token": {
				Description: "The API token that grant you access to Commerce Cloud APIs. Can be specified with the `SAPCC_AUTH_TOKEN` environment variable.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Computed:    false,
			},
			"subscription_id": {
				Description: "The Subscription Id associated with the SAP Commerce Cloud. Can be specified with the `SAPCC_SUBSCRIPTION_ID` environment variable.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    false,
			},
			"api_baseurl": {
				Description: "The base url for SAP Commerce Cloud API. Default: `https://portalrotapi.hana.ondemand.com/v2/subscriptions`",
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

	var apiBaseURL string
	if config.APIBaseURL.Unknown || config.APIBaseURL.Null || config.APIBaseURL.Value == "" {
		apiBaseURL = defaultAPIBaseURL
	} else {
		apiBaseURL = config.APIBaseURL.Value
	}

	var subscriptionID string

	if config.SubscriptionID.Unknown {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Can not create the provider.",
			Detail:   "Cannot use unknown value as for 'subscription_id'",
		})
	}

	if config.SubscriptionID.Null {
		subscriptionID = os.Getenv("SAPCC_SUBSCRIPTION_ID")
	} else {
		subscriptionID = config.SubscriptionID.Value
	}

	if subscriptionID == "" {
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

	httpClient, err := client.NewClient(fmt.Sprintf("%s/%s", apiBaseURL, subscriptionID), authToken)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating http client",
		})

		return
	}
	p.client = httpClient
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
