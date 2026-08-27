package template

import "github.com/hashicorp/terraform-plugin-framework/types"

type templateResourceModel struct {
	Name         types.String `tfsdk:"name"`
	UID          types.String `tfsdk:"uid"`
	Comment      types.String `tfsdk:"comment"`
	TemplateType types.String `tfsdk:"template_type"`
	URI          types.Object `tfsdk:"uri"`
	Tags         types.List   `tfsdk:"tags"`
	Content      types.String `tfsdk:"content"`
	BuiltIn      types.Bool   `tfsdk:"built_in"`
}
