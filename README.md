*Important note: This repo is not currently being developed. Feel free to fork and develop it*
----
![experimental](https://img.shields.io/badge/status-alpha-important) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/foryouandyourcustomers/terraform-provider-sapcc) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/foryouandyourcustomers/terraform-provider-sapcc/release?label=release-action) ![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/foryouandyourcustomers/terraform-provider-sapcc?include_prereleases&label=release&sort=semver) [![GitHub license](https://img.shields.io/github/license/foryouandyourcustomers/terraform-provider-sapcc)](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/blob/master/LICENSE)  
# Terraform Provider for SAP Commerce Cloud API

An **experimental** terraform provider for interacting with SAP Commerce Cloud. The provider uses the [new plugin framework](https://github.com/hashicorp/terraform-plugin-framework) and is considered to be in the _early stages_ of development. Although, we try to keep the provider as stable as possible, the provider is **not** ready to be used in the production because the upstream framework is still in `alpha` and there _will_(!) be breaking changes.

For any breaking changes, we expect them to be released as `pre-releases` (`*beta.*` tag) first and then finally released. 


## Using the provider
Simple & straight forward usage of the provider. Detailed examples can be found under (`examples/`) or in [*_test.go](./internal/provider/)
```terraform
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.2"
      source  = "foryouandyourcustomers/sapcc"
    }
  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = var.base_url        # optional
  auth_token      = var.auth_token      # or export SAPCC_AUTH_TOKEN=..
  subscription_id = var.subscription_id # or export SAPCC_SUBSCRIPTION_ID=..
}

data "sapcc_build" "build" {
  # the build code
  code = "20210819.0"
}

resource "sapcc_deployment" "deployment" {
  # the build code to be used for deployment
  build_code = data.sapcc_build.build.code
  # the environment code to be used for deployment
  environment_code = "d0"
  # the strategy used for deployment
  strategy = "ROLLING_UPDATE"
  # the database update mode for deployment
  database_update_mode = "NONE"
}
```


## Building the provider

### Structure of repo
- Resources (`internal/provider`),
- Examples (`examples/`) and generated documentation (`docs/`)
- WireMock's responses from commerce cloud apis are in `mocks`

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0.3
- [Go](https://golang.org/doc/install) >= 1.17
- [Golangci-lint](https://golangci-lint.run/usage/install)

1. Clone the repository
1. Enter the repository directory
1. Install provider
```sh
$ make install
```
this installs the provider under `~/.terraform.d/plugins/registry.terraform.io/foryouandyourcustomers/sapcc/<version>/<arch>`

## Testing the provider
At the moment, there's no official testing framework (see [#issue 113](https://github.com/hashicorp/terraform-plugin-framework/issues/113)). The [helper](./helper) library temporarily fills the gap of acceptance tests. Once there is an official test suite for writing tests (unit and acceptance tests), this helper **will be** deprecated.  

The repo comes with a mock server for emulating the SAP Commerce Cloud behaviour. The intension is to be use this server as a way to develop the provider instead of directly interacting with SAP Commerce cloud itself. For this case, the [mock responses](./mocks/mappings) have been designed carefully based on the official API documentation and responses.

To run the acceptance tests, you need to either provide terraform cli path with `TF_ACC_TERRAFORM_EXEC_PATH` or provide the version of terraform to run the tests against with `TF_ACC_TERRAFORM_VERSION`, this is mostly useful in case of running with CI.

```shell
TF_ACC_TERRAFORM_EXEC_PATH=/path/to/terraform make testacc   
```



## TODOs
- [ ] Add unit tests
- [x] Improve testing
- [x] Integrate acceptance testing in the source
- [ ] Cleanup code
- [x] Create a dedicated the `http` client

## Roadmap

- [X] Fetch build details
- [X] Fetch deployment Details
- [X] Create deployment
- [X] Add deployment progress during deployment
- [X] Update deployment with the change build code
- [ ] Create new build
- [ ] Add build logs during build
- [ ] Cancel build?
- [ ] Cancel deployment?
- [ ] Create and Delete Customer properties
