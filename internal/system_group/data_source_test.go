package system_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_system_group.foo", "name", "foo-ds-system-group"),
				),
			},
		},
	})
}

const testAccSystemGroupDataSourceBasic = `
resource "cobbler_system_group" "foo" {
  name = "foo-ds-system-group"
}

data "cobbler_system_group" "foo" {
  name = cobbler_system_group.foo.name
}
`
