package provider

import (
	"terraform-provider-sapcc/helper"
	"testing"

	. "github.com/franela/goblin"
)

func TestAccDataDeployment_Basic(t *testing.T) {
	g := Goblin(t)

	dataBuild := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.1"
      source  = "fyayc/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "demo"
}

data "sapcc_deployment" "default_deploy" {
  code = "000000"
}
`, "default_deploy")

	g.Describe(`data "sapcc_deployment" "default_deploy"`, func() {
		g.It("Should exist ", func() {
			g.Assert(dataBuild.Name).IsNotNil("Data source is nil")
		})
		g.It("Should match the proper type ", func() {
			g.Assert(dataBuild.Type).Equal("sapcc_deployment", "Data source must match to `sapcc_build`")
		})
		g.It("Should match provided subscription code", func() {
			g.Assert(dataBuild.AttributeValues["subscription_code"]).Equal("demo")
		})
		g.It("Should match the default build code", func() {
			g.Assert(dataBuild.AttributeValues["build_code"]).Equal("000000.0")
		})
		g.It("Should match the default environment code", func() {
			g.Assert(dataBuild.AttributeValues["environment_code"]).Equal("d0")
		})
		g.It("Should match the default deployment strategy", func() {
			g.Assert(dataBuild.AttributeValues["strategy"]).Equal("ROLLING_UPDATE")
		})
	})
}
