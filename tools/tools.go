// +build tools

// tools is a dummy package that will be ignored for builds, but included for pulling  github.com/hashicorp/terraform-plugin-docs
package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
