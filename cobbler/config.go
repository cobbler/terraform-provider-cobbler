package cobbler

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
		caCert, _, err := Read(c.CACertFile)
		if err != nil {
			return fmt.Errorf("error reading CA Cert: %s", err)
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
		return fmt.Errorf("failed to login: %s", err)
	}

	c.cobblerClient = client

	return nil
}
