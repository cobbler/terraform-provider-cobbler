package inherit_test

import (
	"context"
	"testing"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringFrom_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[string]{IsInherited: true}

	obj := inherit.StringFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if obj.IsNull() || obj.IsUnknown() {
		t.Fatal("expected non-null, non-unknown object")
	}

	attrs := obj.Attributes()
	inherited, ok := attrs["inherited"].(types.Bool)
	if !ok {
		t.Fatal("inherited attr not a types.Bool")
	}
	if !inherited.ValueBool() {
		t.Error("expected inherited to be true")
	}
	val, ok := attrs["value"].(types.String)
	if !ok {
		t.Fatal("value attr not a types.String")
	}
	if !val.IsNull() {
		t.Error("expected value to be null when inherited")
	}
}

func TestStringFrom_Value(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[string]{Data: "hello", IsInherited: false}

	obj := inherit.StringFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if inherited.ValueBool() {
		t.Error("expected inherited to be false")
	}
	val := attrs["value"].(types.String)
	if val.ValueString() != "hello" {
		t.Errorf("expected value 'hello', got %q", val.ValueString())
	}
}

func TestStringRoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[string]{Data: "test-value", IsInherited: false}
	obj := inherit.StringFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringFrom diagnostics: %v", diags)
	}

	result := inherit.StringTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringTo diagnostics: %v", diags)
	}
	if result.Data != original.Data {
		t.Errorf("expected Data %q, got %q", original.Data, result.Data)
	}
	if result.IsInherited != original.IsInherited {
		t.Errorf("expected IsInherited %v, got %v", original.IsInherited, result.IsInherited)
	}
}

func TestStringRoundTrip_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[string]{IsInherited: true}
	obj := inherit.StringFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringFrom diagnostics: %v", diags)
	}

	result := inherit.StringTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringTo diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected IsInherited to be true")
	}
}

func TestStringTo_NullObject(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	obj := types.ObjectNull(inherit.StringAttrTypes)
	result := inherit.StringTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected null object to produce IsInherited=true")
	}
}

func TestBoolFrom_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[bool]{IsInherited: true}

	obj := inherit.BoolFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if !inherited.ValueBool() {
		t.Error("expected inherited to be true")
	}
	val := attrs["value"].(types.Bool)
	if !val.IsNull() {
		t.Error("expected value to be null when inherited")
	}
}

func TestBoolFrom_Value(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[bool]{Data: true, IsInherited: false}

	obj := inherit.BoolFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if inherited.ValueBool() {
		t.Error("expected inherited to be false")
	}
	val := attrs["value"].(types.Bool)
	if !val.ValueBool() {
		t.Error("expected value to be true")
	}
}

func TestBoolRoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[bool]{Data: true, IsInherited: false}
	obj := inherit.BoolFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("BoolFrom diagnostics: %v", diags)
	}

	result := inherit.BoolTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("BoolTo diagnostics: %v", diags)
	}
	if result.Data != original.Data {
		t.Errorf("expected Data %v, got %v", original.Data, result.Data)
	}
	if result.IsInherited != original.IsInherited {
		t.Errorf("expected IsInherited %v, got %v", original.IsInherited, result.IsInherited)
	}
}

func TestBoolRoundTrip_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[bool]{IsInherited: true}
	obj := inherit.BoolFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("BoolFrom diagnostics: %v", diags)
	}

	result := inherit.BoolTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("BoolTo diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected IsInherited to be true")
	}
}

func TestBoolTo_NullObject(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	obj := types.ObjectNull(inherit.BoolAttrTypes)
	result := inherit.BoolTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected null object to produce IsInherited=true")
	}
}

func TestStringListFrom_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[[]string]{IsInherited: true}

	obj := inherit.StringListFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if !inherited.ValueBool() {
		t.Error("expected inherited to be true")
	}
	val := attrs["value"].(types.List)
	if !val.IsNull() {
		t.Error("expected value to be null when inherited")
	}
}

