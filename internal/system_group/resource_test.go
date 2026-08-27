package system_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "name", "foo-resource-system-group-basic"),
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "comment", "A system group"),
				),
			},
			{
				ResourceName:                         "cobbler_system_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccSystemGroupResourceBasic = `
resource "cobbler_system_group" "foo" {
  name    = "foo-resource-system-group-basic"
  comment = "A system group"
}
`

// TestAccSystemGroupResource_change exercises the Update lifecycle: changing the
// comment and adding/removing members from the group's "items" list.
func TestAccSystemGroupResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "comment", "A system group"),
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_system_group.foo", "items.0", "cobbler_system.a", "uid"),
				),
			},
			{
				// Add a second member and change the comment.
				Config: testAccSystemGroupResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "comment", "A system group updated"),
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "items.#", "2"),
					resource.TestCheckResourceAttrPair("cobbler_system_group.foo", "items.0", "cobbler_system.a", "uid"),
					resource.TestCheckResourceAttrPair("cobbler_system_group.foo", "items.1", "cobbler_system.b", "uid"),
				),
			},
			{
				// Remove the first member, keeping only the second.
				Config: testAccSystemGroupResourceChange3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_system_group.foo", "items.0", "cobbler_system.b", "uid"),
				),
			},
			{
				ResourceName:                         "cobbler_system_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-group-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccSystemGroupResourceChangeSystems = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-system-group-change"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-system-group-change"
  distro = cobbler_distro.foo.uid
}

resource "cobbler_system" "a" {
  name    = "foo-resource-system-group-change-a"
  profile = cobbler_profile.foo.uid
}

resource "cobbler_system" "b" {
  name    = "foo-resource-system-group-change-b"
  profile = cobbler_profile.foo.uid
}
`

const testAccSystemGroupResourceChange1 = testAccSystemGroupResourceChangeSystems + `
resource "cobbler_system_group" "foo" {
  name    = "foo-resource-system-group-change"
  comment = "A system group"
  items   = [cobbler_system.a.uid]
}
`

const testAccSystemGroupResourceChange2 = testAccSystemGroupResourceChangeSystems + `
resource "cobbler_system_group" "foo" {
  name    = "foo-resource-system-group-change"
  comment = "A system group updated"
  items   = [cobbler_system.a.uid, cobbler_system.b.uid]
}
`

const testAccSystemGroupResourceChange3 = testAccSystemGroupResourceChangeSystems + `
resource "cobbler_system_group" "foo" {
  name    = "foo-resource-system-group-change"
  comment = "A system group updated"
  items   = [cobbler_system.b.uid]
}
`
