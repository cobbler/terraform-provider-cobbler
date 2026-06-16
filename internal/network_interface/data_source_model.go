package network_interface

import "github.com/hashicorp/terraform-plugin-framework/types"

type networkInterfaceDataSourceModel struct {
	Name            types.String `tfsdk:"name"`
	System          types.String `tfsdk:"system"`
	SystemName      types.String `tfsdk:"system_name"`
	Comment         types.String `tfsdk:"comment"`
	MacAddress      types.String `tfsdk:"mac_address"`
	InterfaceType   types.String `tfsdk:"interface_type"`
	InterfaceMaster types.String `tfsdk:"interface_master"`
	BondingOpts     types.String `tfsdk:"bonding_opts"`
	BridgeOpts      types.String `tfsdk:"bridge_opts"`
	ConnectedMode   types.Bool   `tfsdk:"connected_mode"`
	Management      types.Bool   `tfsdk:"management"`
	Static          types.Bool   `tfsdk:"static"`
	DHCPTag         types.String `tfsdk:"dhcp_tag"`
	IfGateway       types.String `tfsdk:"if_gateway"`
	MTU             types.String `tfsdk:"mtu"`
	VirtBridge      types.Object `tfsdk:"virt_bridge"`
	IPv4            types.Object `tfsdk:"ipv4"`
	IPv6            types.Object `tfsdk:"ipv6"`
	DNS             types.Object `tfsdk:"dns"`
}
