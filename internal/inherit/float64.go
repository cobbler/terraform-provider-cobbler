package inherit

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Float64AttrTypes = map[string]attr.Type{
	"value":     types.Float64Type,
	"inherited": types.BoolType,
}

func Float64From(_ context.Context, v cobbler.Value[float64], diags *diag.Diagnostics) types.Object {
	if v.IsInherited {
		obj, d := types.ObjectValue(Float64AttrTypes, map[string]attr.Value{
			"value":     types.Float64Null(),
			"inherited": types.BoolValue(true),
		})
		diags.Append(d...)
		return obj
	}
	obj, d := types.ObjectValue(Float64AttrTypes, map[string]attr.Value{
		"value":     types.Float64Value(v.Data),
		"inherited": types.BoolValue(false),
	})
	diags.Append(d...)
	return obj
}

func Float64To(_ context.Context, obj types.Object, diags *diag.Diagnostics) cobbler.Value[float64] {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.Value[float64]{IsInherited: true}
	}
	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok || inherited.IsNull() || inherited.IsUnknown() {
		return cobbler.Value[float64]{IsInherited: true}
	}
	if inherited.ValueBool() {
		return cobbler.Value[float64]{IsInherited: true}
	}
	val, ok := attrs["value"].(types.Float64)
	if !ok || val.IsNull() || val.IsUnknown() {
		return cobbler.Value[float64]{}
	}
	return cobbler.Value[float64]{Data: val.ValueFloat64()}
}
