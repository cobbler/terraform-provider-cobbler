package network_interface_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkInterfaceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "name", "eth0@foo-resource-network-interface-basic"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "mac_address", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "static", "true"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.address", "1.2.3.4"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.netmask", "255.255.255.0"),
				),
			},
			{
				ResourceName:                         "cobbler_network_interface.eth0",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "eth0@foo-resource-network-interface-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccNetworkInterfaceResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.address", "1.2.3.4"),
				),
			},
			{
				Config: testAccNetworkInterfaceResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.address", "1.2.3.5"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "dns.name", "host.example.com"),
				),
			},
		},
	})
}

const testAccNetworkInterfaceDistroProfileSystem = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-network-interface"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-network-interface"
  distro = cobbler_distro.foo.name
}
`

const testAccNetworkInterfaceResourceBasic = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-basic"
  profile = cobbler_profile.foo.name
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0@${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "1.2.3.4"
    netmask = "255.255.255.0"
  }
}
`

const testAccNetworkInterfaceResourceChange1 = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-change"
  profile = cobbler_profile.foo.name
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0@${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "1.2.3.4"
    netmask = "255.255.255.0"
  }
}
`

const testAccNetworkInterfaceResourceChange2 = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-change"
  profile = cobbler_profile.foo.name
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0@${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "1.2.3.5"
    netmask = "255.255.255.0"
  }
  dns = {
    name = "host.example.com"
  }
}
`
