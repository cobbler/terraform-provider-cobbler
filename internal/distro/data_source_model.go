package distro

import "github.com/hashicorp/terraform-plugin-framework/types"

type distroDataSourceModel struct {
	Name              types.String `tfsdk:"name"`
	Arch              types.String `tfsdk:"arch"`
	Breed             types.String `tfsdk:"breed"`
	Comment           types.String `tfsdk:"comment"`
	Initrd            types.String `tfsdk:"initrd"`
	Kernel            types.String `tfsdk:"kernel"`
	RemoteBootInitrd  types.String `tfsdk:"remote_boot_initrd"`
	RemoteBootKernel  types.String `tfsdk:"remote_boot_kernel"`
	OSVersion         types.String `tfsdk:"os_version"`
	BootFiles         types.Object `tfsdk:"boot_files"`
	BootLoaders       types.Object `tfsdk:"boot_loaders"`
	FetchableFiles    types.Object `tfsdk:"fetchable_files"`
	KernelOptions     types.Object `tfsdk:"kernel_options"`
	KernelOptionsPost types.Object `tfsdk:"kernel_options_post"`
	MgmtClasses       types.Object `tfsdk:"mgmt_classes"`
	Owners            types.Object `tfsdk:"owners"`
	TemplateFiles     types.Map    `tfsdk:"template_files"`
}
