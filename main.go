package main

import (
	"github.com/cobbler/terraform-provider-cobbler/cobbler"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cobbler.Provider})
}
