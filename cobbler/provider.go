package cobbler

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": {
					Description: "The url to the Cobbler service. This can also be specified with the `COBBLER_URL` shell environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("COBBLER_URL", nil),
				},

				"username": {
					Description: "The username to the Cobbler service. This can also be specified with the `COBBLER_USERNAME` shell environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("COBBLER_USERNAME", nil),
				},

				"password": {
					Description: "The password to the Cobbler service. This can also be specified with the `COBBLER_PASSWORD` shell environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("COBBLER_PASSWORD", nil),
				},

				"insecure": {
					Description: "The url to the Cobbler service. This can also be specified with the `COBBLER_URL` shell environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("COBBLER_INSECURE", nil),
				},

				"cacert_file": {
					Description: "The path or contents of an SSL CA certificate. This can also be specified with the `COBBLER_CACERT_FILE`shell environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
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
		}
		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := Config{
			CACertFile: d.Get("cacert_file").(string),
			Insecure:   d.Get("insecure").(bool),
			URL:        d.Get("url").(string),
			Username:   d.Get("username").(string),
			Password:   d.Get("password").(string),
		}

		if err := config.loadAndValidate(); err != nil {
			return nil, diag.FromErr(err)
		}

		return &config, nil
	}
}
