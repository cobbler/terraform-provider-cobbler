package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	cobbler "github.com/wearespindle/cobblerclient"
)

func TestAccCobblerProfile_basic(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func TestAccCobblerProfile_change(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileChange1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
				),
			},
			{
				Config: testAccCobblerProfileChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func TestAccCobblerProfile_withRepo(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerProfileWithRepo,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
				),
			},
		},
	})
}

func testAccCobblerCheckProfileDestroy(s *terraform.State) error {
	config := testAccCobblerProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_profile" {
			continue
		}

		if _, err := config.cobblerClient.GetProfile(rs.Primary.ID); err == nil {
			return fmt.Errorf("Profile still exists")
		}
	}

	return nil
}

func testAccCobblerCheckProfileExists(t *testing.T, n string, profile *cobbler.Profile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccCobblerProvider.Meta().(*Config)

		found, err := config.cobblerClient.GetProfile(rs.Primary.ID)
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		comment = "I am a profile again"
		distro = cobbler_distro.foo.name
		repos = ["Ubuntu-18.04-x86_64"]
	}`
