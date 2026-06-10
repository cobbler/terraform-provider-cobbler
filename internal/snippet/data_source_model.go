package snippet

import "github.com/hashicorp/terraform-plugin-framework/types"

type snippetDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Body types.String `tfsdk:"body"`
}
