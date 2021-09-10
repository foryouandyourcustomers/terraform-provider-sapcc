package provider

import (
	"terraform-provider-sapcc/helper"
	"testing"

	. "github.com/franela/goblin"
)

func TestAccDataDeployment_Basic(t *testing.T) {
	g := Goblin(t)

	dataDeploy, _, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.0"
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
			g.Assert(dataDeploy.Name).IsNotNil("Data source is nil")
		})
		g.It("Should match the proper type ", func() {
			g.Assert(dataDeploy.Type).Equal("sapcc_deployment", "Data source must match to `sapcc_build`")
		})
		g.It("Should match provided subscription code", func() {
			g.Assert(dataDeploy.AttributeValues["subscription_code"]).Equal("demo")
		})
		g.It("Should match the default build code", func() {
			g.Assert(dataDeploy.AttributeValues["build_code"]).Equal("000000.0")
		})
		g.It("Should match the default environment code", func() {
			g.Assert(dataDeploy.AttributeValues["environment_code"]).Equal("d0")
		})
		g.It("Should match the default code strategy", func() {
			g.Assert(dataDeploy.AttributeValues["strategy"]).Equal("ROLLING_UPDATE")
		})
	})
}

func TestAccDataDeployment_FailedDeployment(t *testing.T) {
	g := Goblin(t)

	_, errors, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.0"
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

data "sapcc_deployment" "build_doesnt_exist" {
  code = "404"
}
`, "")

	g.Describe(`data "sapcc_deployment" "build_doesnt_exist"`, func() {
		g.It("Testing unknown builds ", func() {
			g.Assert(errors).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors)).IsNotZero("Expecting at least one error")
			g.Assert(errors[0]).Equal("Deployment or progress not found; code '404'; Check logs or report it provider developer")
		})
	})

	_, errors, _ = helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.0"
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

data "sapcc_deployment" "build_unauth" {
  code = "401"
}
`, "")

	g.Describe(`data "sapcc_deployment" "build_unauth"`, func() {
		g.It("Testing authorized access", func() {
			g.Assert(errors).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors)).IsNotZero("Expecting at least one error")
			g.Assert(errors[0]).Equal("Unauthorized, credentials invalid for code '401', please verify your 'auth_token' and 'subscription_id'")
		})
	})

	_, errors, _ = helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.0"
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

data "sapcc_deployment" "build_api_error" {
  code = "500"
}
`, "")

	g.Describe(`data "sapcc_deployment" "build_api_error"`, func() {
		g.It("Testing upstream Api Error ", func() {
			g.Assert(errors).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors)).IsNotZero("Expecting at least one error")
			g.Assert(errors[0]).Equal("Unexpected http status 500 for code '500' from upstream api; won't continue. expected 200 ")
		})
	})
}
