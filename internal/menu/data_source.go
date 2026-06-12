package menu

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &MenuDataSource{}

type MenuDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &MenuDataSource{}
}

func (d *MenuDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_menu"
}

func (d *MenuDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler boot menu.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the menu.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"parent": schema.StringAttribute{
				Description: "The name of the parent menu.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name shown in the boot menu.",
				Computed:    true,
			},
			"autoinstall_meta": schema.SingleNestedAttribute{
				Description: "Autoinstall template metadata.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
			"fetchable_files": schema.SingleNestedAttribute{
				Description: "Templates for tftp or wget.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
			"boot_files": schema.SingleNestedAttribute{
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"mgmt_classes": schema.SingleNestedAttribute{
				Description: "Management classes for external config management.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
			"mgmt_parameters": schema.SingleNestedAttribute{
				Description: "Parameters for external management systems.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{Computed: true},
				},
			},
		},
	}
}

func (d *MenuDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MenuDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data menuDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	menuPtr, err := d.client.GetMenu(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Menu", err.Error())
		return
	}
	menu := *menuPtr

	data.Name = types.StringValue(menu.Name)
	data.Comment = types.StringValue(menu.Comment)
	data.Parent = types.StringValue(menu.Parent)
	data.DisplayName = types.StringValue(menu.DisplayName)
	data.AutoinstallMeta = inherit.StringMapFrom(ctx, menu.AutoinstallMeta, &resp.Diagnostics)
	data.FetchableFiles = inherit.StringMapFrom(ctx, menu.FetchableFiles, &resp.Diagnostics)
	data.BootFiles = inherit.StringMapFrom(ctx, menu.BootFiles, &resp.Diagnostics)
	data.MgmtClasses = inherit.StringListFrom(ctx, menu.MgmtClasses, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, menu.Owners, &resp.Diagnostics)
	data.MgmtParameters = inherit.StringMapFrom(ctx, menu.MgmtParameters, &resp.Diagnostics)

	templateFiles, d2 := types.MapValueFrom(ctx, types.StringType, menu.TemplateFiles.Data)
	resp.Diagnostics.Append(d2...)
	data.TemplateFiles = templateFiles

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
