package template_file

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &TemplateFileDataSource{}

type TemplateFileDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &TemplateFileDataSource{}
}

func (d *TemplateFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_file"
}

func (d *TemplateFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler template file.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the template file.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body of the template file.",
				Computed:    true,
			},
		},
	}
}

func (d *TemplateFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TemplateFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templateFileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateFile, err := d.client.GetTemplateFile(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler TemplateFile", err.Error())
		return
	}

	data.Name = types.StringValue(templateFile.Name)
	data.Body = types.StringValue(templateFile.Body)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
