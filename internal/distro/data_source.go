package distro

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DistroDataSource{}

type DistroDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &DistroDataSource{}
}

func (d *DistroDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_distro"
}

func (d *DistroDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler distro.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the distro.",
				Required:    true,
			},
			"uid": schema.StringAttribute{
				Description: "Server-assigned UID for this distro. Use this as the value for `cobbler_profile.distro`.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the distro.",
				Computed:    true,
			},
			"breed": schema.StringAttribute{
				Description: "The \"breed\" of distribution.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"initrd": schema.StringAttribute{
				Description: "Absolute path to initrd on filesystem.",
				Computed:    true,
			},
			"kernel": schema.StringAttribute{
				Description: "Absolute path to kernel on filesystem.",
				Computed:    true,
			},
			"remote_boot_initrd": schema.StringAttribute{
				Description: "URL the bootloader directly retrieves and boots from.",
				Computed:    true,
			},
			"remote_boot_kernel": schema.StringAttribute{
				Description: "URL the bootloader directly retrieves and boots from.",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "The version of the distro.",
				Computed:    true,
			},
			"boot_loaders": schema.SingleNestedAttribute{
				Description: "Boot loaders.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options to use with the kernel.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"kernel_options_post": schema.SingleNestedAttribute{
				Description: "Post install kernel options.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
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
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *DistroDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DistroDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data distroDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	distroPtr, err := d.client.GetDistro(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Distro", err.Error())
		return
	}
	distro := *distroPtr

	data.Name = types.StringValue(distro.Name)
	data.UID = types.StringValue(distro.Uid)
	data.Arch = types.StringValue(distro.Arch)
	data.Breed = types.StringValue(distro.Breed)
	data.Comment = types.StringValue(distro.Comment)
	data.Initrd = types.StringValue(distro.Initrd)
	data.Kernel = types.StringValue(distro.Kernel)
	data.RemoteBootInitrd = types.StringValue(distro.RemoteBootInitrd)
	data.RemoteBootKernel = types.StringValue(distro.RemoteBootKernel)
	data.OSVersion = types.StringValue(distro.OSVersion)
	data.BootLoaders = inherit.StringListFrom(ctx, distro.BootLoaders, &resp.Diagnostics)
	data.KernelOptions = inherit.StringMapFrom(ctx, distro.KernelOptions, &resp.Diagnostics)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, distro.KernelOptionsPost, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, distro.Owners, &resp.Diagnostics)
	templateFiles, d2 := types.MapValueFrom(ctx, types.StringType, distro.TemplateFiles)
	resp.Diagnostics.Append(d2...)
	data.TemplateFiles = templateFiles

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
