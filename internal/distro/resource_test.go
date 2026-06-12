package distro_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDistroResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro.foo", "name", "foo-resource-distro-basic"),
					resource.TestCheckResourceAttr("cobbler_distro.foo", "breed", "ubuntu"),
					resource.TestCheckResourceAttr("cobbler_distro.foo", "os_version", "focal"),
					resource.TestCheckResourceAttr("cobbler_distro.foo", "arch", "x86_64"),
				),
			},
			{
				ResourceName:                         "cobbler_distro.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccDistroResource_basicInherit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroResourceBasicInherit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro.foo", "name", "foo-resource-distro-basic-inherit"),
					resource.TestCheckResourceAttr("cobbler_distro.foo", "boot_loaders.inherited", "true"),
				),
			},
			{
				ResourceName:                         "cobbler_distro.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-basic-inherit",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccDistroResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro.foo", "comment", "I am a distro"),
				),
			},
			{
				Config: testAccDistroResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro.foo", "comment", "I am a distro again"),
				),
			},
			{
				ResourceName:                         "cobbler_distro.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccDistroResourceBasic = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-distro-basic"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}
`

const testAccDistroResourceBasicInherit = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-distro-basic-inherit"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
  boot_loaders = {
    inherited = true
  }
}
`

const testAccDistroResourceChange1 = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-distro-change"
  comment    = "I am a distro"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}
`

const testAccDistroResourceChange2 = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-distro-change"
  comment    = "I am a distro again"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}
`
