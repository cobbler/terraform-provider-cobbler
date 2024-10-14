package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	cobbler "github.com/cobbler/cobblerclient"
)

func TestAccCobblerSystem_basic(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile
	var system cobbler.System

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo", &system),
				),
			},
		},
	})
}

func TestAccCobblerSystem_multi(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile
	var system cobbler.System

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo.45", &system),
				),
			},
		},
	})
}

func TestAccCobblerSystem_change(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile
	var system cobbler.System

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemChange1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo", &system),
				),
			},
			{
				Config: testAccCobblerSystemChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo", &system),
				),
			},
		},
	})
}

func TestAccCobblerSystem_removeInterface(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile
	var system cobbler.System

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemRemoveInterface1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo", &system),
				),
			},
			{
				Config: testAccCobblerSystemRemoveInterface2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists("cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists("cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists("cobbler_system.foo", &system),
				),
			},
		},
	})
}

func testAccCobblerCheckSystemDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_system" {
			continue
		}

		if _, err := cobblerApiClient.GetSystem(rs.Primary.ID, false, false); err == nil {
			//goland:noinspection GoErrorStringFormat
			return fmt.Errorf("System still exists")
		}
	}

	return nil
}

func testAccCobblerCheckSystemExists(n string, system *cobbler.System) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		found, err := cobblerApiClient.GetSystem(rs.Primary.ID, false, false)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			//goland:noinspection GoErrorStringFormat
			return fmt.Errorf("System not found")
		}

		*system = *found

		return nil
	}
}

var testAccCobblerSystemBasic = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}

	resource "cobbler_system" "foo" {
		name = "foo"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		comment = "I'm a system"
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.4"
			netmask = "255.255.255.0"
		}

		interface {
			name = "eth1"
			mac_address = "aa:bb:cc:dd:ee:fa"
			static = true
			ip_address = "1.2.3.5"
			netmask = "255.255.255.0"
		}

	}`

var testAccCobblerSystemMulti = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}

	resource "cobbler_system" "foo" {
		count = 50
		name = "${format("foo-%d", count.index)}"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		comment = "I'm a system"
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:${format("%d", count.index)}"
		}

		interface {
			name = "eth1"
			mac_address = "aa:bb:cc:dd:ef:${format("%d", count.index)}"
		}
	}`

var testAccCobblerSystemChange1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}

	resource "cobbler_system" "foo" {
		name = "foo"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		comment = "I'm a system"
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.4"
			netmask = "255.255.255.0"
		}

		interface {
			name = "eth1"
			mac_address = "aa:bb:cc:dd:ee:fa"
			static = true
			ip_address = "1.2.3.5"
			netmask = "255.255.255.0"
		}

	}`

var testAccCobblerSystemChange2 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}

	resource "cobbler_system" "foo" {
		name = "foo"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		comment = "I'm a system again"
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.6"
			netmask = "255.255.255.0"
		}

		interface {
			name = "eth1"
			mac_address = "aa:bb:cc:dd:ee:fa"
			static = true
			ip_address = "1.2.3.5"
			netmask = "255.255.255.0"
		}

	}`

var testAccCobblerSystemRemoveInterface1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name

	}

	resource "cobbler_system" "foo" {
		name = "foo"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.4"
			netmask = "255.255.255.0"
		}

		interface {
			name = "eth1"
			mac_address = "aa:bb:cc:dd:ee:fa"
			static = true
			ip_address = "1.2.3.5"
			netmask = "255.255.255.0"
			management = true
		}

	}`

var testAccCobblerSystemRemoveInterface2 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "focal"
		arch = "x86_64"
        boot_loaders = ["ipxe"]
		kernel = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
		initrd = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
	}

	resource "cobbler_profile" "foo" {
		name = "foo"
		distro = cobbler_distro.foo.name
	}

	resource "cobbler_system" "foo" {
		name = "foo"
		profile = "${cobbler_profile.foo.name}"
		name_servers = ["8.8.8.8", "8.8.4.4"]
		power_id = "foo"

		interface {
			name = "default"
		}

		interface {
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.4"
			netmask = "255.255.255.0"
		}
	}`
