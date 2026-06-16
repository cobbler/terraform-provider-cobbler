package system

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SystemDataSource{}

type SystemDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &SystemDataSource{}
}

func (d *SystemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (d *SystemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to get the details of a Cobbler system.",
		Attributes: map[string]dsschema.Attribute{
			"name": dsschema.StringAttribute{
				Description: "The name of the system.",
				Required:    true,
			},
			"uid": dsschema.StringAttribute{
				Description: "Server-assigned UID for this system.",
				Computed:    true,
			},
			"autoinstall": dsschema.StringAttribute{
				Description: "Template remote kickstarts or preseeds.",
				Computed:    true,
			},
			"comment": dsschema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"gateway": dsschema.StringAttribute{
				Description: "Network gateway.",
				Computed:    true,
			},
			"hostname": dsschema.StringAttribute{
				Description: "Hostname of the system.",
				Computed:    true,
			},
			"image": dsschema.StringAttribute{
				Description: "The Cobbler UID of the parent image (if no profile is used).",
				Computed:    true,
			},
			"ipv6_default_device": dsschema.StringAttribute{
				Description: "IPv6 default device.",
				Computed:    true,
			},
			"name_servers": dsschema.ListAttribute{
				Description: "Name servers.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"name_servers_search": dsschema.ListAttribute{
				Description: "Name server search settings.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"netboot_enabled": dsschema.BoolAttribute{
				Description: "(Re)install this machine at next boot.",
				Computed:    true,
			},
			"next_server_v4": dsschema.StringAttribute{
				Description: "The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded.",
				Computed:    true,
			},
			"next_server_v6": dsschema.StringAttribute{
				Description: "The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded.",
				Computed:    true,
			},
			"power_address": dsschema.StringAttribute{
				Description: "Power management address.",
				Computed:    true,
			},
			"power_id": dsschema.StringAttribute{
				Description: "Usually a plug number or blade name if power type requires it.",
				Computed:    true,
			},
			"power_pass": dsschema.StringAttribute{
				Description: "Power management password.",
				Computed:    true,
				Sensitive:   true,
			},
			"power_type": dsschema.StringAttribute{
				Description: "Power management type.",
				Computed:    true,
			},
			"power_user": dsschema.StringAttribute{
				Description: "Power management user.",
				Computed:    true,
			},
			"profile": dsschema.StringAttribute{
				Description: "The Cobbler UID of the parent profile.",
				Computed:    true,
			},
			"proxy": dsschema.StringAttribute{
				Description: "Proxy URL.",
				Computed:    true,
			},
			"status": dsschema.StringAttribute{
				Description: "System status (development, testing, acceptance, production).",
				Computed:    true,
			},
			"virt_disk_driver": dsschema.StringAttribute{
				Description: "The virtual machine disk driver.",
				Computed:    true,
			},
			"virt_path": dsschema.StringAttribute{
				Description: "The virtual machine path.",
				Computed:    true,
			},
			"virt_pxe_boot": dsschema.BoolAttribute{
				Description: "Use PXE to build this virtual machine.",
				Computed:    true,
			},
			"virt_type": dsschema.StringAttribute{
				Description: "The type of virtual machine.",
				Computed:    true,
			},
			"autoinstall_meta": dsschema.SingleNestedAttribute{
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"boot_loaders": dsschema.SingleNestedAttribute{
				Description: "Must be either `grub`, `pxe`, or `ipxe`.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.ListAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"enable_ipxe": dsschema.SingleNestedAttribute{
				Description: "Use iPXE instead of PXELINUX for advanced booting options.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.BoolAttribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"kernel_options": dsschema.SingleNestedAttribute{
				Description: "Kernel options for the system.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"kernel_options_post": dsschema.SingleNestedAttribute{
				Description: "Post install kernel options.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.MapAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"owners": dsschema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.ListAttribute{
						Description: "The value.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"template_files": dsschema.MapAttribute{
				Description: "File mappings for built-in config management. Not inheritable.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"virt_auto_boot": dsschema.SingleNestedAttribute{
				Description: "Auto boot virtual machines.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.BoolAttribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"virt_cpus": dsschema.SingleNestedAttribute{
				Description: "The number of virtual CPUs.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.Int64Attribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"virt_file_size": dsschema.SingleNestedAttribute{
				Description: "The virtual machine file size.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.Float64Attribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
			"virt_ram": dsschema.SingleNestedAttribute{
				Description: "The amount of RAM for the virtual machine.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value": dsschema.Int64Attribute{
						Description: "The value.",
						Computed:    true,
					},
					"inherited": dsschema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *SystemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data systemDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemPtr, err := d.client.GetSystem(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler System", err.Error())
		return
	}
	s := *systemPtr

	data.Name = types.StringValue(s.Name)
	data.UID = types.StringValue(s.Uid)
	data.Autoinstall = types.StringValue(s.Autoinstall)
	data.Comment = types.StringValue(s.Comment)
	data.Gateway = types.StringValue(s.Gateway)
	data.Hostname = types.StringValue(s.Hostname)
	data.Image = types.StringValue(s.Image)
	data.IPv6DefaultDevice = types.StringValue(s.IPv6DefaultDevice)
	data.NetbootEnabled = types.BoolValue(s.NetbootEnabled)
	data.NextServerV4 = types.StringValue(s.TFTP.NextServerV4)
	data.NextServerV6 = types.StringValue(s.TFTP.NextServerV6)
	data.PowerAddress = types.StringValue(s.Power.Address)
	data.PowerID = types.StringValue(s.Power.ID)
	data.PowerPass = types.StringValue(s.Power.Password)
	data.PowerType = types.StringValue(s.Power.Type)
	data.PowerUser = types.StringValue(s.Power.User)
	data.Profile = types.StringValue(s.Profile)
	data.Proxy = types.StringValue(s.Proxy)
	data.Status = types.StringValue(s.Status)
	data.VirtDiskDriver = types.StringValue(s.Virt.DiskDriver)
	data.VirtPath = types.StringValue(s.Virt.Path)
	data.VirtPXEBoot = types.BoolValue(s.VirtPXEBoot)
	data.VirtType = types.StringValue(s.Virt.Type)

	nameServersList, diag := types.ListValueFrom(ctx, types.StringType, s.DNS.NameServers.Data)
	resp.Diagnostics.Append(diag...)
	data.NameServers = nameServersList

	nameServersSearchList, diag := types.ListValueFrom(ctx, types.StringType, s.DNS.NameServersSearch)
	resp.Diagnostics.Append(diag...)
	data.NameServersSearch = nameServersSearchList

	templateFiles, diag := types.MapValueFrom(ctx, types.StringType, s.TemplateFiles)
	resp.Diagnostics.Append(diag...)
	data.TemplateFiles = templateFiles

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, s.AutoinstallMeta, &resp.Diagnostics)
	data.BootLoaders = inherit.StringListFrom(ctx, s.BootLoaders, &resp.Diagnostics)
	data.EnableIPXE = inherit.BoolFrom(ctx, s.EnableIPXE, &resp.Diagnostics)
	data.KernelOptions = inherit.StringMapFrom(ctx, s.KernelOptions, &resp.Diagnostics)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, s.KernelOptionsPost, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, s.Owners, &resp.Diagnostics)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, s.Virt.AutoBoot, &resp.Diagnostics)
	data.VirtCPUs = inherit.IntFrom(ctx, s.Virt.Cpus, &resp.Diagnostics)
	data.VirtFileSize = inherit.Float64From(ctx, s.Virt.FileSize, &resp.Diagnostics)
	data.VirtRAM = inherit.IntFrom(ctx, s.Virt.Ram, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
