package menu_test

import (
	"regexp"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMenuResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMenuResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_menu.foo", "name", "foo-resource-menu-basic"),
					resource.TestCheckResourceAttr("cobbler_menu.foo", "display_name", "Basic Test Menu"),
					resource.TestCheckResourceAttr("cobbler_menu.foo", "comment", "A basic menu"),
				),
			},
			{
				ResourceName:                         "cobbler_menu.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-menu-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccMenuResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMenuResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_menu.foo", "comment", "First comment"),
				),
			},
			{
				Config: testAccMenuResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_menu.foo", "comment", "Second comment"),
				),
			},
			{
				ResourceName:                         "cobbler_menu.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-menu-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccMenuResource_inherit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMenuResourceInherit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_menu.foo", "name", "foo-resource-menu-inherit"),
					resource.TestCheckResourceAttr("cobbler_menu.foo", "owners.inherited", "true"),
				),
			},
			{
				ResourceName:                         "cobbler_menu.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-menu-inherit",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccMenuResource_importNotFound exercises the zero-match branch of the resource's
// post-import name-to-uid resolution fallback: importing an ID that doesn't match any Menu on
// the server must surface a clear "not found" diagnostic instead of a raw client error.
func TestAccMenuResource_importNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMenuResourceImportNotFound,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_menu.notfound", "name", "foo-resource-menu-import-not-found"),
				),
			},
			{
				ResourceName:  "cobbler_menu.notfound",
				ImportState:   true,
				ImportStateId: "does-not-exist-menu-resource",
				ExpectError:   regexp.MustCompile(`Cobbler Menu not found`),
			},
		},
	})
}

const testAccMenuResourceImportNotFound = `
resource "cobbler_menu" "notfound" {
  name = "foo-resource-menu-import-not-found"
}
`

const testAccMenuResourceBasic = `
resource "cobbler_menu" "foo" {
  name         = "foo-resource-menu-basic"
  display_name = "Basic Test Menu"
  comment      = "A basic menu"
}
`

const testAccMenuResourceChange1 = `
resource "cobbler_menu" "foo" {
  name    = "foo-resource-menu-change"
  comment = "First comment"
}
`

const testAccMenuResourceChange2 = `
resource "cobbler_menu" "foo" {
  name    = "foo-resource-menu-change"
  comment = "Second comment"
}
`

const testAccMenuResourceInherit = `
resource "cobbler_menu" "foo" {
  name = "foo-resource-menu-inherit"
  owners = {
    inherited = true
  }
}
`
