module github.com/wearespindle/terraform-provider-cobbler

go 1.13

require (
	github.com/hashicorp/terraform v0.12.20
	github.com/wearespindle/cobblerclient v0.0.0
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc // indirect
)

replace github.com/wearespindle/cobblerclient => ../cobblerclient
