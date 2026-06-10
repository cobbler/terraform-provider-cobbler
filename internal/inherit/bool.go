package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BoolAttrTypes is the attribute type map for a {value bool, inherited bool} nested object.
var BoolAttrTypes = map[string]attr.Type{
	"value":     types.BoolType,
	"inherited": types.BoolType,
}

// BoolFrom converts a cobbler Value[bool] to a Terraform types.Object.
func BoolFrom(_ context.Context, v cobbler.Value[bool], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(BoolAttrTypes, map[string]attr.Value{
			"value":     types.BoolNull(),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	obj, d := types.ObjectValue(BoolAttrTypes, map[string]attr.Value{
		"value":     types.BoolValue(v.Data),
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

// BoolTo converts a Terraform types.Object back to a cobbler Value[bool].
func BoolTo(_ context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[bool] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[bool]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[bool]{IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[bool]{IsInherited: true}
	}
	val, ok := attrs["value"].(types.Bool)
	if !ok || val.IsNull() || val.IsUnknown() {
		return cobbler.Value[bool]{}
	}
	return cobbler.Value[bool]{Data: val.ValueBool()}
}
