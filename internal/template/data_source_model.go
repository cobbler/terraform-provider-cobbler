package template

import "github.com/hashicorp/terraform-plugin-framework/types"

type templateDataSourceModel struct {
	Name         types.String `tfsdk:"name"`
	Comment      types.String `tfsdk:"comment"`
	TemplateType types.String `tfsdk:"template_type"`
	URI          types.Object `tfsdk:"uri"`
	Tags         types.List   `tfsdk:"tags"`
	Content      types.String `tfsdk:"content"`
	BuiltIn      types.Bool   `tfsdk:"built_in"`
}
