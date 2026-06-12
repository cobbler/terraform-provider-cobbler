package system

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// interfaceAttrTypes defines the types for each interface attribute
var interfaceAttrTypes = map[string]attr.Type{
	"cnames":               types.ListType{ElemType: types.StringType},
	"dhcp_tag":             types.StringType,
	"dns_name":             types.StringType,
	"bonding_opts":         types.StringType,
	"bridge_opts":          types.StringType,
	"gateway":              types.StringType,
	"interface_type":       types.StringType,
	"interface_master":     types.StringType,
	"ip_address":           types.StringType,
	"ipv6_address":         types.StringType,
	"ipv6_secondaries":     types.ListType{ElemType: types.StringType},
	"ipv6_mtu":             types.StringType,
	"ipv6_static_routes":   types.ListType{ElemType: types.StringType},
	"ipv6_default_gateway": types.StringType,
	"mac_address":          types.StringType,
	"management":           types.BoolType,
	"netmask":              types.StringType,
	"static":               types.BoolType,
	"static_routes":        types.ListType{ElemType: types.StringType},
	"virt_bridge":          types.StringType,
}

// InterfaceMapAttribute returns the schema.MapNestedAttribute for the interface map.
func InterfaceMapAttribute() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Description: "A map of network interfaces, keyed by interface name (e.g. \"eth0\").",
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"cnames": schema.ListAttribute{
					Description: "Canonical name records.",
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
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
				"dns_name": schema.StringAttribute{
					Description: "DNS name.",
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
				"gateway": schema.StringAttribute{
					Description: "Per-interface gateway.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"interface_type": schema.StringAttribute{
					Description: "The type of interface: NA, master, slave, bond, bond_slave, bridge, bridge_slave, bonded_bridge_slave, infiniband, bmc.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("na"),
					Validators: []validator.String{
						stringvalidator.OneOfCaseInsensitive("na", "master", "slave", "bond", "bond_slave", "bridge", "bridge_slave", "bonded_bridge_slave", "infiniband", "bmc"),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"interface_master": schema.StringAttribute{
					Description: "The master interface when slave.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"ip_address": schema.StringAttribute{
					Description: "The IP address of the interface.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6_address": schema.StringAttribute{
					Description: "The IPv6 address of the interface.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6_secondaries": schema.ListAttribute{
					Description: "IPv6 secondaries.",
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6_mtu": schema.StringAttribute{
					Description: "The MTU of the IPv6 address.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6_static_routes": schema.ListAttribute{
					Description: "Static routes for the IPv6 interface.",
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6_default_gateway": schema.StringAttribute{
					Description: "The default gateway for the IPv6 address / interface.",
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
				"management": schema.BoolAttribute{
					Description: "Whether this interface is a management interface.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
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
				"static": schema.BoolAttribute{
					Description: "Whether the interface should be static or DHCP.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"static_routes": schema.ListAttribute{
					Description: "Static routes for the interface.",
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				"virt_bridge": schema.StringAttribute{
					Description: "The virtual bridge to attach to.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}
}

// InterfaceMapFromAPI converts a cobbler.Interfaces map to a Terraform types.Map.
func InterfaceMapFromAPI(ctx context.Context, ifaces cobbler.Interfaces, diags *diag.Diagnostics) types.Map {
	if len(ifaces) == 0 {
		emptyMap, d := types.MapValue(types.ObjectType{AttrTypes: interfaceAttrTypes}, map[string]attr.Value{})
		diags.Append(d...)
		return emptyMap
	}

	elements := make(map[string]attr.Value, len(ifaces))
	for name, iface := range ifaces {
		obj, d := interfaceToObject(ctx, iface)
		diags.Append(d...)
		elements[name] = obj
	}
	m, d := types.MapValue(types.ObjectType{AttrTypes: interfaceAttrTypes}, elements)
	diags.Append(d...)
	return m
}

func interfaceToObject(ctx context.Context, iface cobbler.Interface) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	cnames, d := stringSliceToList(ctx, iface.CNAMEs)
	diags.Append(d...)
	ipv6Secondaries, d := stringSliceToList(ctx, iface.IPv6Secondaries)
	diags.Append(d...)
	ipv6StaticRoutes, d := stringSliceToList(ctx, iface.IPv6StaticRoutes)
	diags.Append(d...)
	staticRoutes, d := stringSliceToList(ctx, iface.StaticRoutes)
	diags.Append(d...)

	obj, d := types.ObjectValue(interfaceAttrTypes, map[string]attr.Value{
		"cnames":               cnames,
		"dhcp_tag":             types.StringValue(iface.DHCPTag),
		"dns_name":             types.StringValue(iface.DNSName),
		"bonding_opts":         types.StringValue(iface.BondingOpts),
		"bridge_opts":          types.StringValue(iface.BridgeOpts),
		"gateway":              types.StringValue(iface.Gateway),
		"interface_type":       types.StringValue(strings.ToLower(iface.InterfaceType)),
		"interface_master":     types.StringValue(iface.InterfaceMaster),
		"ip_address":           types.StringValue(iface.IPAddress),
		"ipv6_address":         types.StringValue(iface.IPv6Address),
		"ipv6_secondaries":     ipv6Secondaries,
		"ipv6_mtu":             types.StringValue(iface.IPv6MTU),
		"ipv6_static_routes":   ipv6StaticRoutes,
		"ipv6_default_gateway": types.StringValue(iface.IPv6DefaultGateway),
		"mac_address":          types.StringValue(iface.MACAddress),
		"management":           types.BoolValue(iface.Management),
		"netmask":              types.StringValue(iface.Netmask),
		"static":               types.BoolValue(iface.Static),
		"static_routes":        staticRoutes,
		"virt_bridge":          types.StringValue(iface.VirtBridge),
	})
	diags.Append(d...)
	return obj, diags
}

// InterfaceMapToAPI converts a Terraform types.Map back to cobbler.Interfaces.
func InterfaceMapToAPI(ctx context.Context, m types.Map, diags *diag.Diagnostics) cobbler.Interfaces {
	if m.IsNull() || m.IsUnknown() {
		return cobbler.Interfaces{}
	}

	var elements map[string]types.Object
	d := m.ElementsAs(ctx, &elements, false)
	diags.Append(d...)
	if diags.HasError() {
		return cobbler.Interfaces{}
	}

	ifaces := make(cobbler.Interfaces, len(elements))
	for name, obj := range elements {
		iface, d := objectToInterface(ctx, obj)
		diags.Append(d...)
		ifaces[name] = iface
	}
	return ifaces
}

func objectToInterface(ctx context.Context, obj types.Object) (cobbler.Interface, diag.Diagnostics) {
	var diags diag.Diagnostics
	attrs := obj.Attributes()

	cnames := listToStringSlice(ctx, attrs["cnames"].(types.List), &diags)
	ipv6Secondaries := listToStringSlice(ctx, attrs["ipv6_secondaries"].(types.List), &diags)
	ipv6StaticRoutes := listToStringSlice(ctx, attrs["ipv6_static_routes"].(types.List), &diags)
	staticRoutes := listToStringSlice(ctx, attrs["static_routes"].(types.List), &diags)

	interfaceType := attrs["interface_type"].(types.String).ValueString()
	if interfaceType == "" {
		interfaceType = "na"
	}

	var management bool
	if boolVal := attrs["management"].(types.Bool); !boolVal.IsNull() && !boolVal.IsUnknown() {
		management = boolVal.ValueBool()
	}

	var static bool
	if boolVal := attrs["static"].(types.Bool); !boolVal.IsNull() && !boolVal.IsUnknown() {
		static = boolVal.ValueBool()
	}

	return cobbler.Interface{
		CNAMEs:             cnames,
		DHCPTag:            attrs["dhcp_tag"].(types.String).ValueString(),
		DNSName:            attrs["dns_name"].(types.String).ValueString(),
		BondingOpts:        attrs["bonding_opts"].(types.String).ValueString(),
		BridgeOpts:         attrs["bridge_opts"].(types.String).ValueString(),
		Gateway:            attrs["gateway"].(types.String).ValueString(),
		InterfaceType:      interfaceType,
		InterfaceMaster:    attrs["interface_master"].(types.String).ValueString(),
		IPAddress:          attrs["ip_address"].(types.String).ValueString(),
		IPv6Address:        attrs["ipv6_address"].(types.String).ValueString(),
		IPv6Secondaries:    ipv6Secondaries,
		IPv6MTU:            attrs["ipv6_mtu"].(types.String).ValueString(),
		IPv6StaticRoutes:   ipv6StaticRoutes,
		IPv6DefaultGateway: attrs["ipv6_default_gateway"].(types.String).ValueString(),
		MACAddress:         attrs["mac_address"].(types.String).ValueString(),
		Management:         management,
		Netmask:            attrs["netmask"].(types.String).ValueString(),
		Static:             static,
		StaticRoutes:       staticRoutes,
		VirtBridge:         attrs["virt_bridge"].(types.String).ValueString(),
	}, diags
}

func stringSliceToList(ctx context.Context, ss []string) (types.List, diag.Diagnostics) {
	if ss == nil {
		ss = []string{}
	}
	return types.ListValueFrom(ctx, types.StringType, ss)
}

func listToStringSlice(ctx context.Context, l types.List, diags *diag.Diagnostics) []string {
	if l.IsNull() || l.IsUnknown() {
		return []string{}
	}
	var ss []string
	d := l.ElementsAs(ctx, &ss, false)
	diags.Append(d...)
	if ss == nil {
		ss = []string{}
	}
	return ss
}
