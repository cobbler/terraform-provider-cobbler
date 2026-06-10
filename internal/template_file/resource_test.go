package template_file_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateFileResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateFileResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_template_file.foo", "name", "foo.seed"),
					resource.TestCheckResourceAttr("cobbler_template_file.foo", "body", "I'm a Template file."),
				),
			},
			{
				ResourceName:                         "cobbler_template_file.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo.seed",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccTemplateFileResourceBasic = `
resource "cobbler_template_file" "foo" {
  name = "foo.seed"
  body = "I'm a Template file."
}
`
