package cobbler

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"

	cobbler "github.com/wearespindle/cobblerclient"
)

// Config defines how to access the Cobbler API.
type Config struct {
	CACertFile string
	Insecure   bool
	URL        string
	Username   string
	Password   string

	cobblerClient cobbler.Client
}

func (c *Config) loadAndValidate() error {
	config := cobbler.ClientConfig{
		URL:      c.URL,
		Username: c.Username,
		Password: c.Password,
	}

	tlsConfig := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsConfig.RootCAs = caCertPool
	}

	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: transport}

	client := cobbler.NewClient(httpClient, config)
	_, err := client.Login()
	if err != nil {
		return fmt.Errorf("Failed to login: %s", err)
	}

	c.cobblerClient = client

	return nil
}
