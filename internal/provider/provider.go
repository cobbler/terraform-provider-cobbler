package provider

import (
	"context"
	"os"

	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/distro"
	"github.com/cobbler/terraform-provider-cobbler/internal/image"
	"github.com/cobbler/terraform-provider-cobbler/internal/profile"
	"github.com/cobbler/terraform-provider-cobbler/internal/repo"
	"github.com/cobbler/terraform-provider-cobbler/internal/snippet"
	"github.com/cobbler/terraform-provider-cobbler/internal/system"
	template_file "github.com/cobbler/terraform-provider-cobbler/internal/template_file"
	"github.com/cobbler/terraform-provider-cobbler/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username to the Cobbler service. This can also be specified with the `COBBLER_USERNAME` shell environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to the Cobbler service. This can also be specified with the `COBBLER_PASSWORD` shell environment variable.",
				Optional:    true,
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

// providerModel maps to the provider schema attributes.
type providerModel struct {
	URL        types.String `tfsdk:"url"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Insecure   types.Bool   `tfsdk:"insecure"`
	CACertFile types.String `tfsdk:"cacert_file"`
}

func (p *CobblerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := data.URL.ValueString()
	if url == "" {
		url = os.Getenv("COBBLER_URL")
	}
	username := data.Username.ValueString()
	if username == "" {
		username = os.Getenv("COBBLER_USERNAME")
	}
	password := data.Password.ValueString()
	if password == "" {
		password = os.Getenv("COBBLER_PASSWORD")
	}
	insecure := data.Insecure.ValueBool()
	if !insecure && os.Getenv("COBBLER_INSECURE") == "true" {
		insecure = true
	}
	cacertFile := data.CACertFile.ValueString()
	if cacertFile == "" {
		cacertFile = os.Getenv("COBBLER_CACERT_FILE")
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Cobbler URL",
			"The provider cannot create the Cobbler client because there is a missing or empty value for the Cobbler URL. "+
				"Set the url value in the configuration or use the COBBLER_URL environment variable.",
		)
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Cobbler Username",
			"The provider cannot create the Cobbler client because there is a missing or empty value for the Cobbler username. "+
				"Set the username value in the configuration or use the COBBLER_USERNAME environment variable.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Cobbler Password",
			"The provider cannot create the Cobbler client because there is a missing or empty value for the Cobbler password. "+
				"Set the password value in the configuration or use the COBBLER_PASSWORD environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	cfg := &clientpkg.Config{
		URL:        url,
		Username:   username,
		Password:   password,
		Insecure:   insecure,
		CACertFile: cacertFile,
	}

	if err := cfg.LoadAndValidate(util.Read); err != nil {
		resp.Diagnostics.AddError("Failed to configure Cobbler client", err.Error())
		return
	}

	resp.ResourceData = cfg
	resp.DataSourceData = cfg
}

func (p *CobblerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		distro.NewResource,
		image.NewResource,
		profile.NewResource,
		repo.NewResource,
		snippet.NewResource,
		system.NewResource,
		template_file.NewResource,
	}
}

func (p *CobblerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		distro.NewDataSource,
		image.NewDataSource,
		profile.NewDataSource,
		repo.NewDataSource,
		snippet.NewDataSource,
		system.NewDataSource,
		template_file.NewDataSource,
	}
}
