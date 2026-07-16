package profile

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ProfileDataSource{}

type ProfileDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &ProfileDataSource{}
}

func (d *ProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (d *ProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler profile.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the profile.",
				Required:    true,
			},
			"uid": schema.StringAttribute{
				Description: "Server-assigned UID for this profile. Use this as the value for `cobbler_profile.parent` or `cobbler_system.profile`.",
				Computed:    true,
			},
			"autoinstall": schema.StringAttribute{
				Description: "Template remote kickstarts or preseeds.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"dhcp_tag": schema.StringAttribute{
				Description: "DHCP tag.",
				Computed:    true,
			},
			"distro": schema.StringAttribute{
				Description: "The Cobbler UID of the parent distribution.",
				Computed:    true,
			},
			"next_server_v4": schema.StringAttribute{
				Description: "The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded.",
				Computed:    true,
			},
			"next_server_v6": schema.StringAttribute{
				Description: "The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded.",
				Computed:    true,
			},
			"parent": schema.StringAttribute{
				Description: "The Cobbler UID of the parent profile this profile inherits settings from.",
				Computed:    true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy URL.",
				Computed:    true,
			},
			"server": schema.StringAttribute{
				Description: "The server-override for the profile.",
				Computed:    true,
			},
			"virt_bridge": schema.StringAttribute{
				Description: "The bridge for virtual machines.",
				Computed:    true,
			},
			"virt_cpus": schema.Int64Attribute{
				Description: "The number of virtual CPUs.",
				Computed:    true,
			},
			"virt_disk_driver": schema.StringAttribute{
				Description: "The virtual machine disk driver.",
				Computed:    true,
			},
			"virt_path": schema.StringAttribute{
				Description: "The virtual machine path.",
				Computed:    true,
			},
			"virt_type": schema.StringAttribute{
				Description: "The type of virtual machine.",
				Computed:    true,
			},
			"repos": schema.ListAttribute{
				Description: "Repos to auto-assign to this profile.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"autoinstall_meta": schema.SingleNestedAttribute{
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"enable_ipxe": schema.SingleNestedAttribute{
				Description: "Use iPXE instead of PXELINUX for advanced booting options.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"enable_menu": schema.SingleNestedAttribute{
				Description: "Enable a boot menu.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options for the profile.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"kernel_options_post": schema.SingleNestedAttribute{
				Description: "Post install kernel options.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"name_servers_search": schema.ListAttribute{
				Description: "Name server search settings. Not inheritable.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"name_servers": schema.SingleNestedAttribute{
				Description: "Name servers.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management. Not inheritable.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"virt_auto_boot": schema.SingleNestedAttribute{
				Description: "Auto boot virtual machines.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"virt_file_size": schema.SingleNestedAttribute{
				Description: "The virtual machine file size.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"virt_ram": schema.SingleNestedAttribute{
				Description: "The amount of RAM for the virtual machine.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *ProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data profileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profilePtr, err := d.client.GetProfile(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Profile", err.Error())
		return
	}
	p := *profilePtr

	data.Name = types.StringValue(p.Name)
	data.UID = types.StringValue(p.Uid)
	data.Autoinstall = types.StringValue(p.Autoinstall)
	data.Comment = types.StringValue(p.Comment)
	data.DHCPTag = types.StringValue(p.DHCPTag)
	data.Distro = types.StringValue(p.Distro)
	data.NextServerV4 = types.StringValue(p.TFTP.NextServerV4)
	data.NextServerV6 = types.StringValue(p.TFTP.NextServerV6)
	data.Parent = types.StringValue(p.Parent)
	data.Proxy = types.StringValue(p.Proxy)
	data.Server = types.StringValue(p.Server)
	data.VirtBridge = types.StringValue(p.VirtBridge)
	data.VirtCPUs = types.Int64Value(int64(p.Virt.Cpus.Data))
	data.VirtDiskDriver = types.StringValue(p.Virt.DiskDriver)
	data.VirtPath = types.StringValue(p.Virt.Path)
	data.VirtType = types.StringValue(p.Virt.Type)

	repoList, diag := types.ListValueFrom(ctx, types.StringType, p.Repos)
	resp.Diagnostics.Append(diag...)
	data.Repos = repoList

	nameServersSearch, diag2 := types.ListValueFrom(ctx, types.StringType, p.DNS.NameServersSearch)
	resp.Diagnostics.Append(diag2...)
	data.NameServersSearch = nameServersSearch

	templateFiles, diag3 := types.MapValueFrom(ctx, types.StringType, p.TemplateFiles)
	resp.Diagnostics.Append(diag3...)
	data.TemplateFiles = templateFiles

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, p.AutoinstallMeta, &resp.Diagnostics)
	data.EnableIPXE = inherit.BoolFrom(ctx, p.EnableIPXE, &resp.Diagnostics)
	data.EnableMenu = inherit.BoolFrom(ctx, p.EnableMenu, &resp.Diagnostics)
	data.KernelOptions = inherit.StringMapFrom(ctx, p.KernelOptions, &resp.Diagnostics)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, p.KernelOptionsPost, &resp.Diagnostics)
	data.NameServers = inherit.StringListFrom(ctx, p.DNS.NameServers, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, p.Owners, &resp.Diagnostics)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, p.Virt.AutoBoot, &resp.Diagnostics)
	data.VirtFileSize = inherit.Float64From(ctx, p.Virt.FileSize, &resp.Diagnostics)
	data.VirtRAM = inherit.IntFrom(ctx, p.Virt.Ram, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
