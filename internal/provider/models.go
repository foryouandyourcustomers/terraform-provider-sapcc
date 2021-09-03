package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Build is the model that describes a build
type Build struct {
	// CreatedBy
	CreatedBy                    types.String `tfsdk:"created_by"`
	SubscriptionCode             types.String `tfsdk:"subscription_code"`
	ApplicationCode              types.String `tfsdk:"application_code"`
	ApplicationDefinitionVersion types.String `tfsdk:"application_definition_version"`
	Branch                       types.String `tfsdk:"branch"`
	Name                         types.String `tfsdk:"name"`
	// Code represents the build code
	Code                types.String    `tfsdk:"code"`
	BuildStartTimestamp types.String    `tfsdk:"build_start_timestamp"`
	BuildEndTimestamp   types.String    `tfsdk:"build_end_timestamp"`
	BuildVersion        types.String    `tfsdk:"build_version"`
	Status              types.String    `tfsdk:"status"`
	Properties          []BuildProperty `tfsdk:"properties"`
}

// BuildProperty describing additional property configured for the build
type BuildProperty struct {
	Key types.String `tfsdk:"key"`
	Val types.String `tfsdk:"val"`
}

// Deployment is the model that describes the deployment
type Deployment struct {
	CreatedBy           types.String       `tfsdk:"created_by"`
	SubscriptionCode    types.String       `tfsdk:"subscription_code"`
	CreatedTimestamp    types.String       `tfsdk:"created_timestamp"`
	BuildCode           types.String       `tfsdk:"build_code"`
	EnvironmentCode     types.String       `tfsdk:"environment_code"`
	DatabaseUpdateMode  types.String       `tfsdk:"database_update_mode"`
	Code                types.String       `tfsdk:"code"`
	Strategy            types.String       `tfsdk:"strategy"`
	ScheduledTimestamp  types.String       `tfsdk:"scheduled_timestamp"`
	DeployedTimestamp   types.String       `tfsdk:"deployed_timestamp"`
	FailedTimestamp     types.String       `tfsdk:"failed_timestamp"`
	UndeployedTimestamp types.String       `tfsdk:"undeployed_timestamp"`
	Status              types.String       `tfsdk:"status"`
	Cancelation         DeployCancellation `tfsdk:"cancelation"`
}

type DeployCancellation struct {
	CancelledBy      types.String `tfsdk:"canceled_by"`
	StartTimestamp   types.String `tfsdk:"start_timestamp"`
	FinishTimestamp  types.String `tfsdk:"finished_timestamp"`
	Failed           types.Bool   `tfsdk:"failed"`
	RollbackDatabase types.Bool   `tfsdk:"rollback_database"`
}
