package cobbler

import (
	"fmt"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCobblerDistro_basic(t *testing.T) {
	var distro cobbler.Distro

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckDistroDestroy,
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

func TestAccCobblerDistro_basic_inherit(t *testing.T) {
	var distro cobbler.Distro

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckDistroDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerDistroBasicInherit,
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
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckDistroDestroy,
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_distro" {
			continue
		}
		if _, err := cobblerApiClient.GetDistro(rs.Primary.ID, false, false); err == nil {
			//goland:noinspection GoErrorStringFormat
			return fmt.Errorf("Distro still exists")
		}
	}
	return nil
}

func testAccCobblerCheckDistroExists(n string, distro *cobbler.Distro) resource.TestCheckFunc { //nolint:unparam
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		found, err := cobblerApiClient.GetDistro(rs.Primary.ID, false, false)
		if err != nil {
			return err
		}
		if found.Name != rs.Primary.ID {
			//goland:noinspection GoErrorStringFormat
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
		boot_loaders = ["ipxe"]
	}`

var testAccCobblerDistroBasicInherit = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
		boot_loaders_inherit = true
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
