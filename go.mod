module github.com/wearespindle/terraform-provider-cobbler

go 1.14

require (
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/wearespindle/cobblerclient v0.0.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
)

replace github.com/wearespindle/cobblerclient => ../cobblerclient
