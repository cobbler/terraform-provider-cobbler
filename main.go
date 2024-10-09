package main

import (
	"flag"
	"fmt"
	"github.com/cobbler/terraform-provider-cobbler/cobbler"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version = "dev"

	// goreleaser can also pass the specific commit if you want
	commit = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cobbler.New(fmt.Sprintf("%s-%s", version, commit)),
		Debug:        debugMode,
	})
}
