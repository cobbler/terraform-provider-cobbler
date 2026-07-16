package network_interface

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ipv4AttrTypes = map[string]attr.Type{
	"address":       types.StringType,
	"netmask":       types.StringType,
	"gateway":       types.StringType,
	"static_routes": types.ListType{ElemType: types.StringType},
}

var ipv6AttrTypes = map[string]attr.Type{
	"address":         types.StringType,
	"prefix":          types.StringType,
	"mtu":             types.StringType,
	"default_gateway": types.StringType,
	"secondaries":     types.ListType{ElemType: types.StringType},
	"static_routes":   types.ListType{ElemType: types.StringType},
}

var dnsAttrTypes = map[string]attr.Type{
	"name":   types.StringType,
	"cnames": types.ListType{ElemType: types.StringType},
}

func parseInterfaceType(s string) cobbler.NetworkInterfaceType {
	switch s {
	case "bond":
		return cobbler.NetworkInterfaceTypeBond
	case "bond_slave":
		return cobbler.NetworkInterfaceTypeBondSlave
	case "bridge":
		return cobbler.NetworkInterfaceTypeBridge
	case "bridge_slave":
		return cobbler.NetworkInterfaceTypeBridgeSlave
	case "bonded_bridge_slave":
		return cobbler.NetworkInterfaceTypeBondedBridgeSlave
	case "infiniband":
		return cobbler.NetworkInterfaceTypeInfiniband
	default:
		return cobbler.NetworkInterfaceTypeNA
	}
}

func stringSliceToList(ctx context.Context, ss []string, diags *diag.Diagnostics) types.List {
	if ss == nil {
		ss = []string{}
	}
	l, d := types.ListValueFrom(ctx, types.StringType, ss)
	diags.Append(d...)
	return l
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

// ipv4FromAPI renders the nested "ipv4" object. There is no Gateway field on
// IPv4Option itself (see NetworkInterface.IfGateway) so the gateway value is
// passed in separately by the caller.
func ipv4FromAPI(ctx context.Context, o cobbler.IPv4Option, ifGateway string, diags *diag.Diagnostics) types.Object {
	obj, d := types.ObjectValue(ipv4AttrTypes, map[string]attr.Value{
		"address":       types.StringValue(o.Address),
		"netmask":       types.StringValue(o.Netmask),
		"gateway":       types.StringValue(ifGateway),
		"static_routes": stringSliceToList(ctx, o.StaticRoutes, diags),
	})
	diags.Append(d...)
	return obj
}

// ipv4ToAPI extracts the IPv4Option fields plus the "gateway" value, which the
// caller is responsible for assigning to NetworkInterface.IfGateway.
func ipv4ToAPI(ctx context.Context, obj types.Object, diags *diag.Diagnostics) (cobbler.IPv4Option, string) {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.IPv4Option{StaticRoutes: []string{}}, ""
	}
	attrs := obj.Attributes()
	address, _ := attrs["address"].(types.String)
	netmask, _ := attrs["netmask"].(types.String)
	gateway, _ := attrs["gateway"].(types.String)
	routes, _ := attrs["static_routes"].(types.List)
	return cobbler.IPv4Option{
		Address:      address.ValueString(),
		Netmask:      netmask.ValueString(),
		StaticRoutes: listToStringSlice(ctx, routes, diags),
	}, gateway.ValueString()
}

// ipv6FromAPI renders the nested "ipv6" object. There is no DefaultGateway
// field on IPv6Option itself (see NetworkInterface.Ipv6DefaultGateway) so the
// gateway value is passed in separately by the caller.
func ipv6FromAPI(ctx context.Context, o cobbler.IPv6Option, defaultGateway string, diags *diag.Diagnostics) types.Object {
	obj, d := types.ObjectValue(ipv6AttrTypes, map[string]attr.Value{
		"address":         types.StringValue(o.Address),
		"prefix":          types.StringValue(o.Prefix),
		"mtu":             types.StringValue(o.MTU),
		"default_gateway": types.StringValue(defaultGateway),
		"secondaries":     stringSliceToList(ctx, o.Secondaries, diags),
		"static_routes":   stringSliceToList(ctx, o.StaticRoutes, diags),
	})
	diags.Append(d...)
	return obj
}

// ipv6ToAPI extracts the IPv6Option fields plus the "default_gateway" value,
// which the caller is responsible for assigning to
// NetworkInterface.Ipv6DefaultGateway.
func ipv6ToAPI(ctx context.Context, obj types.Object, diags *diag.Diagnostics) (cobbler.IPv6Option, string) {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.IPv6Option{Secondaries: []string{}, StaticRoutes: []string{}}, ""
	}
	attrs := obj.Attributes()
	address, _ := attrs["address"].(types.String)
	prefix, _ := attrs["prefix"].(types.String)
	mtu, _ := attrs["mtu"].(types.String)
	gw, _ := attrs["default_gateway"].(types.String)
	secondaries, _ := attrs["secondaries"].(types.List)
	routes, _ := attrs["static_routes"].(types.List)
	return cobbler.IPv6Option{
		Address:      address.ValueString(),
		Prefix:       prefix.ValueString(),
		MTU:          mtu.ValueString(),
		Secondaries:  listToStringSlice(ctx, secondaries, diags),
		StaticRoutes: listToStringSlice(ctx, routes, diags),
	}, gw.ValueString()
}

func dnsFromAPI(ctx context.Context, o cobbler.DNSInterfaceOption, diags *diag.Diagnostics) types.Object {
	obj, d := types.ObjectValue(dnsAttrTypes, map[string]attr.Value{
		"name":   types.StringValue(o.Name),
		"cnames": stringSliceToList(ctx, o.CNames, diags),
	})
	diags.Append(d...)
	return obj
}

