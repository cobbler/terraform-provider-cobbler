package distro_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDistroDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_distro.foo", "name", "foo-data-source-distro-basic"),
					resource.TestCheckResourceAttrSet("data.cobbler_distro.foo", "breed"),
					resource.TestCheckResourceAttrSet("data.cobbler_distro.foo", "arch"),
				),
			},
		},
	})
}

const testAccDistroDataSourceBasic = `
resource "cobbler_distro" "foo" {
  name       = "foo-data-source-distro-basic"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

data "cobbler_distro" "foo" {
  name = cobbler_distro.foo.name
}
`
