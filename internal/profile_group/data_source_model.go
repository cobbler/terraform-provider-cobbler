package profile_group

import "github.com/hashicorp/terraform-plugin-framework/types"

type profileGroupDataSourceModel struct {
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Items   types.List   `tfsdk:"items"`
}
