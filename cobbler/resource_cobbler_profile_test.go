package cobbler

import (
	"fmt"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCobblerProfile_basic(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func TestAccCobblerProfile_change(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileChange1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
				),
			},
			{
				Config: testAccCobblerProfileChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func TestAccCobblerProfile_withRepo(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileWithRepo,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func testAccCobblerCheckProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_profile" {
			continue
		}

		if _, err := cobblerApiClient.GetProfile(rs.Primary.ID); err == nil {
			return fmt.Errorf("Profile still exists")
		}
	}

	return nil
}

func testAccCobblerCheckProfileExists(n string, profile *cobbler.Profile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		found, err := cobblerApiClient.GetProfile(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Profile not found")
		}

		*profile = *found

		return nil
	}
}

var testAccCobblerProfileBasic = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		comment = "No comment"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}`

var testAccCobblerProfileChange1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		comment = "I am a profile"
		distro = cobbler_distro.foo.name
	}`

var testAccCobblerProfileChange2 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro again"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		comment = "I am a profile again"
		distro = cobbler_distro.foo.name
	}`

var testAccCobblerProfileWithRepo = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro all over again"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		comment = "I am a profile again"
		distro = cobbler_distro.foo.name
		repos = ["Ubuntu-20.04-x86_64"]
	}`
