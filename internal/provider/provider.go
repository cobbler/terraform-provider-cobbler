package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &CobblerProvider{}

type CobblerProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CobblerProvider{version: version}
	}
}

func (p *CobblerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cobbler"
	resp.Version = p.version
}

func (p *CobblerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Cobbler provider is used to manage Cobbler resources.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The url to the Cobbler service. This can also be specified with the `COBBLER_URL` shell environment variable.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username to the Cobbler service. This can also be specified with the `COBBLER_USERNAME` shell environment variable.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to the Cobbler service. This can also be specified with the `COBBLER_PASSWORD` shell environment variable.",
				Required:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "If set to true, SSL certificate errors are ignored. This can also be specified with the `COBBLER_INSECURE` shell environment variable.",
				Optional:    true,
			},
			"cacert_file": schema.StringAttribute{
				Description: "The path or contents of an SSL CA certificate. This can also be specified with the `COBBLER_CACERT_FILE` shell environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *CobblerProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// TODO: implement in Step 2
}

func (p *CobblerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *CobblerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
