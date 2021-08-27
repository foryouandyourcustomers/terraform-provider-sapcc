module terraform-provider-sapcc

go 1.16

require (
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-plugin v1.4.1 // indirect
	github.com/hashicorp/terraform-plugin-framework v0.2.1-0.20210817164910-fad6afe33058
	github.com/hashicorp/terraform-plugin-go v0.3.1
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mitchellh/go-testing-interface v1.0.4 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210326060303-6b1517762897 // indirect
	golang.org/x/sys v0.0.0-20210502180810-71e4cd670f79 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200711021454-869866162049 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// the support for v6 hasnt been approved/merged yet
// https://github.com/hashicorp/terraform-plugin-docs/pull/79
replace github.com/hashicorp/terraform-plugin-docs v0.4.0 => github.com/bill-rich/terraform-plugin-docs v0.4.1-0.20210819000645-ca71b522f3de
