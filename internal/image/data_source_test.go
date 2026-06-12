package image_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_image.foo", "name", "foo-data-source-image-basic"),
					resource.TestCheckResourceAttrSet("data.cobbler_image.foo", "breed"),
					resource.TestCheckResourceAttrSet("data.cobbler_image.foo", "arch"),
					resource.TestCheckResourceAttrSet("data.cobbler_image.foo", "image_type"),
				),
			},
		},
	})
}

const testAccImageDataSourceBasic = `
resource "cobbler_image" "foo" {
  name       = "foo-data-source-image-basic"
  file       = "/var/www/cobbler/images/foo-data-source.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
}

data "cobbler_image" "foo" {
  name = cobbler_image.foo.name
}
`
