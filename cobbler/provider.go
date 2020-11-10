package cobbler

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider does the talking to the Cobbler API.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cobbler URL",
				DefaultFunc: schema.EnvDefaultFunc("COBBLER_URL", nil),
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username for accessing Cobbler.",
				DefaultFunc: schema.EnvDefaultFunc("COBBLER_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password for accessing Cobbler.",
				DefaultFunc: schema.EnvDefaultFunc("COBBLER_PASSWORD", nil),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore SSL certificate warnings and errors.",
				DefaultFunc: schema.EnvDefaultFunc("COBBLER_INSECURE", nil),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path or contents of an SSL CA certificate",
				DefaultFunc: schema.EnvDefaultFunc("COBBLER_CACERT_FILE", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cobbler_distro":        resourceDistro(),
			"cobbler_template_file": resourceTemplateFile(),
			"cobbler_profile":       resourceProfile(),
			"cobbler_repo":          resourceRepo(),
			"cobbler_snippet":       resourceSnippet(),
			"cobbler_system":        resourceSystem(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		CACertFile: d.Get("cacert_file").(string),
		Insecure:   d.Get("insecure").(bool),
		URL:        d.Get("url").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
