package client_test

import (
	"errors"
	"fmt"
	"testing"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestAddClientError_genericError(t *testing.T) {
	var diags diag.Diagnostics
	client.AddClientError(&diags, "Create failed", errors.New("something went wrong"))

	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Summary() != "Create failed" {
		t.Errorf("expected summary %q, got %q", "Create failed", diags[0].Summary())
	}
	if diags[0].Detail() != "something went wrong" {
		t.Errorf("expected detail %q, got %q", "something went wrong", diags[0].Detail())
	}
}

func TestAddClientError_inheritanceUnsupported(t *testing.T) {
	ie := &cobbler.InheritanceUnsupportedError{
		Field:          "enable_ipxe",
		ServerVersion:  cobbler.CobblerVersion{Major: 3, Minor: 3, Patch: 0},
		MinimumVersion: cobbler.CobblerVersion{Major: 3, Minor: 3, Patch: 5},
	}

	var diags diag.Diagnostics
	client.AddClientError(&diags, "ignored summary", ie)

	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	wantSummary := "Inheritance not supported: enable_ipxe"
	if diags[0].Summary() != wantSummary {
		t.Errorf("expected summary %q, got %q", wantSummary, diags[0].Summary())
	}
	if diags[0].Detail() != ie.Error() {
		t.Errorf("expected detail %q, got %q", ie.Error(), diags[0].Detail())
	}
}

func TestAddClientError_wrappedInheritanceUnsupported(t *testing.T) {
	ie := &cobbler.InheritanceUnsupportedError{
		Field:          "enable_ipxe",
		ServerVersion:  cobbler.CobblerVersion{Major: 3, Minor: 3, Patch: 0},
		MinimumVersion: cobbler.CobblerVersion{Major: 3, Minor: 3, Patch: 5},
	}
	wrapped := fmt.Errorf("outer context: %w", ie)

	var diags diag.Diagnostics
	client.AddClientError(&diags, "ignored summary", wrapped)

	if !diags.HasError() {
		t.Fatal("expected diagnostics to have an error")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	wantSummary := "Inheritance not supported: enable_ipxe"
	if diags[0].Summary() != wantSummary {
		t.Errorf("expected summary %q, got %q", wantSummary, diags[0].Summary())
	}
	if diags[0].Detail() != ie.Error() {
		t.Errorf("expected detail %q, got %q", ie.Error(), diags[0].Detail())
	}
}
