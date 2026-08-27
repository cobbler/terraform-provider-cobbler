package distro_group_test

import (
	"regexp"
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

// TestAccDistroGroupDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccDistroGroupDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDistroGroupDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler DistroGroup not found`),
			},
		},
	})
}

const testAccDistroGroupDataSourceNotFound = `
data "cobbler_distro_group" "notfound" {
  name = "does-not-exist-distro-group-data-source"
}
`

const testAccDistroGroupDataSourceBasic = `
resource "cobbler_distro_group" "foo" {
  name = "foo-ds-distro-group"
}

data "cobbler_distro_group" "foo" {
  name = cobbler_distro_group.foo.name
}
`
