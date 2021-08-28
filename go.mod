module terraform-provider-sapcc

go 1.16

require (
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/docker/docker v20.10.8+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/franela/goblin v0.0.0-20210519012713-85d372ac71e2
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-plugin v1.4.1 // indirect
	github.com/hashicorp/terraform-exec v0.14.0
	github.com/hashicorp/terraform-json v0.12.0
	github.com/hashicorp/terraform-plugin-framework v0.2.1-0.20210817164910-fad6afe33058
	github.com/hashicorp/terraform-plugin-go v0.3.1
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mitchellh/go-testing-interface v1.0.4 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
)

// the support for v6 hasnt been approved/merged yet
// https://github.com/hashicorp/terraform-plugin-docs/pull/79
replace github.com/hashicorp/terraform-plugin-docs v0.4.0 => github.com/bill-rich/terraform-plugin-docs v0.4.1-0.20210819000645-ca71b522f3de
