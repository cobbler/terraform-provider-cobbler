package system_test

import (
	"regexp"
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

// TestAccSystemDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccSystemDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSystemDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler System not found`),
			},
		},
	})
}

const testAccSystemDataSourceNotFound = `
data "cobbler_system" "notfound" {
  name = "does-not-exist-system-data-source"
}
`

const testAccSystemDataSourceBasic = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name    = "foo"
  profile = cobbler_profile.foo.uid
}

data "cobbler_system" "foo" {
  name = cobbler_system.foo.name
}
`
