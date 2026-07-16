package profile_group_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 4, 0, 0)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileGroupDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_profile_group.foo", "name", "foo-ds-profile-group"),
				),
			},
		},
	})
}

const testAccProfileGroupDataSourceBasic = `
resource "cobbler_profile_group" "foo" {
  name = "foo-ds-profile-group"
}

data "cobbler_profile_group" "foo" {
  name = cobbler_profile_group.foo.name
}
`
