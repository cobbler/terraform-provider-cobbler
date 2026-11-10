package system_test

import (
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/cobbler/terraform-provider-cobbler/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "name", "foo-resource-system-basic"),
					resource.TestCheckResourceAttrPair("cobbler_system.foo", "profile", "cobbler_profile.foo", "uid"),
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system"),
				),
			},
			{
				ResourceName:                         "cobbler_system.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-basic",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccSystemResource_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourceChange1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system"),
				),
			},
			{
				Config: testAccSystemResourceChange2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "I'm a system again"),
				),
			},
			{
				ResourceName:                         "cobbler_system.foo",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "foo-resource-system-change",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccSystemResource_updatePersistsAfterRestart is a regression test for a bug where
// updating a cobbler_system resource (via modify_system XML-RPC calls) never called
// save_system afterward, so the edit only ever lived in cobblerd's in-memory handle and
// was silently lost on the next cobblerd restart/reload - even though the very next
// Terraform read (served from the same still-running cobblerd process) looked correct.
//
// Catching this requires actually restarting cobblerd between the edit and the
// verification read, which this repo's other acceptance tests don't do. Since that's not
// something the Terraform testing framework can drive directly, this test is opt-in:
// it's skipped unless COBBLER_RESTART_CMD is set to a shell command that restarts the
// cobblerd instance behind COBBLER_URL (e.g. "docker restart cobbler-dev").
func TestAccSystemResource_updatePersistsAfterRestart(t *testing.T) {
	restartCmd := os.Getenv("COBBLER_RESTART_CMD")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5)
			if restartCmd == "" {
				t.Skip("COBBLER_RESTART_CMD must be set (e.g. to a `docker restart <container>` command) " +
					"to exercise the cobblerd-restart persistence regression test")
			}
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourcePersistRestart1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "before restart"),
				),
			},
			{
				// This update is the one that must persist to disk via save_system.
				Config: testAccSystemResourcePersistRestart2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.foo", "comment", "after restart"),
				),
			},
			{
				// No Terraform-visible config change here: this step only restarts
				// cobblerd (forcing it to reload its collections from disk) and then
				// re-plans against the same config. If the prior update was never
				// saved to disk, cobblerd comes back with the old "before restart"
				// value, the refresh picks that up, and the plan is non-empty.
				PreConfig: func() {
					if err := exec.Command("sh", "-c", restartCmd).Run(); err != nil {
						t.Fatalf("failed to run COBBLER_RESTART_CMD %q: %v", restartCmd, err)
					}
					waitForCobblerReady(t)
				},
				Config:             testAccSystemResourcePersistRestart2,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// waitForCobblerReady polls COBBLER_URL until cobblerd accepts connections again after a
// restart, or fails the test if it doesn't come back within the timeout.
func waitForCobblerReady(t *testing.T) {
	t.Helper()
	url := os.Getenv("COBBLER_URL")
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		client := http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("cobblerd at %s did not become ready again after restart", url)
}

const testAccSystemResourcePersistRestart1 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name     = "foo-resource-system-persist-restart"
  profile  = cobbler_profile.foo.uid
  comment  = "before restart"
  power_id = "foo"
}
`

const testAccSystemResourcePersistRestart2 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name     = "foo-resource-system-persist-restart"
  profile  = cobbler_profile.foo.uid
  comment  = "after restart"
  power_id = "foo"
}
`

// TestAccSystemResource_importNotFound exercises the zero-match branch of the resource's
// post-import name-to-uid resolution fallback: importing an ID that doesn't match any System
// on the server must surface a clear "not found" diagnostic instead of a raw client error.
func TestAccSystemResource_importNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.SkipIfCobblerVersionLessThan(t, 3, 3, 5) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemResourceImportNotFound,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cobbler_system.notfound", "name", "foo-resource-system-import-not-found"),
				),
			},
			{
				ResourceName:  "cobbler_system.notfound",
				ImportState:   true,
				ImportStateId: "does-not-exist-system-resource",
				ExpectError:   regexp.MustCompile(`Cobbler System not found`),
			},
		},
	})
}

const testAccSystemResourceImportNotFound = testAccSystemDistroProfile + `
resource "cobbler_system" "notfound" {
  name    = "foo-resource-system-import-not-found"
  profile = cobbler_profile.foo.uid
}
`

// testAccSystemDistroProfile is the shared distro+profile config used by system tests.
const testAccSystemDistroProfile = `
resource "cobbler_distro" "foo" {
  name       = "foo-resource-system"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "foo" {
  name   = "foo-resource-profile"
  distro = cobbler_distro.foo.uid
}
`

const testAccSystemResourceBasic = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-basic"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"
  power_id     = "foo"
}
`

const testAccSystemResourceChange1 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-change"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"
  power_id     = "foo"
}
`

const testAccSystemResourceChange2 = testAccSystemDistroProfile + `
resource "cobbler_system" "foo" {
  name         = "foo-resource-system-change"
  profile      = cobbler_profile.foo.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system again"
  power_id     = "foo"
}
`
