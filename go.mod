module github.com/cobbler/terraform-provider-cobbler

go 1.15

require (
	github.com/cobbler/cobblerclient v0.4.2
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	golang.org/x/tools v0.0.0-20201121010211-780cb80bd7fb // indirect
)

// replace github.com/cobbler/cobblerclient => ../cobblerclient
