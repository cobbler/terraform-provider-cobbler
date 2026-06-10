package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IntAttrTypes = map[string]attr.Type{
	"value":     types.Int64Type,
	"inherited": types.BoolType,
}

func IntFrom(_ context.Context, v cobbler.Value[int], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(IntAttrTypes, map[string]attr.Value{
			"value":     types.Int64Null(),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	obj, d := types.ObjectValue(IntAttrTypes, map[string]attr.Value{
		"value":     types.Int64Value(int64(v.Data)),
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

func IntTo(_ context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[int] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[int]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[int]{IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[int]{IsInherited: true}
	}
	val, ok := attrs["value"].(types.Int64)
	if !ok || val.IsNull() || val.IsUnknown() {
		return cobbler.Value[int]{}
	}
	return cobbler.Value[int]{Data: int(val.ValueInt64())}
}
