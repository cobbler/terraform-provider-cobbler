package image_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImageResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_image.foo", "name", "foo-resource-image-basic"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "breed", "ubuntu"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "os_version", "focal"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "arch", "x86_64"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "image_type", "iso"),
				),
			},
			{
				ResourceName:                         "cobbler_image.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-image-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccImageResource_basicInherit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageResourceBasicInherit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_image.foo", "name", "foo-resource-image-basic-inherit"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "virt_file_size.inherited", "true"),
					resource.TestCheckResourceAttr("cobbler_image.foo", "virt_ram.inherited", "true"),
				),
			},
			{
				ResourceName:                         "cobbler_image.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-image-basic-inherit",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccImageResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_image.foo", "comment", "I am an image"),
				),
			},
			{
				Config: testAccImageResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_image.foo", "comment", "I am an image again"),
				),
			},
			{
				ResourceName:                         "cobbler_image.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-image-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccImageResourceBasic = `
resource "cobbler_image" "foo" {
  name       = "foo-resource-image-basic"
  file       = "/var/www/cobbler/images/foo-basic.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
}
`

const testAccImageResourceBasicInherit = `
resource "cobbler_image" "foo" {
  name       = "foo-resource-image-basic-inherit"
  file       = "/var/www/cobbler/images/foo-inherit.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
  virt_file_size = {
    inherited = true
  }
  virt_ram = {
    inherited = true
  }
}
`

const testAccImageResourceChange1 = `
resource "cobbler_image" "foo" {
  name       = "foo-resource-image-change"
  comment    = "I am an image"
  file       = "/var/www/cobbler/images/foo-change.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
}
`

const testAccImageResourceChange2 = `
resource "cobbler_image" "foo" {
  name       = "foo-resource-image-change"
  comment    = "I am an image again"
  file       = "/var/www/cobbler/images/foo-change.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
}
`
