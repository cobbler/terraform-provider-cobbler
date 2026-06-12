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
				Description: "Parent image (if no profile is used).",
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
				Description: "Parent profile.",
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
			"interface": dsschema.MapNestedAttribute{
				Description: "A map of network interfaces, keyed by interface name (e.g. \"eth0\").",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"cnames": dsschema.ListAttribute{
							Description: "Canonical name records.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"dhcp_tag": dsschema.StringAttribute{
							Description: "DHCP tag.",
							Computed:    true,
						},
						"dns_name": dsschema.StringAttribute{
							Description: "DNS name.",
							Computed:    true,
						},
						"bonding_opts": dsschema.StringAttribute{
							Description: "Options for bonded interfaces.",
							Computed:    true,
						},
						"bridge_opts": dsschema.StringAttribute{
							Description: "Options for bridge interfaces.",
							Computed:    true,
						},
						"gateway": dsschema.StringAttribute{
							Description: "Per-interface gateway.",
							Computed:    true,
						},
						"interface_type": dsschema.StringAttribute{
							Description: "The type of interface.",
							Computed:    true,
						},
						"interface_master": dsschema.StringAttribute{
							Description: "The master interface when slave.",
							Computed:    true,
						},
						"ip_address": dsschema.StringAttribute{
							Description: "The IP address of the interface.",
							Computed:    true,
						},
						"ipv6_address": dsschema.StringAttribute{
							Description: "The IPv6 address of the interface.",
							Computed:    true,
						},
						"ipv6_secondaries": dsschema.ListAttribute{
							Description: "IPv6 secondaries.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"ipv6_mtu": dsschema.StringAttribute{
							Description: "The MTU of the IPv6 address.",
							Computed:    true,
						},
						"ipv6_static_routes": dsschema.ListAttribute{
							Description: "Static routes for the IPv6 interface.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"ipv6_default_gateway": dsschema.StringAttribute{
							Description: "The default gateway for the IPv6 address / interface.",
							Computed:    true,
						},
						"mac_address": dsschema.StringAttribute{
							Description: "The MAC address of the interface.",
							Computed:    true,
						},
						"management": dsschema.BoolAttribute{
							Description: "Whether this interface is a management interface.",
							Computed:    true,
						},
						"netmask": dsschema.StringAttribute{
							Description: "The IPv4 netmask of the interface.",
							Computed:    true,
						},
						"static": dsschema.BoolAttribute{
							Description: "Whether the interface should be static or DHCP.",
							Computed:    true,
						},
						"static_routes": dsschema.ListAttribute{
							Description: "Static routes for the interface.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"virt_bridge": dsschema.StringAttribute{
							Description: "The virtual bridge to attach to.",
							Computed:    true,
						},
					},
				},
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
			"boot_files": dsschema.SingleNestedAttribute{
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
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
			"fetchable_files": dsschema.SingleNestedAttribute{
				Description: "Templates for tftp or wget.",
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
			"mgmt_classes": dsschema.SingleNestedAttribute{
				Description: "For external configuration management.",
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
			"mgmt_parameters": dsschema.SingleNestedAttribute{
				Description: "Parameters which will be handed to your management application.",
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
			"template_files": dsschema.SingleNestedAttribute{
				Description: "File mappings for built-in config management.",
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

	ifaces, err := systemPtr.GetInterfaces()
	if err != nil {
		resp.Diagnostics.AddError("Error getting interfaces", err.Error())
		return
	}

	data.Name = types.StringValue(s.Name)
	data.Autoinstall = types.StringValue(s.Autoinstall)
	data.Comment = types.StringValue(s.Comment)
	data.Gateway = types.StringValue(s.Gateway)
	data.Hostname = types.StringValue(s.Hostname)
	data.Image = types.StringValue(s.Image)
	data.IPv6DefaultDevice = types.StringValue(s.IPv6DefaultDevice)
	data.NetbootEnabled = types.BoolValue(s.NetbootEnabled)
	data.NextServerV4 = types.StringValue(s.NextServerv4)
	data.NextServerV6 = types.StringValue(s.NextServerv6)
	data.PowerAddress = types.StringValue(s.PowerAddress)
	data.PowerID = types.StringValue(s.PowerID)
	data.PowerPass = types.StringValue(s.PowerPass)
	data.PowerType = types.StringValue(s.PowerType)
	data.PowerUser = types.StringValue(s.PowerUser)
	data.Profile = types.StringValue(s.Profile)
	data.Proxy = types.StringValue(s.Proxy)
	data.Status = types.StringValue(s.Status)
	data.VirtDiskDriver = types.StringValue(s.VirtDiskDriver)
	data.VirtPath = types.StringValue(s.VirtPath)
	data.VirtPXEBoot = types.BoolValue(s.VirtPXEBoot)
	data.VirtType = types.StringValue(s.VirtType)

	nameServersList, diag := types.ListValueFrom(ctx, types.StringType, s.NameServers)
	resp.Diagnostics.Append(diag...)
	data.NameServers = nameServersList

	nameServersSearchList, diag := types.ListValueFrom(ctx, types.StringType, s.NameServersSearch)
	resp.Diagnostics.Append(diag...)
	data.NameServersSearch = nameServersSearchList

	data.Interface = InterfaceMapFromAPI(ctx, ifaces, &resp.Diagnostics)

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, s.AutoinstallMeta, &resp.Diagnostics)
	data.BootFiles = inherit.StringMapFrom(ctx, s.BootFiles, &resp.Diagnostics)
	data.BootLoaders = inherit.StringListFrom(ctx, s.BootLoaders, &resp.Diagnostics)
	data.EnableIPXE = inherit.BoolFrom(ctx, s.EnableIPXE, &resp.Diagnostics)
	data.FetchableFiles = inherit.StringMapFrom(ctx, s.FetchableFiles, &resp.Diagnostics)
	data.KernelOptions = inherit.StringMapFrom(ctx, s.KernelOptions, &resp.Diagnostics)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, s.KernelOptionsPost, &resp.Diagnostics)
	data.MgmtClasses = inherit.StringListFrom(ctx, s.MgmtClasses, &resp.Diagnostics)
	data.MgmtParameters = inherit.StringMapFrom(ctx, s.MgmtParameters, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, s.Owners, &resp.Diagnostics)
	data.TemplateFiles = inherit.StringMapFrom(ctx, s.TemplateFiles, &resp.Diagnostics)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, s.VirtAutoBoot, &resp.Diagnostics)
	data.VirtCPUs = inherit.IntFrom(ctx, s.VirtCPUs, &resp.Diagnostics)
	data.VirtFileSize = inherit.Float64From(ctx, s.VirtFileSize, &resp.Diagnostics)
	data.VirtRAM = inherit.IntFrom(ctx, s.VirtRAM, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

