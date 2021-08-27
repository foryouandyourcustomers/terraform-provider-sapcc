# Terraform Provider for SAP Commerce Cloud API

Am experimental terraform provider for experimenting with SAP Commerce Cloud API. The provider uses the [new plugin framework](https://github.com/hashicorp/terraform-plugin-framework) and is considered to be in the early stages of development. The provider is not ready to be used in the production as the upstream framework is still in `alpha` and there _will_ be breaking changes.

 - A resource, and a data source (`sapcc`),
 - Examples (`examples/`) and generated documentation (`docs/`)
 - WireMock's responses from commerce cloud apis are in `sapcc-api-mocks` 
 - Miscellaneous meta files.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.0.3
-	[Go](https://golang.org/doc/install) >= 1.16

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider 
```sh
$ make install
```
## TODOs
- [ ] Add unit tests
- [ ] Cleanup code
- [ ] Create a dedicated the `http` client 

## SAP Commerce Cloud API Support

- [X] Fetch build details
- [X] Fetch deployment Details
- [ ] Create new build 
- [ ] Add build logs during build
- [X] Create Deployment 
- [ ] Add deployment progress during deployment
- [ ] Cancel build?
- [ ] Cancel deployment?
- [ ] Create and Delete Customer properties
