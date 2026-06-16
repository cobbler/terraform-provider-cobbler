package template_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_template.foo", "name", "foo-ds-template"),
					resource.TestCheckResourceAttr("data.cobbler_template.foo", "template_type", "jinja2"),
				),
			},
		},
	})
}

const testAccTemplateDataSourceBasic = `
resource "cobbler_template" "foo" {
  name    = "foo-ds-template"
  uri = {
    schema = "file"
    path   = "foo-ds-template.j2"
  }
  content = "# ds test content\n"
}

data "cobbler_template" "foo" {
  name = cobbler_template.foo.name
}
`
