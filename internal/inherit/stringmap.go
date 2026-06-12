package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringMapAttrTypes is the attribute type map for a {value map[string]string, inherited bool} nested object.
var StringMapAttrTypes = map[string]attr.Type{
	"value":     types.MapType{ElemType: types.StringType},
	"inherited": types.BoolType,
}

// StringMapFrom converts a cobbler Value[map[string]interface{}] to a Terraform types.Object.
// Note: the cobbler API uses map[string]interface{} but values are always strings.
func StringMapFrom(ctx context.Context, v cobbler.Value[map[string]interface{}], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(StringMapAttrTypes, map[string]attr.Value{
			"value":     types.MapNull(types.StringType),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	elems := make(map[string]attr.Value, len(v.Data))
	for k, val := range v.Data {
		if s, ok := val.(string); ok {
			elems[k] = types.StringValue(s)
		} else {
			elems[k] = types.StringValue("")
		}
	}
	mapVal, d := types.MapValue(types.StringType, elems)
	diags.Append(d...)
	obj, d := types.ObjectValue(StringMapAttrTypes, map[string]attr.Value{
		"value":     mapVal,
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

// StringMapTo converts a Terraform types.Object back to a cobbler Value[map[string]interface{}].
func StringMapTo(ctx context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[map[string]interface{}] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[map[string]interface{}]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[map[string]interface{}]{Data: make(map[string]interface{}), IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[map[string]interface{}]{Data: make(map[string]interface{}), IsInherited: true}
	}
	mapVal, ok := attrs["value"].(types.Map)
	if !ok || mapVal.IsNull() || mapVal.IsUnknown() {
		return cobbler.Value[map[string]interface{}]{Data: map[string]interface{}{}}
	}
	var m map[string]string
	d := mapVal.ElementsAs(ctx, &m, false)
	diags.Append(d...)
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = v
	}
	return cobbler.Value[map[string]interface{}]{Data: result}
}
