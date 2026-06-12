package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	cobbler "github.com/cobbler/cobblerclient"
)

// Config defines how to access the Cobbler API.
type Config struct {
	CACertFile string
	Insecure   bool
	URL        string
	Username   string
	Password   string

	CobblerClient cobbler.Client
}

// LoadAndValidate configures the Cobbler client, performs TLS setup, and logs in.
// The readFile parameter is a function that reads a file path or returns the string as-is.
func (c *Config) LoadAndValidate(readFile func(string) (string, bool, error)) error {
	config := cobbler.ClientConfig{
		URL:      c.URL,
		Username: c.Username,
		Password: c.Password,
	}

	tlsConfig := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := readFile(c.CACertFile)
		if err != nil {
			return fmt.Errorf("error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsConfig.RootCAs = caCertPool
	}

	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true //nolint:gosec
	}

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: transport}

	client := cobbler.NewClient(httpClient, config)
	_, err := client.Login()
	if err != nil {
		return fmt.Errorf("failed to login: %s", err)
	}

	c.CobblerClient = client
	return nil
}