func dnsToAPI(ctx context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.DNSInterfaceOption {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.DNSInterfaceOption{CNames: []string{}}
	}
	attrs := obj.Attributes()
	name, _ := attrs["name"].(types.String)
	cnames, _ := attrs["cnames"].(types.List)
	return cobbler.DNSInterfaceOption{
		Name:   name.ValueString(),
		CNames: listToStringSlice(ctx, cnames, diags),
	}
}

// modelToInterface converts a Terraform model into a cobblerclient NetworkInterface.
func modelToInterface(ctx context.Context, data networkInterfaceResourceModel, diags *diag.Diagnostics) cobbler.NetworkInterface {
	iface := cobbler.NewNetworkInterface()
	iface.Name = data.Name.ValueString()
	iface.Comment = data.Comment.ValueString()
	iface.MacAddress = data.MacAddress.ValueString()
	iface.InterfaceType = parseInterfaceType(data.InterfaceType.ValueString())
	iface.InterfaceMaster = data.InterfaceMaster.ValueString()
	iface.BondingOpts = data.BondingOpts.ValueString()
	iface.BridgeOpts = data.BridgeOpts.ValueString()
	if !data.ConnectedMode.IsNull() && !data.ConnectedMode.IsUnknown() {
		iface.ConnectedMode = data.ConnectedMode.ValueBool()
	}
	if !data.Management.IsNull() && !data.Management.IsUnknown() {
		iface.Management = data.Management.ValueBool()
	}
	if !data.Static.IsNull() && !data.Static.IsUnknown() {
		iface.Static = data.Static.ValueBool()
	}
	iface.DHCPTag = data.DHCPTag.ValueString()
	iface.MTU = data.MTU.ValueString()
	iface.VirtBridge = inherit.StringTo(ctx, data.VirtBridge, diags)

	ipv4Opt, ipv4Gateway := ipv4ToAPI(ctx, data.IPv4, diags)
	iface.IPv4 = ipv4Opt
	iface.IfGateway = data.IfGateway.ValueString()
	if iface.IfGateway == "" {
		iface.IfGateway = ipv4Gateway
	}

	ipv6Opt, ipv6DefaultGateway := ipv6ToAPI(ctx, data.IPv6, diags)
	iface.IPv6 = ipv6Opt
	iface.Ipv6DefaultGateway = ipv6DefaultGateway

	iface.DNS = dnsToAPI(ctx, data.DNS, diags)
	return iface
}

// interfaceToModel populates a resource model from a NetworkInterface.
func interfaceToModel(ctx context.Context, iface cobbler.NetworkInterface, data *networkInterfaceResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(iface.Name)
	data.System = types.StringValue(iface.SystemUid)
	data.SystemName = types.StringValue(iface.SystemName)
	data.Comment = types.StringValue(iface.Comment)
	data.MacAddress = types.StringValue(iface.MacAddress)
	data.InterfaceType = types.StringValue(iface.InterfaceType)
	data.InterfaceMaster = types.StringValue(iface.InterfaceMaster)
	data.BondingOpts = types.StringValue(iface.BondingOpts)
	data.BridgeOpts = types.StringValue(iface.BridgeOpts)
	data.ConnectedMode = types.BoolValue(iface.ConnectedMode)
	data.Management = types.BoolValue(iface.Management)
	data.Static = types.BoolValue(iface.Static)
	data.DHCPTag = types.StringValue(iface.DHCPTag)
	data.IfGateway = types.StringValue(iface.IfGateway)
	data.MTU = types.StringValue(iface.MTU)
	data.VirtBridge = inherit.StringFrom(ctx, iface.VirtBridge, diags)
	data.IPv4 = ipv4FromAPI(ctx, iface.IPv4, iface.IfGateway, diags)
	data.IPv6 = ipv6FromAPI(ctx, iface.IPv6, iface.Ipv6DefaultGateway, diags)
	data.DNS = dnsFromAPI(ctx, iface.DNS, diags)
}

// interfaceToDataSourceModel populates a data source model from a NetworkInterface.
func interfaceToDataSourceModel(ctx context.Context, iface cobbler.NetworkInterface, data *networkInterfaceDataSourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(iface.Name)
	data.System = types.StringValue(iface.SystemUid)
	data.SystemName = types.StringValue(iface.SystemName)
	data.Comment = types.StringValue(iface.Comment)
	data.MacAddress = types.StringValue(iface.MacAddress)
	data.InterfaceType = types.StringValue(iface.InterfaceType)
	data.InterfaceMaster = types.StringValue(iface.InterfaceMaster)
	data.BondingOpts = types.StringValue(iface.BondingOpts)
	data.BridgeOpts = types.StringValue(iface.BridgeOpts)
	data.ConnectedMode = types.BoolValue(iface.ConnectedMode)
	data.Management = types.BoolValue(iface.Management)
	data.Static = types.BoolValue(iface.Static)
	data.DHCPTag = types.StringValue(iface.DHCPTag)
	data.IfGateway = types.StringValue(iface.IfGateway)
	data.MTU = types.StringValue(iface.MTU)
	data.VirtBridge = inherit.StringFrom(ctx, iface.VirtBridge, diags)
	data.IPv4 = ipv4FromAPI(ctx, iface.IPv4, iface.IfGateway, diags)
	data.IPv6 = ipv6FromAPI(ctx, iface.IPv6, iface.Ipv6DefaultGateway, diags)
	data.DNS = dnsFromAPI(ctx, iface.DNS, diags)
}
