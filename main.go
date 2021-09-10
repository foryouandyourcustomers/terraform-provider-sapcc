package main

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

import (
	"context"
	"fmt"
	"os"
	"terraform-provider-sapcc/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	// this is supplied by the goreleaser - check .goreleaser.yml
	version = "dev"
)

func main() {
	err := tfsdk.Serve(context.Background(), func() tfsdk.Provider {
		return provider.New(version)
	}, tfsdk.ServeOpts{
		Name: "sapcc",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not load the plugin")
		return
	}
}
