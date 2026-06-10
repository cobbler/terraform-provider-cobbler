package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringListAttrTypes is the attribute type map for a {value []string, inherited bool} nested object.
var StringListAttrTypes = map[string]attr.Type{
	"value":     types.ListType{ElemType: types.StringType},
	"inherited": types.BoolType,
}

// StringListFrom converts a cobbler Value[[]string] to a Terraform types.Object.
func StringListFrom(ctx context.Context, v cobbler.Value[[]string], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(StringListAttrTypes, map[string]attr.Value{
			"value":     types.ListNull(types.StringType),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	elems := make([]attr.Value, len(v.Data))
	for i, s := range v.Data {
		elems[i] = types.StringValue(s)
	}
	listVal, d := types.ListValue(types.StringType, elems)
	diags.Append(d...)
	obj, d := types.ObjectValue(StringListAttrTypes, map[string]attr.Value{
		"value":     listVal,
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

// StringListTo converts a Terraform types.Object back to a cobbler Value[[]string].
func StringListTo(ctx context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[[]string] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[[]string]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[[]string]{IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[[]string]{IsInherited: true}
	}
	listVal, ok := attrs["value"].(types.List)
	if !ok || listVal.IsNull() || listVal.IsUnknown() {
		return cobbler.Value[[]string]{Data: []string{}}
	}
	var strs []string
	d := listVal.ElementsAs(ctx, &strs, false)
	diags.Append(d...)
	if strs == nil {
		strs = []string{}
	}
	return cobbler.Value[[]string]{Data: strs}
}
