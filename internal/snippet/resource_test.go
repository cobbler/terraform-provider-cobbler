package snippet_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnippetResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnippetResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_snippet.foo", "name", "foo"),
					resource.TestCheckResourceAttr("cobbler_snippet.foo", "body", "I'm a Snippet."),
				),
			},
			{
				ResourceName:                         "cobbler_snippet.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccSnippetResourceBasic = `
resource "cobbler_snippet" "foo" {
  name = "foo"
  body = "I'm a Snippet."
}
`
