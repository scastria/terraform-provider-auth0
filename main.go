package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-auth0/auth0"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: auth0.Provider,
	})
}
