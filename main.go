package main

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"os"
	"terraform-provider-sapcc/sapcc"
)

func main() {
	err := tfsdk.Serve(context.Background(), sapcc.New, tfsdk.ServeOpts{
		Name: "sapcc",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not load the plugin")
		return
	}
}
