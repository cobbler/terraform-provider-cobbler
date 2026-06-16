package system_group

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SystemGroupDataSource{}

type SystemGroupDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &SystemGroupDataSource{}
}

func (d *SystemGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_group"
}

func (d *SystemGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to look up a Cobbler system group (4.0.0+).",
		Attributes: map[string]dsschema.Attribute{
			"name":    dsschema.StringAttribute{Description: "Name of the group.", Required: true},
			"comment": dsschema.StringAttribute{Description: "Free form text description.", Computed: true},
			"items":   dsschema.ListAttribute{Description: "Distro names in the group.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *SystemGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data systemGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	g, err := d.client.GetSystemGroup(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler SystemGroup", err.Error())
		return
	}

	groupToDataSourceModel(ctx, *g, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
