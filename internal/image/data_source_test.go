package image_test

import (
	"regexp"
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

// TestAccImageDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccImageDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccImageDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler Image not found`),
			},
		},
	})
}

const testAccImageDataSourceNotFound = `
data "cobbler_image" "notfound" {
  name = "does-not-exist-image-data-source"
}
`

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
