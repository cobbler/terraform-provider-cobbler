package network_interface

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &NetworkInterfaceDataSource{}

type NetworkInterfaceDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &NetworkInterfaceDataSource{}
}

func (d *NetworkInterfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (d *NetworkInterfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to look up a Cobbler network interface (Cobbler 4.0.0+).",
		Attributes: map[string]dsschema.Attribute{
			"name": dsschema.StringAttribute{
				Description: "The fully qualified interface name in the form `<ifname>@<system-name>`.",
				Required:    true,
			},
			"system":           dsschema.StringAttribute{Description: "The UID of the parent system.", Computed: true},
			"system_name":      dsschema.StringAttribute{Description: "The name of the parent system.", Computed: true},
			"comment":          dsschema.StringAttribute{Description: "Free form text description.", Computed: true},
			"mac_address":      dsschema.StringAttribute{Description: "The MAC address of the interface.", Computed: true},
			"interface_type":   dsschema.StringAttribute{Description: "Type of interface.", Computed: true},
			"interface_master": dsschema.StringAttribute{Description: "Master interface name when this is a slave.", Computed: true},
			"bonding_opts":     dsschema.StringAttribute{Description: "Options for bonded interfaces.", Computed: true},
			"bridge_opts":      dsschema.StringAttribute{Description: "Options for bridge interfaces.", Computed: true},
			"connected_mode":   dsschema.BoolAttribute{Description: "Whether InfiniBand connected-mode is enabled.", Computed: true},
			"management":       dsschema.BoolAttribute{Description: "Whether this is a management interface.", Computed: true},
			"static":           dsschema.BoolAttribute{Description: "Static (true) or DHCP (false).", Computed: true},
			"dhcp_tag":         dsschema.StringAttribute{Description: "DHCP tag.", Computed: true},
			"if_gateway":       dsschema.StringAttribute{Description: "Per-interface gateway.", Computed: true},
			"mtu":              dsschema.StringAttribute{Description: "The interface MTU.", Computed: true},
			"virt_bridge": dsschema.SingleNestedAttribute{
				Description: "The virtual bridge to attach to (inheritable).",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"value":     dsschema.StringAttribute{Computed: true},
					"inherited": dsschema.BoolAttribute{Computed: true},
				},
			},
			"ipv4": dsschema.SingleNestedAttribute{
				Description: "Per-interface IPv4 configuration.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"address":       dsschema.StringAttribute{Computed: true},
					"netmask":       dsschema.StringAttribute{Computed: true},
					"gateway":       dsschema.StringAttribute{Computed: true},
					"static_routes": dsschema.ListAttribute{Computed: true, ElementType: types.StringType},
				},
			},
			"ipv6": dsschema.SingleNestedAttribute{
				Description: "Per-interface IPv6 configuration.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"address":         dsschema.StringAttribute{Computed: true},
					"prefix":          dsschema.StringAttribute{Computed: true},
					"mtu":             dsschema.StringAttribute{Computed: true},
					"default_gateway": dsschema.StringAttribute{Computed: true},
					"secondaries":     dsschema.ListAttribute{Computed: true, ElementType: types.StringType},
					"static_routes":   dsschema.ListAttribute{Computed: true, ElementType: types.StringType},
				},
			},
			"dns": dsschema.SingleNestedAttribute{
				Description: "Per-interface DNS configuration.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"name":   dsschema.StringAttribute{Computed: true},
					"cnames": dsschema.ListAttribute{Computed: true, ElementType: types.StringType},
				},
			},
		},
	}
}

func (d *NetworkInterfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data networkInterfaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := d.client.GetNetworkInterface(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler NetworkInterface", err.Error())
		return
	}

	interfaceToDataSourceModel(ctx, *iface, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
