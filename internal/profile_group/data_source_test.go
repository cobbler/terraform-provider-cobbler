package profile_group_test

import (
	"regexp"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileGroupDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_profile_group.foo", "name", "foo-ds-profile-group"),
				),
			},
		},
	})
}

// TestAccProfileGroupDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccProfileGroupDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProfileGroupDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler ProfileGroup not found`),
			},
		},
	})
}

const testAccProfileGroupDataSourceNotFound = `
data "cobbler_profile_group" "notfound" {
  name = "does-not-exist-profile-group-data-source"
}
`

const testAccProfileGroupDataSourceBasic = `
resource "cobbler_profile_group" "foo" {
  name = "foo-ds-profile-group"
}

data "cobbler_profile_group" "foo" {
  name = cobbler_profile_group.foo.name
}
`
