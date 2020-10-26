module github.com/wearespindle/terraform-provider-cobbler

go 1.15

require (
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/wearespindle/cobblerclient v0.4.0
)

// replace github.com/wearespindle/cobblerclient => ../cobblerclient
