package system_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_system.foo", "name", "foo"),
					resource.TestCheckResourceAttrSet("data.cobbler_system.foo", "profile"),
				),
			},
		},
	})
}

const testAccSystemDataSourceBasic = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name    = "foo"
  profile = cobbler_profile.foo.name
}

data "cobbler_system" "foo" {
  name = cobbler_system.foo.name
}
`
