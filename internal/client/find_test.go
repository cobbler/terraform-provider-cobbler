package client_test

import (
	"strings"
	"testing"

	"github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// fakeItem stands in for the cobblerclient item types (e.g. *cobbler.Distro) that
// ResolveUnique is used with; only the uid extracted via uidOf matters to the helper.
type fakeItem struct {
	Uid string
}

func uidOfFakeItem(i fakeItem) string {
	return i.Uid
}

func TestResolveUnique_singleMatch(t *testing.T) {
	var diags diag.Diagnostics
	matches := []fakeItem{{Uid: "abc123"}}

	uid, ok := client.ResolveUnique(&diags, "Distro", "d1", matches, uidOfFakeItem)

	if !ok {
		t.Fatal("expected ok to be true")
	}
	if uid != "abc123" {
		t.Errorf("expected uid %q, got %q", "abc123", uid)
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestResolveUnique_zeroMatches(t *testing.T) {
	var diags diag.Diagnostics
	matches := []fakeItem{}

	uid, ok := client.ResolveUnique(&diags, "Distro", "d1", matches, uidOfFakeItem)

	if ok {
		t.Fatal("expected ok to be false")
	}
	if uid != "" {
		t.Errorf("expected empty uid, got %q", uid)
	}
	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	wantSummary := "Cobbler Distro not found"
	if diags[0].Summary() != wantSummary {
		t.Errorf("expected summary %q, got %q", wantSummary, diags[0].Summary())
	}
	if !strings.Contains(diags[0].Detail(), "d1") {
		t.Errorf("expected detail to mention name %q, got %q", "d1", diags[0].Detail())
	}
	if !strings.Contains(diags[0].Detail(), "Distro") {
		t.Errorf("expected detail to mention type %q, got %q", "Distro", diags[0].Detail())
	}
}

func TestResolveUnique_multipleMatches(t *testing.T) {
	var diags diag.Diagnostics
	matches := []fakeItem{{Uid: "abc123"}, {Uid: "def456"}, {Uid: "ghi789"}}

	uid, ok := client.ResolveUnique(&diags, "NetworkInterface", "eth0", matches, uidOfFakeItem)

	if ok {
		t.Fatal("expected ok to be false")
	}
	if uid != "" {
		t.Errorf("expected empty uid, got %q", uid)
	}
	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	wantSummary := "Ambiguous Cobbler NetworkInterface name"
	if diags[0].Summary() != wantSummary {
		t.Errorf("expected summary %q, got %q", wantSummary, diags[0].Summary())
	}
	detail := diags[0].Detail()
	if !strings.Contains(detail, "3") {
		t.Errorf("expected detail to mention match count 3, got %q", detail)
	}
	if !strings.Contains(detail, "NetworkInterface") {
		t.Errorf("expected detail to mention type %q, got %q", "NetworkInterface", detail)
	}
	if !strings.Contains(detail, "eth0") {
		t.Errorf("expected detail to mention name %q, got %q", "eth0", detail)
	}
}

func TestResolveUnique_emptyTypeNameAndName(t *testing.T) {
	var diags diag.Diagnostics
	matches := []fakeItem{}

	uid, ok := client.ResolveUnique(&diags, "", "", matches, uidOfFakeItem)

	if ok {
		t.Fatal("expected ok to be false")
	}
	if uid != "" {
		t.Errorf("expected empty uid, got %q", uid)
	}
	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	wantSummary := "Cobbler  not found"
	if diags[0].Summary() != wantSummary {
		t.Errorf("expected summary %q, got %q", wantSummary, diags[0].Summary())
	}
	wantDetail := `No  was found with name "".`
	if diags[0].Detail() != wantDetail {
		t.Errorf("expected detail %q, got %q", wantDetail, diags[0].Detail())
	}
}

func TestResolveUnique_uidOfCalledOnlyOnSingleMatch(t *testing.T) {
	// uidOf must not be invoked when there isn't exactly one match, since callers may
	// rely on it only being safe to dereference fields on an actual match.
	calls := 0
	uidOf := func(i fakeItem) string {
		calls++
		return i.Uid
	}

	var diags diag.Diagnostics
	if _, ok := client.ResolveUnique(&diags, "Distro", "d1", []fakeItem{}, uidOf); ok {
		t.Fatal("expected ok to be false for zero matches")
	}
	if calls != 0 {
		t.Errorf("expected uidOf not to be called for zero matches, called %d times", calls)
	}

	diags = nil
	if _, ok := client.ResolveUnique(&diags, "Distro", "d1", []fakeItem{{Uid: "a"}, {Uid: "b"}}, uidOf); ok {
		t.Fatal("expected ok to be false for multiple matches")
	}
	if calls != 0 {
		t.Errorf("expected uidOf not to be called for multiple matches, called %d times", calls)
	}

	diags = nil
	if _, ok := client.ResolveUnique(&diags, "Distro", "d1", []fakeItem{{Uid: "a"}}, uidOf); !ok {
		t.Fatal("expected ok to be true for a single match")
	}
	if calls != 1 {
		t.Errorf("expected uidOf to be called exactly once for a single match, called %d times", calls)
	}
}
