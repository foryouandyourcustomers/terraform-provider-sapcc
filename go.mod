module terraform-provider-sapcc

go 1.16

require (
	cloud.google.com/go v0.61.0 // indirect
	github.com/aws/aws-sdk-go v1.25.3 // indirect
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-plugin v1.4.1 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-framework v0.2.1-0.20210817164910-fad6afe33058
	github.com/hashicorp/terraform-plugin-go v0.3.1
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mitchellh/go-testing-interface v1.0.4 // indirect
)

// the support for v6 hasnt been approved/merged yet
// https://github.com/hashicorp/terraform-plugin-docs/pull/79
replace github.com/hashicorp/terraform-plugin-docs v0.4.0 => github.com/bill-rich/terraform-plugin-docs v0.4.1-0.20210819000645-ca71b522f3de
