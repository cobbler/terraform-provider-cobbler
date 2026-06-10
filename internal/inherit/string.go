package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringAttrTypes is the attribute type map for a {value string, inherited bool} nested object.
var StringAttrTypes = map[string]attr.Type{
	"value":     types.StringType,
	"inherited": types.BoolType,
}

// StringFrom converts a cobbler Value[string] to a Terraform types.Object.
func StringFrom(_ context.Context, v cobbler.Value[string], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(StringAttrTypes, map[string]attr.Value{
			"value":     types.StringNull(),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	obj, d := types.ObjectValue(StringAttrTypes, map[string]attr.Value{
		"value":     types.StringValue(v.Data),
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

// StringTo converts a Terraform types.Object back to a cobbler Value[string].
func StringTo(_ context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[string] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[string]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[string]{IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[string]{IsInherited: true}
	}
	val, ok := attrs["value"].(types.String)
	if !ok || val.IsNull() || val.IsUnknown() {
		return cobbler.Value[string]{}
	}
	return cobbler.Value[string]{Data: val.ValueString()}
}
