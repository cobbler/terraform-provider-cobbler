package template_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTemplateResource_basic (and _change, and the template data source test) are currently
// expected to fail on "content" after apply: Cobbler 4.0.0a3's Template.content setter
// (cobbler/items/template.py) writes the file using uri.path directly instead of resolving it
// against autoinstall_templates_dir like the path validator does, so the write lands in
// cobblerd's cwd rather than the real template directory - and since to_dict() never includes
// "content", the in-memory value doesn't survive the save_template round trip either. This is an
// upstream Cobbler bug, not a client/provider bug; re-check against a newer Cobbler 4.0.0 build.
func TestAccTemplateResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_template.foo", "name", "foo-resource-template-basic"),
					resource.TestCheckResourceAttr("cobbler_template.foo", "template_type", "jinja"),
					resource.TestCheckResourceAttr("cobbler_template.foo", "uri.schema", "file"),
					resource.TestCheckResourceAttr("cobbler_template.foo", "content", "# original content\n"),
				),
			},
			{
				ResourceName:                         "cobbler_template.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-template-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccTemplateResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_template.foo", "content", "# first revision\n"),
				),
			},
			{
				Config: testAccTemplateResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_template.foo", "content", "# second revision\n"),
				),
			},
		},
	})
}

const testAccTemplateResourceBasic = `
resource "cobbler_template" "foo" {
  name    = "foo-resource-template-basic"
  uri = {
    schema = "file"
    path   = "foo-resource-template-basic.j2"
  }
  content = "# original content\n"
}
`

const testAccTemplateResourceChange1 = `
resource "cobbler_template" "foo" {
  name    = "foo-resource-template-change"
  uri = {
    schema = "file"
    path   = "foo-resource-template-change.j2"
  }
  content = "# first revision\n"
}
`

const testAccTemplateResourceChange2 = `
resource "cobbler_template" "foo" {
  name    = "foo-resource-template-change"
  uri = {
    schema = "file"
    path   = "foo-resource-template-change.j2"
  }
  content = "# second revision\n"
}
`
