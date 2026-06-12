package snippet_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnippetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnippetDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_snippet.foo", "name", "foo"),
					resource.TestCheckResourceAttrSet("data.cobbler_snippet.foo", "body"),
				),
			},
		},
	})
}

const testAccSnippetDataSourceBasic = `
resource "cobbler_snippet" "foo" {
  name = "foo"
  body = "I'm a Snippet."
}

data "cobbler_snippet" "foo" {
  name = cobbler_snippet.foo.name
}
`
