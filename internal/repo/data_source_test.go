package repo_test

import (
	"regexp"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepoDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepoDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cobbler_repo.foo", "name", "foo"),
					resource.TestCheckResourceAttrSet("data.cobbler_repo.foo", "breed"),
					resource.TestCheckResourceAttrSet("data.cobbler_repo.foo", "mirror"),
				),
			},
		},
	})
}

// TestAccRepoDataSource_notFound exercises the zero-match branch of the data source's
// name-to-uid resolution: looking up a name that doesn't exist on the server must surface a
// clear "not found" diagnostic instead of a raw client error.
func TestAccRepoDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRepoDataSourceNotFound,
				ExpectError: regexp.MustCompile(`Cobbler Repo not found`),
			},
		},
	})
}

const testAccRepoDataSourceNotFound = `
data "cobbler_repo" "notfound" {
  name = "does-not-exist-repo-data-source"
}
`

const testAccRepoDataSourceBasic = `
resource "cobbler_repo" "foo" {
  name           = "foo"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}

data "cobbler_repo" "foo" {
  name = cobbler_repo.foo.name
}
`
