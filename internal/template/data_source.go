package template

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &TemplateDataSource{}

type TemplateDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &TemplateDataSource{}
}

func (d *TemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *TemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to look up a Cobbler template (4.0.0+).",
		Attributes: map[string]dsschema.Attribute{
			"name":          dsschema.StringAttribute{Description: "The name of the template.", Required: true},
			"comment":       dsschema.StringAttribute{Description: "Free form text description.", Computed: true},
			"template_type": dsschema.StringAttribute{Description: "The template engine.", Computed: true},
			"uri": dsschema.SingleNestedAttribute{
				Description: "Where the template's content lives.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"schema": dsschema.StringAttribute{Computed: true},
					"path":   dsschema.StringAttribute{Computed: true},
				},
			},
			"tags":     dsschema.ListAttribute{Description: "Tags.", Computed: true, ElementType: types.StringType},
			"content":  dsschema.StringAttribute{Description: "The template body.", Computed: true},
			"built_in": dsschema.BoolAttribute{Description: "Whether the template is built into Cobbler.", Computed: true},
		},
	}
}

func (d *TemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpl, err := d.client.GetTemplate(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Template", err.Error())
		return
	}

	templateToDataSourceModel(ctx, *tpl, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
