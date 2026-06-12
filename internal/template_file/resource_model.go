package template_file

import "github.com/hashicorp/terraform-plugin-framework/types"

type templateFileResourceModel struct {
	Name types.String `tfsdk:"name"`
	Body types.String `tfsdk:"body"`
}
