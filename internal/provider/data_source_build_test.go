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
