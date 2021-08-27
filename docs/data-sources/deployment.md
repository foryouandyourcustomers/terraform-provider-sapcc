---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sapcc_deployment Data Source - terraform-provider-sapcc"
subcategory: ""
description: |-
  Fetches the  Commerce Cloud deployment details for the provided deployment code. More information on the configuration parameters at getDeployment api https://help.sap.com/viewer/452dcbb0e00f47e88a69cdaeb87a925d/v1905/en-US/d86d3539bd284410bc83817297a117ac.html
---

# sapcc_deployment (Data Source)

Fetches the  Commerce Cloud deployment details for the provided deployment `code`. More information on the configuration parameters at [getDeployment api](https://help.sap.com/viewer/452dcbb0e00f47e88a69cdaeb87a925d/v1905/en-US/d86d3539bd284410bc83817297a117ac.html)

## Example Usage

```terraform
data "sapcc_deployment" "deployment" {
  # the deployment code
  code = "1230182"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **code** (String) The deployment code assigned to this deployment.

### Optional

- **build_code** (String) The build code associated with this deployment.
- **cancelation** (Object) If the deployment was cancelled, the cancellation details. (see [below for nested schema](#nestedatt--cancelation))
- **created_by** (String) The User Id of the user who created this build.
- **created_timestamp** (String) Build start timestamp in UTC.
- **database_update_mode** (String) The database update mode for the deployment.
- **deployed_timestamp** (String) Deploy start timestamp in UTC.
- **environment_code** (String) The environment code of the environment of the deployment.
- **failed_timestamp** (String) If the deployment fails, the failed timestamp.
- **scheduled_timestamp** (String) Timestamp when the deployment was initially scheduled.
- **status** (String) Status of the Deployment.
- **strategy** (String) The strategy used for this deployment.
- **subscription_code** (String) The subscription id associated to the build.
- **undeployed_timestamp** (String) If the deployment was rolledback, the rollback timestamp.

<a id="nestedatt--cancelation"></a>
### Nested Schema for `cancelation`

Optional:

- **canceled_by** (String)
- **failed** (Boolean)
- **finished_timestamp** (String)
- **rollback_database** (Boolean)
- **start_timestamp** (String)

