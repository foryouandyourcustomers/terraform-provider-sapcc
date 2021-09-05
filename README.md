# Terraform Provider for SAP Commerce Cloud API

Am experimental terraform provider for experimenting with SAP Commerce Cloud API. The provider uses the [new plugin framework](https://github.com/hashicorp/terraform-plugin-framework) and is considered to be in the early stages of development. The provider is not ready to be used in the production as the upstream framework is still in `alpha` and there _will_ be breaking changes.

- A resource, and a data source (`internal/provider`),
- Examples (`examples/`) and generated documentation (`docs/`)
- WireMock's responses from commerce cloud apis are in `mocks`

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0.3
- [Go](https://golang.org/doc/install) >= 1.16
- [Golangci-lint](https://golangci-lint.run/usage/install)

## Installing The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider
```sh
$ make build
```

You can install the provider locally by running `make install`, this installs the provider under `~/.terraform.d/plugins/registry.terraform.io/fyayc/sapcc/<version>/<arch>`

## Testing the provider
At the moment, there's no official testing framework (see [#issue 113](https://github.com/hashicorp/terraform-plugin-framework/issues/113). The [helper](./helper) library provides a way for run acceptance tests aganist the mock server. The Mock server responses have been designed carefully based on the official API documentation and responses.

To run the acceptance tests, you need to either provide terraform cli path with `TF_ACC_TERRAFORM_EXEC_PATH` or provide the version of terraform to run the tests aganist with `TF_ACC_TERRAFORM_VERSION`, this is mostly useful in case of running with CI.

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
- [ ] Create new build
- [ ] Add build logs during build
- [X] Create Deployment
- [ ] Add deployment progress during deployment
- [ ] Cancel build?
- [ ] Cancel deployment?
- [ ] Create and Delete Customer properties