func TestStringListFrom_Value(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[[]string]{Data: []string{"a", "b", "c"}, IsInherited: false}

	obj := inherit.StringListFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if inherited.ValueBool() {
		t.Error("expected inherited to be false")
	}
	listVal := attrs["value"].(types.List)
	if listVal.IsNull() || listVal.IsUnknown() {
		t.Fatal("expected non-null list")
	}
	var strs []string
	d := listVal.ElementsAs(ctx, &strs, false)
	if d.HasError() {
		t.Fatalf("ElementsAs diagnostics: %v", d)
	}
	if len(strs) != 3 || strs[0] != "a" || strs[1] != "b" || strs[2] != "c" {
		t.Errorf("expected [a b c], got %v", strs)
	}
}

func TestStringListRoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[[]string]{Data: []string{"x", "y"}, IsInherited: false}
	obj := inherit.StringListFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringListFrom diagnostics: %v", diags)
	}

	result := inherit.StringListTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringListTo diagnostics: %v", diags)
	}
	if len(result.Data) != 2 || result.Data[0] != "x" || result.Data[1] != "y" {
		t.Errorf("expected [x y], got %v", result.Data)
	}
	if result.IsInherited {
		t.Error("expected IsInherited to be false")
	}
}

func TestStringListRoundTrip_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[[]string]{IsInherited: true}
	obj := inherit.StringListFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringListFrom diagnostics: %v", diags)
	}

	result := inherit.StringListTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringListTo diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected IsInherited to be true")
	}
}

func TestStringListFrom_EmptySlice(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[[]string]{Data: []string{}, IsInherited: false}

	obj := inherit.StringListFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	listVal := attrs["value"].(types.List)
	if listVal.IsNull() {
		t.Error("expected non-null list for empty slice")
	}
	var strs []string
	d := listVal.ElementsAs(ctx, &strs, false)
	if d.HasError() {
		t.Fatalf("ElementsAs diagnostics: %v", d)
	}
	if len(strs) != 0 {
		t.Errorf("expected empty slice, got %v", strs)
	}
}

func TestStringMapFrom_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[map[string]interface{}]{IsInherited: true}

	obj := inherit.StringMapFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if !inherited.ValueBool() {
		t.Error("expected inherited to be true")
	}
	val := attrs["value"].(types.Map)
	if !val.IsNull() {
		t.Error("expected value to be null when inherited")
	}
}

func TestStringMapFrom_Value(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	v := cobbler.Value[map[string]interface{}]{
		Data:        map[string]interface{}{"key1": "val1", "key2": "val2"},
		IsInherited: false,
	}

	obj := inherit.StringMapFrom(ctx, v, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	inherited := attrs["inherited"].(types.Bool)
	if inherited.ValueBool() {
		t.Error("expected inherited to be false")
	}
	mapVal := attrs["value"].(types.Map)
	if mapVal.IsNull() || mapVal.IsUnknown() {
		t.Fatal("expected non-null map")
	}
	var m map[string]string
	d := mapVal.ElementsAs(ctx, &m, false)
	if d.HasError() {
		t.Fatalf("ElementsAs diagnostics: %v", d)
	}
	if m["key1"] != "val1" || m["key2"] != "val2" {
		t.Errorf("unexpected map contents: %v", m)
	}
}

func TestStringMapRoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[map[string]interface{}]{
		Data:        map[string]interface{}{"a": "1", "b": "2"},
		IsInherited: false,
	}
	obj := inherit.StringMapFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringMapFrom diagnostics: %v", diags)
	}

	result := inherit.StringMapTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringMapTo diagnostics: %v", diags)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Data))
	}
	if result.Data["a"] != "1" || result.Data["b"] != "2" {
		t.Errorf("unexpected map contents: %v", result.Data)
	}
	if result.IsInherited {
		t.Error("expected IsInherited to be false")
	}
}

func TestStringMapRoundTrip_Inherited(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	original := cobbler.Value[map[string]interface{}]{IsInherited: true}
	obj := inherit.StringMapFrom(ctx, original, &diags)
	if diags.HasError() {
		t.Fatalf("StringMapFrom diagnostics: %v", diags)
	}

	result := inherit.StringMapTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("StringMapTo diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected IsInherited to be true")
	}
}

func TestStringMapTo_NullObject(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	obj := types.ObjectNull(inherit.StringMapAttrTypes)
	result := inherit.StringMapTo(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsInherited {
		t.Error("expected null object to produce IsInherited=true")
	}
}
