package cobbler

import (
	"fmt"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
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
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
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
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
				),
			},
			{
				Config: testAccCobblerDistroChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
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

func testAccCobblerCheckDistroExists(n string, distro *cobbler.Distro) resource.TestCheckFunc {
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
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}`

var testAccCobblerDistroChange1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}`

var testAccCobblerDistroChange2 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		comment = "I am a distro again"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}`
