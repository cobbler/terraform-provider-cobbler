package system_group_test

import (
	"regexp"
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

// TestAccSystemGroupDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccSystemGroupDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSystemGroupDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler SystemGroup not found`),
			},
		},
	})
}

const testAccSystemGroupDataSourceNotFound = `
data "cobbler_system_group" "notfound" {
  name = "does-not-exist-system-group-data-source"
}
`

const testAccSystemGroupDataSourceBasic = `
resource "cobbler_system_group" "foo" {
  name = "foo-ds-system-group"
}

data "cobbler_system_group" "foo" {
  name = cobbler_system_group.foo.name
}
`
