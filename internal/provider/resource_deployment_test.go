package provider

import (
	"terraform-provider-sapcc/helper"
	"testing"

	. "github.com/franela/goblin"
)

func TestAccResourceDeployment_Basic(t *testing.T) {
	g := Goblin(t)

	resourceDeploy, _, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.1"
      source  = "foryouandyourcustomers/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "demo"
}

resource "sapcc_deployment" "deployment" {
  build_code = "000000.0"
  environment_code = "d0"
  strategy = "ROLLING_UPDATE"
  database_update_mode = "NONE"
}
`, "deployment")

	g.Describe(`resource "sapcc_deployment" "deployment"`, func() {
		g.It("Should exist ", func() {
			g.Assert(resourceDeploy.Name).IsNotNil("Resource is nil")
		})
		g.It("Should match the proper type ", func() {
			g.Assert(resourceDeploy.Type).Equal("sapcc_deployment", "Resource must match to `sapcc_build`")
		})
		g.It("Should match the default build code", func() {
			g.Assert(resourceDeploy.AttributeValues["build_code"]).Equal("000000.0")
		})
		g.It("Should match the default environment code", func() {
			g.Assert(resourceDeploy.AttributeValues["environment_code"]).Equal("d0")
		})
		g.It("Should match the default deployment code", func() {
			g.Assert(resourceDeploy.AttributeValues["code"]).Equal("000000")
		})
	})
}

func TestAccResourceDeployment_FailedDeployment(t *testing.T) {
	g := Goblin(t)

	_, errors, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.1"
      source  = "foryouandyourcustomers/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "demo"
}

resource "sapcc_deployment" "deployment" {
  build_code = "404"
  environment_code = "d0"
  strategy = "ROLLING_UPDATE"
  database_update_mode = "NONE"
}
`, "")

	g.Describe(`resource "sapcc_deployment" "deployment"`, func() {
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
      version = "~> 0.0.1"
      source  = "foryouandyourcustomers/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "demo"
}

resource "sapcc_deployment" "deployment" {
  build_code = "401"
  environment_code = "d0"
  strategy = "ROLLING_UPDATE"
  database_update_mode = "NONE"
}
`, "")

	g.Describe(`resource "sapcc_deployment" "deployment"`, func() {
		g.It("Testing authorized access", func() {
			g.Assert(errors).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors)).IsNotZero("Expecting at least one error")
			g.Assert(errors[0]).Equal("Unauthorized, credentials invalid for code '401', please verify your 'auth_token' and 'subscription_id'")
		})
	})
}

func TestAccResourceDeployment_VarInput(t *testing.T) {

	g := Goblin(t)

	resourceDeploy, _, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.1"
      source  = "foryouandyourcustomers/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}
provider "sapcc" {
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "demo"
}

variable "build_code" {
	default = "000000"
}

variable "environment_code" {
	default = "d0"
}

variable "strategy" {
  default = "ROLLING_UPDATE"
}

variable "database_update_mode" {
  default = "NONE"
}

resource "sapcc_deployment" "deployment" {
  build_code = var.build_code
  environment_code = var.environment_code
  strategy = var.strategy
  database_update_mode = var.database_update_mode
}
`, "")

	g.Describe(`resource "sapcc_deployment" "deployment"`, func() {
		g.It("Should exist ", func() {
			g.Assert(resourceDeploy.Name).IsNotNil("Resource is nil")
		})
		g.It("Should match the proper type ", func() {
			g.Assert(resourceDeploy.Type).Equal("sapcc_deployment", "Resource must match to `sapcc_build`")
		})
		g.It("Should match the default build code", func() {
			g.Assert(resourceDeploy.AttributeValues["build_code"]).Equal("000000.0")
		})
		g.It("Should match the default environment code", func() {
			g.Assert(resourceDeploy.AttributeValues["environment_code"]).Equal("d0")
		})
		g.It("Should match the default deployment code", func() {
			g.Assert(resourceDeploy.AttributeValues["code"]).Equal("000000")
		})
	})
}
