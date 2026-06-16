package client

import (
	"errors"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// AddClientError appends err to diags. For InheritanceUnsupportedError it uses a descriptive
// summary so the user can immediately identify which field to fix and set explicitly.
func AddClientError(diags *diag.Diagnostics, summary string, err error) {
	var ie *cobbler.InheritanceUnsupportedError
	if errors.As(err, &ie) {
		diags.AddError(
			"Inheritance not supported: "+ie.Field,
			ie.Error(),
		)
		return
	}
	diags.AddError(summary, err.Error())
}
