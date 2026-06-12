package menu_test

import (
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

const testAccMenuDataSourceBasic = `
resource "cobbler_menu" "foo" {
  name         = "foo-datasource-menu-basic"
  display_name = "Data Source Menu"
}

data "cobbler_menu" "foo" {
  name = cobbler_menu.foo.name
}
`
