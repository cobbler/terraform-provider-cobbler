package acctest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/cobbler/terraform-provider-cobbler/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// CobblerApiClient is a pre-configured Cobbler client for use in acceptance tests.
var CobblerApiClient cobbler.Client

func init() {
	CobblerApiClient = cobbler.NewClient(&http.Client{}, cobbler.ClientConfig{
		URL:      os.Getenv("COBBLER_URL"),
		Username: os.Getenv("COBBLER_USERNAME"),
		Password: os.Getenv("COBBLER_PASSWORD"),
	})
	_, _ = CobblerApiClient.Login()
}

// ProtoV6ProviderFactories is a map of provider factories for use in acceptance tests.
// The key is the short provider type name. The testing framework builds the full
// address as "registry.terraform.io/<TF_ACC_PROVIDER_NAMESPACE>/cobbler".
// Set TF_ACC_PROVIDER_NAMESPACE=cobbler (see Makefile/CI) to avoid the legacy
// "-" wildcard namespace that OpenTofu rejects.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cobbler": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// PreCheck verifies that required environment variables are set for acceptance tests.
func PreCheck(t *testing.T) {
	t.Helper()
	for _, env := range []string{"COBBLER_USERNAME", "COBBLER_PASSWORD", "COBBLER_URL"} {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}

// SkipIfCobblerVersionLessThan skips the test if the running Cobbler server's version
// is older than the given major.minor.patch. Use this to gate tests on features that
// were introduced in a specific Cobbler release.
func SkipIfCobblerVersionLessThan(t *testing.T, major, minor, patch int) {
	t.Helper()
	ev, err := CobblerApiClient.ExtendedVersion()
	if err != nil {
		t.Skipf("could not determine Cobbler version: %v", err)
	}
	required := &cobbler.CobblerVersion{Major: major, Minor: minor, Patch: patch}
	tuple := ev.VersionTuple
	if len(tuple) < 3 {
		t.Skipf("unexpected Cobbler version tuple %v", tuple)
	}
	actual := &cobbler.CobblerVersion{Major: tuple[0], Minor: tuple[1], Patch: tuple[2]}
	if actual.LessThan(required) {
		t.Skipf("test requires Cobbler >= %s, server is %s", required, fmt.Sprintf("%d.%d.%d", tuple[0], tuple[1], tuple[2]))
	}
}
