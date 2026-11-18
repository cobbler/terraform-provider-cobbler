package system_group

import "github.com/hashicorp/terraform-plugin-framework/types"

type systemGroupResourceModel struct {
	Name    types.String `tfsdk:"name"`
	UID     types.String `tfsdk:"uid"`
	Comment types.String `tfsdk:"comment"`
	Parent  types.String `tfsdk:"parent"`
	Items   types.List   `tfsdk:"items"`
}
