package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	cobbler "github.com/jtopjian/cobblerclient"
)

func TestAccCobblerRepo_basic(t *testing.T) {
	var repo cobbler.Repo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerRepo_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckRepoExists(t, "cobbler_repo.foo", &repo),
				),
			},
		},
	})
}

func TestAccCobblerRepo_change(t *testing.T) {
	var repo cobbler.Repo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerRepo_change_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckRepoExists(t, "cobbler_repo.foo", &repo),
				),
			},
			{
				Config: testAccCobblerRepo_change_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckRepoExists(t, "cobbler_repo.foo", &repo),
				),
			},
		},
	})
}

func testAccCobblerCheckRepoDestroy(s *terraform.State) error {
	config := testAccCobblerProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_repo" {
			continue
		}

		if _, err := config.cobblerClient.GetRepo(rs.Primary.ID); err == nil {
			return fmt.Errorf("Repo still exists")
		}
	}

	return nil
}

func testAccCobblerCheckRepoExists(t *testing.T, n string, repo *cobbler.Repo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccCobblerProvider.Meta().(*Config)

		found, err := config.cobblerClient.GetRepo(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Repo not found")
		}

		*repo = *found

		return nil
	}
}

var testAccCobblerRepo_basic = `
  resource "cobbler_repo" "foo" {
    name = "foo"
    breed = "apt"
    arch = "x86_64"
    apt_components = ["main"]
    apt_dists = ["trusty"]
    mirror = "http://us.archive.ubuntu.com/ubuntu/"
  }`

var testAccCobblerRepo_change_1 = `
  resource "cobbler_repo" "foo" {
    name = "foo"
    comment = "I am a repo"
    breed = "apt"
    arch = "x86_64"
    apt_components = ["main"]
    apt_dists = ["trusty"]
    mirror = "http://us.archive.ubuntu.com/ubuntu/"
  }`

var testAccCobblerRepo_change_2 = `
  resource "cobbler_repo" "foo" {
    name = "foo"
    comment = "I am a repo again"
    breed = "apt"
    arch = "x86_64"
    apt_components = ["main"]
    apt_dists = ["trusty"]
    mirror = "http://us.archive.ubuntu.com/ubuntu/"
  }`
