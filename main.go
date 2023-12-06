package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-auth1/auth1"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: auth1.Provider,
	})
}
