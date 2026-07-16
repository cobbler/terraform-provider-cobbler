package distro_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDistroGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDistroGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "name", "foo-resource-distro-group-basic"),
					resource.TestCheckResourceAttr("cobbler_distro_group.foo", "comment", "A distro group"),
				),
			},
			{
				ResourceName:                         "cobbler_distro_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-distro-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccDistroGroupResourceBasic = `
resource "cobbler_distro_group" "foo" {
  name    = "foo-resource-distro-group-basic"
  comment = "A distro group"
}
`
