package repo_test

import (
	"regexp"
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

// TestAccRepoResource_importNotFound exercises the zero-match branch of the resource's
// post-import name-to-uid resolution fallback: importing an ID that doesn't match any Repo on
// the server must surface a clear "not found" diagnostic instead of a raw client error.
func TestAccRepoResource_importNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepoResourceImportNotFound,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.notfound", "name", "foo-resource-repo-import-not-found"),
				),
			},
			{
				ResourceName:  "cobbler_repo.notfound",
				ImportState:   true,
				ImportStateId: "does-not-exist-repo-resource",
				ExpectError:   regexp.MustCompile(`Cobbler Repo not found`),
			},
		},
	})
}

const testAccRepoResourceImportNotFound = `
resource "cobbler_repo" "notfound" {
  name           = "foo-resource-repo-import-not-found"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}
`

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

// TestAccRepoResource_createrepoFlagsExplicit reproduces the "Provider produced inconsistent
// result after apply" error for string-typed value sub-attributes of inherited nested objects.
func TestAccRepoResource_createrepoFlagsExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepoResourceCreaterepoFlagsExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "name", "foo-createrepo-flags"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "createrepo_flags.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "createrepo_flags.value", "--no-database"),
				),
			},
			{
				Config: testAccRepoResourceCreaterepoFlagsInherited,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "createrepo_flags.inherited", "true"),
				),
			},
			{
				Config: testAccRepoResourceCreaterepoFlagsExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_repo.foo", "createrepo_flags.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_repo.foo", "createrepo_flags.value", "--no-database"),
				),
			},
			{
				ResourceName:                         "cobbler_repo.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-createrepo-flags",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccRepoResourceCreaterepoFlagsExplicit = `
resource "cobbler_repo" "foo" {
  name           = "foo-createrepo-flags"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
  createrepo_flags = {
    inherited = false
    value     = "--no-database"
  }
}
`

const testAccRepoResourceCreaterepoFlagsInherited = `
resource "cobbler_repo" "foo" {
  name           = "foo-createrepo-flags"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
  createrepo_flags = {
    inherited = true
  }
}
`
