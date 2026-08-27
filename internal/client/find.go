package client

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ResolveUnique extracts the uid from the single element of matches (as returned by one of
// the cobblerclient Find<Type> calls keyed on {"name": name}), adding a clear diagnostic
// error if there isn't exactly one match.
//
// Cobbler item names are not guaranteed unique across the whole server for every item type
// (e.g. NetworkInterface names are only unique per-system), so a name-based lookup can
// legitimately return zero or more than one result; callers must not silently pick one.
func ResolveUnique[T any](diags *diag.Diagnostics, typeName, name string, matches []T, uidOf func(T) string) (string, bool) {
	switch len(matches) {
	case 0:
		diags.AddError(
			fmt.Sprintf("Cobbler %s not found", typeName),
			fmt.Sprintf("No %s was found with name %q.", typeName, name),
		)
		return "", false
	case 1:
		return uidOf(matches[0]), true
	default:
		diags.AddError(
			fmt.Sprintf("Ambiguous Cobbler %s name", typeName),
			fmt.Sprintf("Found %d %s items with name %q; expected exactly one. Cobbler item names are not guaranteed to be globally unique for every item type.", len(matches), typeName, name),
		)
		return "", false
	}
}
