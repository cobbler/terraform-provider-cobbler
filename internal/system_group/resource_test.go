package system_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "name", "foo-resource-system-group-basic"),
					resource.TestCheckResourceAttr("cobbler_system_group.foo", "comment", "A system group"),
				),
			},
			{
				ResourceName:                         "cobbler_system_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccSystemGroupResourceBasic = `
resource "cobbler_system_group" "foo" {
  name    = "foo-resource-system-group-basic"
  comment = "A system group"
}
`
