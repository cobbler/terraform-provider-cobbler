package profile_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "name", "foo-resource-profile-group-basic"),
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "comment", "A profile group"),
				),
			},
			{
				ResourceName:                         "cobbler_profile_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileGroupResourceBasic = `
resource "cobbler_profile_group" "foo" {
  name    = "foo-resource-profile-group-basic"
  comment = "A profile group"
}
`

// TestAccProfileGroupResource_change exercises the Update lifecycle: changing the
// comment and adding/removing members from the group's "items" list.
func TestAccProfileGroupResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileGroupResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "comment", "A profile group"),
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_profile_group.foo", "items.0", "cobbler_profile.a", "uid"),
				),
			},
			{
				// Add a second member and change the comment.
				Config: testAccProfileGroupResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "comment", "A profile group updated"),
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "items.#", "2"),
					resource.TestCheckResourceAttrPair("cobbler_profile_group.foo", "items.0", "cobbler_profile.a", "uid"),
					resource.TestCheckResourceAttrPair("cobbler_profile_group.foo", "items.1", "cobbler_profile.b", "uid"),
				),
			},
			{
				// Remove the first member, keeping only the second.
				Config: testAccProfileGroupResourceChange3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "items.#", "1"),
					resource.TestCheckResourceAttrPair("cobbler_profile_group.foo", "items.0", "cobbler_profile.b", "uid"),
				),
			},
			{
				ResourceName:                         "cobbler_profile_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-group-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileGroupResourceChangeProfiles = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-group-change"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "a" {
  name   = "foo-resource-profile-group-change-a"
  distro = cobbler_distro.foo.uid
}

resource "cobbler_profile" "b" {
  name   = "foo-resource-profile-group-change-b"
  distro = cobbler_distro.foo.uid
}
`

const testAccProfileGroupResourceChange1 = testAccProfileGroupResourceChangeProfiles + `
resource "cobbler_profile_group" "foo" {
  name    = "foo-resource-profile-group-change"
  comment = "A profile group"
  items   = [cobbler_profile.a.uid]
}
`

const testAccProfileGroupResourceChange2 = testAccProfileGroupResourceChangeProfiles + `
resource "cobbler_profile_group" "foo" {
  name    = "foo-resource-profile-group-change"
  comment = "A profile group updated"
  items   = [cobbler_profile.a.uid, cobbler_profile.b.uid]
}
`

const testAccProfileGroupResourceChange3 = testAccProfileGroupResourceChangeProfiles + `
resource "cobbler_profile_group" "foo" {
  name    = "foo-resource-profile-group-change"
  comment = "A profile group updated"
  items   = [cobbler_profile.b.uid]
}
`
