package provider

// The new framework doesn't yet provide the testing "resource" helper for acceptance testing
import (
	"os"
	"terraform-provider-sapcc/helper"
	"testing"

	. "github.com/franela/goblin"
)

// Start the mock server to simulate real server
func TestMain(m *testing.M) {
	helper.StartMockServer()
	// os.Exit() does not respect defer statements
	ret := m.Run()

	helper.StopMockServer()
	os.Exit(ret)
}

func TestAccDataBuild_Basic(t *testing.T) {
	g := Goblin(t)

	dataBuild, _, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.2"
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

data "sapcc_build" "default_build" {
  code = "000000.0"
}
`, "default_build")

	g.Describe(`data "sapcc_build" "default_build"`, func() {

		g.It("Should exist ", func() {
			g.Assert(dataBuild.Name).IsNotNil("Data source is nil")
		})
		g.It("Should match the proper type ", func() {
			g.Assert(dataBuild.Type).Equal("sapcc_build", "Data source must match to `sapcc_build`")
		})
		g.It("Should match the default build name", func() {
			g.Assert(dataBuild.AttributeValues["name"]).Equal("ci-build-0")
		})
		g.It("Should match provided subscription code", func() {
			g.Assert(dataBuild.AttributeValues["subscription_code"]).Equal("demo")
		})
	})
}

func TestAccDataBuild_FailedBuilds(t *testing.T) {
	g := Goblin(t)

	_, errors, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.2"
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

data "sapcc_build" "build_doesnt_exist" {
  code = "404"
}
`, "")

	g.Describe(`data "sapcc_build" "build_doesnt_exist"`, func() {
		g.It("Testing unknown deployments ", func() {
			g.Assert(errors).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors)).IsNotZero("Expecting at least one error")
			g.Assert(errors[0]).Equal("Build '404' not found")
		})
	})

	_, errors2, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.2"
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

data "sapcc_build" "build_unauth" {
  code = "401"
}
`, "")

	g.Describe(`data "sapcc_build" "build_unauth"`, func() {
		g.It("Testing authorized access", func() {
			g.Assert(errors2).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors2)).IsNotZero("Expecting at least one error")
			g.Assert(errors2[0]).Equal("Unauthorized, credentials invalid for build '401', please verify your 'auth_token' and 'subscription_id'")
		})
	})

	_, errors3, _ := helper.ResourceTest(t, `
terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.2"
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

data "sapcc_build" "build_api_error" {
  code = "500"
}
`, "")

	g.Describe(`data "sapcc_build" "build_api_error"`, func() {
		g.It("Testing upstream Api Error ", func() {
			g.Assert(errors3).IsNotNil("Expecting errors not be nil")
			g.Assert(len(errors3)).IsNotZero("Expecting at least one error")
			g.Assert(errors3[0]).Equal("Unexpected http status 500 for build '500' from upstream api; won't continue. expected 200 ")
		})
	})
}
