package profile

import "github.com/hashicorp/terraform-plugin-framework/types"

type profileDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Autoinstall    types.String `tfsdk:"autoinstall"`
	Comment        types.String `tfsdk:"comment"`
	DHCPTag        types.String `tfsdk:"dhcp_tag"`
	Distro         types.String `tfsdk:"distro"`
	NextServerV4   types.String `tfsdk:"next_server_v4"`
	NextServerV6   types.String `tfsdk:"next_server_v6"`
	Parent         types.String `tfsdk:"parent"`
	Proxy          types.String `tfsdk:"proxy"`
	Server         types.String `tfsdk:"server"`
	VirtBridge     types.String `tfsdk:"virt_bridge"`
	VirtCPUs       types.Int64  `tfsdk:"virt_cpus"`
	VirtDiskDriver types.String `tfsdk:"virt_disk_driver"`
	VirtPath       types.String `tfsdk:"virt_path"`
	VirtType       types.String `tfsdk:"virt_type"`
	Repos          types.List   `tfsdk:"repos"`
	// Inheritable:
	AutoinstallMeta   types.Object `tfsdk:"autoinstall_meta"`
	BootFiles         types.Object `tfsdk:"boot_files"`
	EnableIPXE        types.Object `tfsdk:"enable_ipxe"`
	EnableMenu        types.Object `tfsdk:"enable_menu"`
	FetchableFiles    types.Object `tfsdk:"fetchable_files"`
	KernelOptions     types.Object `tfsdk:"kernel_options"`
	KernelOptionsPost types.Object `tfsdk:"kernel_options_post"`
	NameServersSearch types.Object `tfsdk:"name_servers_search"`
	NameServers       types.Object `tfsdk:"name_servers"`
	Owners            types.Object `tfsdk:"owners"`
	TemplateFiles     types.Object `tfsdk:"template_files"`
	VirtAutoBoot      types.Object `tfsdk:"virt_auto_boot"`
	VirtFileSize      types.Object `tfsdk:"virt_file_size"`
	VirtRAM           types.Object `tfsdk:"virt_ram"`
}
