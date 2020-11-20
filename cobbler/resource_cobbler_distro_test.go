package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	cobbler "github.com/cobbler/cobblerclient"
)

func TestAccCobblerDistro_basic(t *testing.T) {
	var distro cobbler.Distro

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckDistroDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerDistroBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
				),
			},
		},
	})
}

func TestAccCobblerDistro_change(t *testing.T) {
	var distro cobbler.Distro

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckDistroDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerDistroChange1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
				),
			},
			{
				Config: testAccCobblerDistroChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
				),
			},
		},
	})
}

func testAccCobblerCheckDistroDestroy(s *terraform.State) error {
	config := testAccCobblerProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_distro" {
			continue
		}

		if _, err := config.cobblerClient.GetDistro(rs.Primary.ID); err == nil {
			return fmt.Errorf("Distro still exists")
		}
	}

	return nil
}

func testAccCobblerCheckDistroExists(t *testing.T, n string, distro *cobbler.Distro) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccCobblerProvider.Meta().(*Config)

		found, err := config.cobblerClient.GetDistro(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Distro not found")
		}

		*distro = *found

		return nil
	}
}

var testAccCobblerDistroBasic = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
	}`

var testAccCobblerDistroChange1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro"
		breed = "ubuntu"
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
	}`

var testAccCobblerDistroChange2 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro again"
		breed = "ubuntu"
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
	}`
