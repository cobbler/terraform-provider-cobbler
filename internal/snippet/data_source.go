package snippet

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SnippetDataSource{}

type SnippetDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &SnippetDataSource{}
}

func (d *SnippetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snippet"
}

func (d *SnippetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler snippet.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the snippet.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body of the snippet.",
				Computed:    true,
			},
		},
	}
}

func (d *SnippetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*clientpkg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			"Expected *client.Config, got unexpected type.")
		return
	}
	d.client = cfg.CobblerClient
}

func (d *SnippetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data snippetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snippet, err := d.client.GetSnippet(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Snippet", err.Error())
		return
	}

	data.Name = types.StringValue(snippet.Name)
	data.Body = types.StringValue(snippet.Body)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
