package system

import "github.com/hashicorp/terraform-plugin-framework/types"

type systemResourceModel struct {
	Name              types.String `tfsdk:"name"`
	UID               types.String `tfsdk:"uid"`
	Autoinstall       types.String `tfsdk:"autoinstall"`
	Comment           types.String `tfsdk:"comment"`
	Gateway           types.String `tfsdk:"gateway"`
	Hostname          types.String `tfsdk:"hostname"`
	Image             types.String `tfsdk:"image"`
	IPv6DefaultDevice types.String `tfsdk:"ipv6_default_device"`
	NameServers       types.List   `tfsdk:"name_servers"`
	NameServersSearch types.List   `tfsdk:"name_servers_search"`
	NetbootEnabled    types.Bool   `tfsdk:"netboot_enabled"`
	NextServerV4      types.String `tfsdk:"next_server_v4"`
	NextServerV6      types.String `tfsdk:"next_server_v6"`
	PowerAddress      types.String `tfsdk:"power_address"`
	PowerID           types.String `tfsdk:"power_id"`
	PowerPass         types.String `tfsdk:"power_pass"`
	PowerType         types.String `tfsdk:"power_type"`
	PowerUser         types.String `tfsdk:"power_user"`
	Profile           types.String `tfsdk:"profile"`
	Proxy             types.String `tfsdk:"proxy"`
	Status            types.String `tfsdk:"status"`
	VirtDiskDriver    types.String `tfsdk:"virt_disk_driver"`
	VirtPath          types.String `tfsdk:"virt_path"`
	VirtPXEBoot       types.Bool   `tfsdk:"virt_pxe_boot"`
	VirtType          types.String `tfsdk:"virt_type"`
	// Inheritable:
	AutoinstallMeta   types.Object `tfsdk:"autoinstall_meta"`
	BootLoaders       types.Object `tfsdk:"boot_loaders"`
	EnableIPXE        types.Object `tfsdk:"enable_ipxe"`
	KernelOptions     types.Object `tfsdk:"kernel_options"`
	KernelOptionsPost types.Object `tfsdk:"kernel_options_post"`
	Owners            types.Object `tfsdk:"owners"`
	TemplateFiles     types.Map    `tfsdk:"template_files"`
	VirtAutoBoot      types.Object `tfsdk:"virt_auto_boot"`
	VirtCPUs          types.Object `tfsdk:"virt_cpus"`
	VirtFileSize      types.Object `tfsdk:"virt_file_size"`
	VirtRAM           types.Object `tfsdk:"virt_ram"`
}
