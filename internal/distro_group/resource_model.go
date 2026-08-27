package distro_group

import "github.com/hashicorp/terraform-plugin-framework/types"

type distroGroupResourceModel struct {
	Name    types.String `tfsdk:"name"`
	UID     types.String `tfsdk:"uid"`
	Comment types.String `tfsdk:"comment"`
	Items   types.List   `tfsdk:"items"`
}
