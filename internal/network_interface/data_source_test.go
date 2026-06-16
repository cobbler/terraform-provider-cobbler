package network_interface_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkInterfaceDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_network_interface.eth0", "name", "eth0@foo-ds-network-interface"),
					resource.TestCheckResourceAttr("data.cobbler_network_interface.eth0", "mac_address", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttrSet("data.cobbler_network_interface.eth0", "system"),
				),
			},
		},
	})
}

const testAccNetworkInterfaceDataSourceBasic = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-ds-network-interface"
  profile = cobbler_profile.foo.name
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0@${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
}

data "cobbler_network_interface" "eth0" {
  name       = cobbler_network_interface.eth0.name
  depends_on = [cobbler_network_interface.eth0]
}
`
