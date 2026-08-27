package distro_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDistroGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "name", "foo-resource-distro-group-basic"),
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "comment", "A distro group"),
				),
			},
			{
				ResourceName:                         "cobbler_distro_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccDistroGroupResourceBasic = `
resource "cobbler_distro_group" "foo" {
  name    = "foo-resource-distro-group-basic"
  comment = "A distro group"
}
`

// TestAccDistroGroupResource_change exercises the Update lifecycle: changing the
// comment and adding/removing members from the group's "items" list.
func TestAccDistroGroupResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroGroupResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "comment", "A distro group"),
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_distro_group.foo", "items.0", "cobbler_distro.a", "uid"),
				),
			},
			{
				// Add a second member and change the comment.
				Config: testAccDistroGroupResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "comment", "A distro group updated"),
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "items.#", "2"),
					resource.TestCheckResourceAttrPair("cobbler_distro_group.foo", "items.0", "cobbler_distro.a", "uid"),
					resource.TestCheckResourceAttrPair("cobbler_distro_group.foo", "items.1", "cobbler_distro.b", "uid"),
				),
			},
			{
				// Remove the first member, keeping only the second.
				Config: testAccDistroGroupResourceChange3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_distro_group.foo", "items.0", "cobbler_distro.b", "uid"),
				),
			},
			{
				ResourceName:                         "cobbler_distro_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-group-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccDistroGroupResourceChangeDistros = `
resource "cobbler_distro" "a" {
  name       = "foo-resource-distro-group-change-a"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_distro" "b" {
  name       = "foo-resource-distro-group-change-b"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}
`

const testAccDistroGroupResourceChange1 = testAccDistroGroupResourceChangeDistros + `
resource "cobbler_distro_group" "foo" {
  name    = "foo-resource-distro-group-change"
  comment = "A distro group"
  items   = [cobbler_distro.a.uid]
}
`

const testAccDistroGroupResourceChange2 = testAccDistroGroupResourceChangeDistros + `
resource "cobbler_distro_group" "foo" {
  name    = "foo-resource-distro-group-change"
  comment = "A distro group updated"
  items   = [cobbler_distro.a.uid, cobbler_distro.b.uid]
}
`

const testAccDistroGroupResourceChange3 = testAccDistroGroupResourceChangeDistros + `
resource "cobbler_distro_group" "foo" {
  name    = "foo-resource-distro-group-change"
  comment = "A distro group updated"
  items   = [cobbler_distro.b.uid]
}
`
