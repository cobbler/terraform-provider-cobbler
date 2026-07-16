package menu

import "github.com/hashicorp/terraform-plugin-framework/types"

type menuResourceModel struct {
	Name            types.String `tfsdk:"name"`
	UID             types.String `tfsdk:"uid"`
	Comment         types.String `tfsdk:"comment"`
	Parent          types.String `tfsdk:"parent"`
	DisplayName     types.String `tfsdk:"display_name"`
	AutoinstallMeta types.Object `tfsdk:"autoinstall_meta"`
	TemplateFiles   types.Map    `tfsdk:"template_files"`
	Owners          types.Object `tfsdk:"owners"`
}
