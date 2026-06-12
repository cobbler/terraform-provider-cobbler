package repo_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepoResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepoResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "name", "foo"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "breed", "apt"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "arch", "x86_64"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "mirror", "http://us.archive.ubuntu.com/ubuntu/"),
				),
			},
			{
				ResourceName:                         "cobbler_repo.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccRepoResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepoResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "comment", "I am a repo"),
				),
			},
			{
				Config: testAccRepoResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "comment", "I am a repo again"),
				),
			},
			{
				ResourceName:                         "cobbler_repo.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccRepoResourceBasic = `
resource "cobbler_repo" "foo" {
  name           = "foo"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}
`

const testAccRepoResourceChange1 = `
resource "cobbler_repo" "foo" {
  name           = "foo"
  comment        = "I am a repo"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}
`

const testAccRepoResourceChange2 = `
resource "cobbler_repo" "foo" {
  name           = "foo"
  comment        = "I am a repo again"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}
`
