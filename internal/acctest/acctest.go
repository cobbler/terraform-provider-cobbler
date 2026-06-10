package acctest

import (
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
