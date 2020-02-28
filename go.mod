module github.com/wearespindle/terraform-provider-cobbler

go 1.14

require (
	github.com/hashicorp/terraform-plugin-sdk v1.1.0
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/wearespindle/cobblerclient v0.0.0
	golang.org/x/tools v0.0.0-20200228135638-5c7c66ced534 // indirect
)

replace github.com/wearespindle/cobblerclient => ../cobblerclient
