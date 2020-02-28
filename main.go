package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/wearespindle/terraform-provider-cobbler/cobbler"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cobbler.Provider})
}
