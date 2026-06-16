package image

import "github.com/hashicorp/terraform-plugin-framework/types"

type imageResourceModel struct {
	Name              types.String `tfsdk:"name"`
	File              types.String `tfsdk:"file"`
	Arch              types.String `tfsdk:"arch"`
	Autoinstall       types.String `tfsdk:"autoinstall"`
	Breed             types.String `tfsdk:"breed"`
	Comment           types.String `tfsdk:"comment"`
	ImageType         types.String `tfsdk:"image_type"`
	OSVersion         types.String `tfsdk:"os_version"`
	BootLoaders       types.List   `tfsdk:"boot_loaders"`
	Menu              types.String `tfsdk:"menu"`
	VirtAutoBoot      types.Bool   `tfsdk:"virt_auto_boot"`
	VirtBridge        types.String `tfsdk:"virt_bridge"`
	VirtCpus          types.Int64  `tfsdk:"virt_cpus"`
	VirtDiskDriver    types.String `tfsdk:"virt_disk_driver"`
	VirtFileSize      types.Object `tfsdk:"virt_file_size"`
	VirtPath          types.String `tfsdk:"virt_path"`
	VirtRam           types.Object `tfsdk:"virt_ram"`
	VirtType          types.String `tfsdk:"virt_type"`
	KernelOptions     types.Object `tfsdk:"kernel_options"`
	KernelOptionsPost types.Object `tfsdk:"kernel_options_post"`
	FetchableFiles    types.Object `tfsdk:"fetchable_files"`
	BootFiles         types.Object `tfsdk:"boot_files"`
	Owners            types.Object `tfsdk:"owners"`
	TemplateFiles     types.Map    `tfsdk:"template_files"`
}
