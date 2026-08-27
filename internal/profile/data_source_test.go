package profile_test

import (
	"regexp"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_profile.foo", "name", "foo-data-source-profile-basic"),
					resource.TestCheckResourceAttrSet("data.cobbler_profile.foo", "distro"),
				),
			},
		},
	})
}

// TestAccProfileDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccProfileDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProfileDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler Profile not found`),
			},
		},
	})
}

const testAccProfileDataSourceNotFound = `
data "cobbler_profile" "notfound" {
  name = "does-not-exist-profile-data-source"
}
`

const testAccProfileDataSourceBasic = `
resource "cobbler_distro" "foo" {
  name       = "foo-data-source-profile-basic"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-data-source-profile-basic"
  distro = cobbler_distro.foo.uid
}

data "cobbler_profile" "foo" {
  name = cobbler_profile.foo.name
}
`
