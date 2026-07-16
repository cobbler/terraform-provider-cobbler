package profile_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileGroupResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "name", "foo-resource-profile-group-basic"),
					resource.TestCheckResourceAttr("cobbler_profile_group.foo", "comment", "A profile group"),
				),
			},
			{
				ResourceName:                         "cobbler_profile_group.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-group-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileGroupResourceBasic = `
resource "cobbler_profile_group" "foo" {
  name    = "foo-resource-profile-group-basic"
  comment = "A profile group"
}
`
