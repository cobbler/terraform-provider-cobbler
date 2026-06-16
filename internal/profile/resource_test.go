package profile_test

import (
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "name", "foo-resource-profile-basic"),
					resource.TestCheckResourceAttrPair("cobbler_profile.foo", "distro", "cobbler_distro.foo", "uid"),
				),
			},
			{
				ResourceName:                         "cobbler_profile.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccProfileResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "comment", "I am a profile"),
				),
			},
			{
				Config: testAccProfileResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "comment", "I am a profile again"),
				),
			},
			{
				ResourceName:                         "cobbler_profile.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileResourceBasic = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-basic"
  breed      = "ubuntu"
  comment    = "No comment"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-basic"
  distro = cobbler_distro.foo.uid
}
`

const testAccProfileResourceChange1 = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-change"
  comment    = "I am a distro"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name    = "foo-resource-profile-change"
  comment = "I am a profile"
  distro  = cobbler_distro.foo.uid
}
`

const testAccProfileResourceChange2 = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-change"
  comment    = "I am a distro again"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name    = "foo-resource-profile-change"
  comment = "I am a profile again"
  distro  = cobbler_distro.foo.uid
}
`

// TestAccProfileResource_enableIpxeExplicit reproduces the "Provider produced inconsistent
// result after apply" error that occurs when toggling a bool-typed inherited field between
// explicit (inherited=false, value=true) and inherited (inherited=true).
// The bug: boolplanmodifier.UseStateForUnknown() on the inner `value` sub-attribute copies
// the prior state value (true) into the plan, but BoolFrom writes types.BoolNull() on apply
// when inherited=true, causing the framework to detect a plan/state mismatch.
func TestAccProfileResource_enableIpxeExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceEnableIpxeExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "name", "foo-resource-profile-enable-ipxe"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "enable_ipxe.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "enable_ipxe.value", "true"),
				),
			},
			{
				// Switch to inherited — triggers "inconsistent result after apply" without the fix.
				Config: testAccProfileResourceEnableIpxeInherited,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "enable_ipxe.inherited", "true"),
				),
			},
			{
				Config: testAccProfileResourceEnableIpxeExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "enable_ipxe.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "enable_ipxe.value", "true"),
				),
			},
			{
				ResourceName:                         "cobbler_profile.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-enable-ipxe",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileResourceEnableIpxeExplicit = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-enable-ipxe"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-enable-ipxe"
  distro = cobbler_distro.foo.uid
  enable_ipxe = {
    inherited = false
    value     = true
  }
}
`

// TestAccProfileResource_virtFileSizeExplicit reproduces the "Provider produced inconsistent
// result after apply" error for float64-typed value sub-attributes of inherited nested objects.
func TestAccProfileResource_virtFileSizeExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceVirtFileSizeExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "name", "foo-resource-profile-virt-file-size"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_file_size.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_file_size.value", "10"),
				),
			},
			{
				Config: testAccProfileResourceVirtFileSizeInherited,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_file_size.inherited", "true"),
				),
			},
			{
				Config: testAccProfileResourceVirtFileSizeExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_file_size.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_file_size.value", "10"),
				),
			},
			{
				ResourceName:                         "cobbler_profile.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-virt-file-size",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccProfileResource_virtRamExplicit reproduces the "Provider produced inconsistent
// result after apply" error for int64-typed value sub-attributes of inherited nested objects.
func TestAccProfileResource_virtRamExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceVirtRamExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "name", "foo-resource-profile-virt-ram"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_ram.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_ram.value", "2048"),
				),
			},
			{
				Config: testAccProfileResourceVirtRamInherited,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_ram.inherited", "true"),
				),
			},
			{
				Config: testAccProfileResourceVirtRamExplicit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_ram.inherited", "false"),
					resource.TestCheckResourceAttr("cobbler_profile.foo", "virt_ram.value", "2048"),
				),
			},
			{
				ResourceName:                         "cobbler_profile.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-profile-virt-ram",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

const testAccProfileResourceVirtFileSizeExplicit = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-virt-file-size"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-virt-file-size"
  distro = cobbler_distro.foo.uid
  virt_file_size = {
    inherited = false
    value     = 10
  }
}
`

const testAccProfileResourceVirtFileSizeInherited = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-virt-file-size"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-virt-file-size"
  distro = cobbler_distro.foo.uid
  virt_file_size = {
    inherited = true
  }
}
`

const testAccProfileResourceVirtRamExplicit = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-virt-ram"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-virt-ram"
  distro = cobbler_distro.foo.uid
  virt_ram = {
    inherited = false
    value     = 2048
  }
}
`

const testAccProfileResourceVirtRamInherited = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-virt-ram"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-virt-ram"
  distro = cobbler_distro.foo.uid
  virt_ram = {
    inherited = true
  }
}
`

const testAccProfileResourceEnableIpxeInherited = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-profile-enable-ipxe"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile-enable-ipxe"
  distro = cobbler_distro.foo.uid
  enable_ipxe = {
    inherited = true
  }
}
`
