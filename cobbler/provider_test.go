package cobbler

import (
	cobbler "github.com/cobbler/cobblerclient"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var cobblerApiClient cobbler.Client
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"cobbler": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func init() {
	cobblerApiClient = cobbler.NewClient(&http.Client{}, cobbler.ClientConfig{
		URL:      os.Getenv("COBBLER_URL"),
		Username: os.Getenv("COBBLER_USERNAME"),
		Password: os.Getenv("COBBLER_PASSWORD"),
	})
	_, _ = cobblerApiClient.Login()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"cobbler": func() (*schema.Provider, error) {
			return New("dev")(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccCobblerPreCheck(t *testing.T) {
	v := os.Getenv("COBBLER_USERNAME")
	if v == "" {
		t.Fatal("COBBLER_USERNAME must be set for acceptance tests.")
	}
	v = os.Getenv("COBBLER_PASSWORD")
	if v == "" {
		t.Fatal("COBBLER_PASSWORD must be set for acceptance tests.")
	}
	v = os.Getenv("COBBLER_URL")
	if v == "" {
		t.Fatal("COBBLER_URL must be set for acceptance tests.")
	}
}
