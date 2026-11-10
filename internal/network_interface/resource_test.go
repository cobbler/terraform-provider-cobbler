package network_interface_test

import (
	"regexp"
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
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "name", "eth0-foo-resource-network-interface-basic"),
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
				ImportStateId:                        "eth0-foo-resource-network-interface-basic",
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

// TestAccNetworkInterfaceResource_gateway exercises the per-address-family default gateway,
// which Cobbler 4.0.0 models as IPv4Option.DefaultGateway / IPv6Option.DefaultGateway (see
// cobbler/items/options/ip.py) rather than as a standalone NetworkInterface field.
func TestAccNetworkInterfaceResource_gateway(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceResourceGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.address", "1.2.3.4"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv4.gateway", "1.2.3.1"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv6.address", "2001:db8::4"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0", "ipv6.default_gateway", "2001:db8::1"),
				),
			},
		},
	})
}

// TestAccNetworkInterfaceResource_importNotFound exercises the zero-match branch of the
// resource's post-import name-to-uid resolution fallback: importing an ID that doesn't match
// any NetworkInterface on the server must surface a clear "not found" diagnostic instead of a
// raw client error.
func TestAccNetworkInterfaceResource_importNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceResourceImportNotFound,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.notfound", "name", "eth0-foo-resource-network-interface-import-not-found"),
				),
			},
			{
				ResourceName:  "cobbler_network_interface.notfound",
				ImportState:   true,
				ImportStateId: "does-not-exist-network-interface-resource",
				ExpectError:   regexp.MustCompile(`Cobbler NetworkInterface not found`),
			},
		},
	})
}

// TestAccNetworkInterfaceResource_importComposite exercises importing a NetworkInterface
// whose name ("eth0") collides with another system's interface of the same name — the
// overwhelmingly common real-world case, since most systems name their primary interface
// eth0. A bare-name import can't disambiguate between the two and fails with "Ambiguous
// Cobbler NetworkInterface name"; a composite "<system_name>/<interface_name>" import ID
// must resolve the system name to a uid first and scope the NetworkInterface lookup to
// that system, importing the correct interface (verified via its distinguishing
// mac_address, not just "something" got imported).
func TestAccNetworkInterfaceResource_importComposite(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceResourceImportComposite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0_a", "mac_address", "aa:bb:cc:dd:ee:a1"),
					resource.TestCheckResourceAttr("cobbler_network_interface.eth0_b", "mac_address", "aa:bb:cc:dd:ee:b1"),
				),
			},
			{
				ResourceName:                         "cobbler_network_interface.eth0_b",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-network-interface-import-composite-b/eth0",
				ImportStateVerifyIdentifierAttribute: "mac_address",
			},
		},
	})
}

const testAccNetworkInterfaceResourceImportComposite = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "composite_a" {
  name    = "foo-resource-network-interface-import-composite-a"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_system" "composite_b" {
  name    = "foo-resource-network-interface-import-composite-b"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "eth0_a" {
  name        = "eth0"
  system      = cobbler_system.composite_a.uid
  mac_address = "aa:bb:cc:dd:ee:a1"
}

resource "cobbler_network_interface" "eth0_b" {
  name        = "eth0"
  system      = cobbler_system.composite_b.uid
  mac_address = "aa:bb:cc:dd:ee:b1"
}
`

const testAccNetworkInterfaceResourceImportNotFound = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "notfound" {
  name    = "foo-resource-network-interface-import-not-found"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "notfound" {
  name        = "eth0-${cobbler_system.notfound.name}"
  system      = cobbler_system.notfound.uid
  mac_address = "aa:bb:cc:dd:ee:00"
}
`

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
  distro = cobbler_distro.foo.uid
}
`

const testAccNetworkInterfaceResourceBasic = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-basic"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0-${cobbler_system.foo.name}"
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
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0-${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "1.2.3.4"
    netmask = "255.255.255.0"
  }
}
`

const testAccNetworkInterfaceResourceGateway = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-gateway"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0-${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:fe"
  static      = true
  ipv4 = {
    address = "1.2.3.4"
    netmask = "255.255.255.0"
    gateway = "1.2.3.1"
  }
  ipv6 = {
    address         = "2001:db8::4"
    prefix          = "64"
    default_gateway = "2001:db8::1"
  }
}
`

const testAccNetworkInterfaceResourceChange2 = testAccNetworkInterfaceDistroProfileSystem + `
resource "cobbler_system" "foo" {
  name    = "foo-resource-network-interface-change"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0-${cobbler_system.foo.name}"
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
