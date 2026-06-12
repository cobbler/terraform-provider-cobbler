package template_file

import "github.com/hashicorp/terraform-plugin-framework/types"

type templateFileDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Body types.String `tfsdk:"body"`
}
