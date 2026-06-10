package template_file_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateFileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateFileDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_template_file.foo", "name", "foo.seed"),
					resource.TestCheckResourceAttrSet("data.cobbler_template_file.foo", "body"),
				),
			},
		},
	})
}

const testAccTemplateFileDataSourceBasic = `
resource "cobbler_template_file" "foo" {
  name = "foo.seed"
  body = "I'm a Template file."
}

data "cobbler_template_file" "foo" {
  name = cobbler_template_file.foo.name
}
`
