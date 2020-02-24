package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	cobbler "github.com/wearespindle/cobblerclient"
)

func TestAccCobblerSystem_basic(t *testing.T) {
	var distro cobbler.Distro
	var profile cobbler.Profile
	var system cobbler.System

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo", &system),
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
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo.45", &system),
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
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemChange1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo", &system),
				),
			},
			{
				Config: testAccCobblerSystemChange2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo", &system),
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
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerSystemRemoveInterface1,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo", &system),
					resource.TestCheckResourceAttr(
						"cobbler_system.foo", "interface.586365610.management", "true"),
				),
			},
			{
				Config: testAccCobblerSystemRemoveInterface2,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckDistroExists(t, "cobbler_distro.foo", &distro),
					testAccCobblerCheckProfileExists(t, "cobbler_profile.foo", &profile),
					testAccCobblerCheckSystemExists(t, "cobbler_system.foo", &system),
				),
			},
		},
	})
}

func testAccCobblerCheckSystemDestroy(s *terraform.State) error {
	config := testAccCobblerProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_system" {
			continue
		}

		if _, err := config.cobblerClient.GetSystem(rs.Primary.ID); err == nil {
			return fmt.Errorf("System still exists")
		}
	}

	return nil
}

func testAccCobblerCheckSystemExists(t *testing.T, n string, system *cobbler.System) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccCobblerProvider.Meta().(*Config)

		found, err := config.cobblerClient.GetSystem(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
			name = "eth0"
		}

		interface {
			name = "eth1"
		}
	}`

var testAccCobblerSystemChange1 = `
	resource "cobbler_distro" "foo" {
		name = "foo"
		breed = "ubuntu"
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
		os_version = "bionic"
		arch = "x86_64"
		boot_loader = "grub"
		kernel = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
		initrd = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
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
			name = "eth0"
			mac_address = "aa:bb:cc:dd:ee:ff"
			static = true
			ip_address = "1.2.3.4"
			netmask = "255.255.255.0"
		}
	}`
