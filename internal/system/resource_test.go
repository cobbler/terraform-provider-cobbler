package system_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "name", "foo-resource-system-basic"),
					resource.TestCheckResourceAttrPair("cobbler_system.foo", "profile", "cobbler_profile.foo", "uid"),
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system"),
				),
			},
			{
				ResourceName:                         "cobbler_system.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccSystemResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system"),
				),
			},
			{
				Config: testAccSystemResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system again"),
				),
			},
			{
				ResourceName:                         "cobbler_system.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// testAccSystemDistroProfile is the shared distro+profile config used by system tests.
const testAccSystemDistroProfile = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-system"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile"
  distro = cobbler_distro.foo.uid
}
`

const testAccSystemResourceBasic = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-basic"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"
  power_id     = "foo"
}
`

const testAccSystemResourceChange1 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-change"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"
  power_id     = "foo"
}
`

const testAccSystemResourceChange2 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-change"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system again"
  power_id     = "foo"
}
`
