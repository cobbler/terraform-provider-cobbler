package cobbler

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cobbler URL",
				DefaultFunc: envDefaultFunc("COBBLER_URL"),
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username for accessing Cobbler.",
				DefaultFunc: envDefaultFunc("COBBLER_USERNAME"),
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password for accessing Cobbler.",
				DefaultFunc: envDefaultFunc("COBBLER_PASSWORD"),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: envDefaultFunc("COBBLER_INSECURE"),
				Description: "Ignore SSL certificate warnings and errors.",
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("COBBLER_CACERT_FILE"),
				Description: "The path or contents of an SSL CA certificate",
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
		URL:        d.Get("URL").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func envDefaultFuncAllowMissing(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		v := os.Getenv(k)
		return v, nil
	}
}
