package network_interface

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &NetworkInterfaceResource{}
var _ resource.ResourceWithImportState = &NetworkInterfaceResource{}

type NetworkInterfaceResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &NetworkInterfaceResource{}
}

func (r *NetworkInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (r *NetworkInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_network_interface` manages a network interface attached to a Cobbler system (Cobbler 4.0.0+).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The fully qualified interface name in the form `<ifname>@<system-name>` (e.g. `eth0@mybox`). Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system": schema.StringAttribute{
				Description: "The Cobbler UID of the parent system. Use `cobbler_system.foo.uid`. Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_name": schema.StringAttribute{
				Description: "The name of the parent system (computed echo from the server).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_address": schema.StringAttribute{
				Description: "The MAC address of the interface.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interface_type": schema.StringAttribute{
				Description: "Type of interface. One of: na, bond, bond_slave, bridge, bridge_slave, bonded_bridge_slave, infiniband.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("na"),
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("na", "bond", "bond_slave", "bridge", "bridge_slave", "bonded_bridge_slave", "infiniband"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interface_master": schema.StringAttribute{
				Description: "Name of the master interface when this interface is a slave.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bonding_opts": schema.StringAttribute{
				Description: "Options for bonded interfaces.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bridge_opts": schema.StringAttribute{
				Description: "Options for bridge interfaces.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connected_mode": schema.BoolAttribute{
				Description: "Whether InfiniBand connected-mode is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"management": schema.BoolAttribute{
				Description: "Whether this interface is a management interface.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"static": schema.BoolAttribute{
				Description: "Whether the interface is static (true) or DHCP (false).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dhcp_tag": schema.StringAttribute{
				Description: "DHCP tag.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"if_gateway": schema.StringAttribute{
				Description: "Per-interface gateway.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mtu": schema.StringAttribute{
				Description: "The interface MTU.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_bridge": schema.SingleNestedAttribute{
				Description: "The virtual bridge to attach to. Inheritable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"ipv4": schema.SingleNestedAttribute{
				Description: "Per-interface IPv4 configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The IPv4 address of the interface.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"netmask": schema.StringAttribute{
						Description: "The IPv4 netmask of the interface.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"gateway": schema.StringAttribute{
						Description: "The IPv4 gateway for the interface.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"static_routes": schema.ListAttribute{
						Description: "Static IPv4 routes for the interface.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "Per-interface IPv6 configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The IPv6 address of the interface.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"prefix": schema.StringAttribute{
						Description: "The IPv6 prefix length.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"mtu": schema.StringAttribute{
						Description: "The IPv6 MTU.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"default_gateway": schema.StringAttribute{
						Description: "The IPv6 default gateway.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"secondaries": schema.ListAttribute{
						Description: "IPv6 secondary addresses.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"static_routes": schema.ListAttribute{
						Description: "Static IPv6 routes for the interface.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dns": schema.SingleNestedAttribute{
				Description: "Per-interface DNS configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "DNS name.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"cnames": schema.ListAttribute{
						Description: "Canonical name records.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *NetworkInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*clientpkg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			"Expected *client.Config, got unexpected type.")
		return
	}
	r.client = cfg.CobblerClient
}

func (r *NetworkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data networkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemUid := data.System.ValueString()
	iface := modelToInterface(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	mu := lockForSystem(systemUid)
	mu.Lock()
	defer mu.Unlock()

	tflog.Debug(ctx, "Cobbler NetworkInterface: Create", map[string]interface{}{
		"name":   iface.Name,
		"system": systemUid,
	})

	created, err := r.client.CreateNetworkInterface(systemUid, iface)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler NetworkInterface", err.Error())
		return
	}

	interfaceToModel(ctx, *created, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data networkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.GetNetworkInterface(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Cobbler NetworkInterface", err.Error())
		return
	}

	interfaceToModel(ctx, *iface, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data networkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemUid := data.System.ValueString()
	iface := modelToInterface(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	mu := lockForSystem(systemUid)
	mu.Lock()
	defer mu.Unlock()

	tflog.Debug(ctx, "Cobbler NetworkInterface: Update", map[string]interface{}{"name": iface.Name})

	if err := r.client.UpdateNetworkInterface(&iface); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler NetworkInterface", err.Error())
		return
	}

	updated, err := r.client.GetNetworkInterface(iface.Name, false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler NetworkInterface after update", err.Error())
		return
	}

	interfaceToModel(ctx, *updated, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data networkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemUid := data.System.ValueString()
	mu := lockForSystem(systemUid)
	mu.Lock()
	defer mu.Unlock()

	tflog.Debug(ctx, "Cobbler NetworkInterface: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteNetworkInterface(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler NetworkInterface", err.Error())
	}
}

func (r *NetworkInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
