package menu

import "github.com/hashicorp/terraform-plugin-framework/types"

type menuDataSourceModel struct {
	Name            types.String `tfsdk:"name"`
	Comment         types.String `tfsdk:"comment"`
	Parent          types.String `tfsdk:"parent"`
	DisplayName     types.String `tfsdk:"display_name"`
	AutoinstallMeta types.Object `tfsdk:"autoinstall_meta"`
	FetchableFiles  types.Object `tfsdk:"fetchable_files"`
	BootFiles       types.Object `tfsdk:"boot_files"`
	TemplateFiles   types.Map    `tfsdk:"template_files"`
	MgmtClasses     types.Object `tfsdk:"mgmt_classes"`
	Owners          types.Object `tfsdk:"owners"`
	MgmtParameters  types.Object `tfsdk:"mgmt_parameters"`
}
