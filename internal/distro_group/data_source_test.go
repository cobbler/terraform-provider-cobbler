package distro_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDistroGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroGroupDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_distro_group.foo", "name", "foo-ds-distro-group"),
				),
			},
		},
	})
}

const testAccDistroGroupDataSourceBasic = `
resource "cobbler_distro_group" "foo" {
  name = "foo-ds-distro-group"
}

data "cobbler_distro_group" "foo" {
  name = cobbler_distro_group.foo.name
}
`
