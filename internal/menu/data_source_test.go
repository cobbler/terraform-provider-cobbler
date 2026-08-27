package menu_test

import (
	"regexp"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMenuDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMenuDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_menu.foo", "name", "foo-datasource-menu-basic"),
					resource.TestCheckResourceAttr("data.cobbler_menu.foo", "display_name", "Data Source Menu"),
				),
			},
		},
	})
}

// TestAccMenuDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccMenuDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMenuDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler Menu not found`),
			},
		},
	})
}

const testAccMenuDataSourceNotFound = `
data "cobbler_menu" "notfound" {
  name = "does-not-exist-menu-data-source"
}
`

const testAccMenuDataSourceBasic = `
resource "cobbler_menu" "foo" {
  name         = "foo-datasource-menu-basic"
  display_name = "Data Source Menu"
}

data "cobbler_menu" "foo" {
  name = cobbler_menu.foo.name
}
`
