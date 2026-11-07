package template

// Unit tests for the "refuse to manage built-in Cobbler templates" guards added to Create, Read,
// Update, and Delete. These run against a stubbed cobbler.Client (no live server / Docker image
// required), unlike the TestAcc* tests in resource_test.go which need a real beta8 Cobbler server.

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// tripwireHTTPClient fails the test the moment the provider issues any Cobbler API call. It
// proves a guard returns before touching the client at all, not merely before completing
// successfully.
type tripwireHTTPClient struct {
	t *testing.T
}

func (c *tripwireHTTPClient) Post(_, _ string, _ io.Reader) (*http.Response, error) {
	c.t.Fatal("unexpected Cobbler API call: the built-in template guard should reject before any RPC")
	return nil, nil
}

// cannedResponseHTTPClient always returns the same raw XML-RPC response body, regardless of the
// request - sufficient for Read's guard, which only needs one get_template round trip before it
// decides to reject.
type cannedResponseHTTPClient struct {
	body []byte
}

func (c *cannedResponseHTTPClient) Post(_, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(c.body)),
	}, nil
}

const builtInTemplateGetResponse = `<?xml version='1.0'?>
<methodResponse>
    <params>
        <param>
            <value>
                <struct>
                    <member>
                        <name>name</name>
                        <value><string>sample_end.ks</string></value>
                    </member>
                    <member>
                        <name>uid</name>
                        <value><string>00000000000000000000000000000001</string></value>
                    </member>
                    <member>
                        <name>built_in</name>
                        <value><boolean>1</boolean></value>
                    </member>
                </struct>
            </value>
        </param>
    </params>
</methodResponse>
`

func validURI(schemaVal string) types.Object {
	return types.ObjectValueMust(uriAttrTypes, map[string]attr.Value{
		"schema": types.StringValue(schemaVal),
		"path":   types.StringValue(""),
	})
}

func emptyTags() types.List {
	return types.ListValueMust(types.StringType, []attr.Value{})
}

func templateSchemaResponse(r *TemplateResource) resource.SchemaResponse {
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp
}

func newPlan(t *testing.T, r *TemplateResource, data templateResourceModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: templateSchemaResponse(r).Schema}
	if diags := plan.Set(context.Background(), &data); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags)
	}
	return plan
}

func newState(t *testing.T, r *TemplateResource, data templateResourceModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: templateSchemaResponse(r).Schema}
	if diags := state.Set(context.Background(), &data); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags)
	}
	return state
}

func TestTemplateResource_Create_RefusesImportlibSchema(t *testing.T) {
	r := &TemplateResource{client: cobbler.NewClient(&tripwireHTTPClient{t: t}, cobbler.ClientConfig{})}

	data := templateResourceModel{
		Name:         types.StringValue("some-builtin-template"),
		UID:          types.StringUnknown(),
		Comment:      types.StringValue(""),
		TemplateType: types.StringValue("jinja"),
		URI:          validURI("importlib"),
		Tags:         emptyTags(),
		Content:      types.StringValue(""),
		BuiltIn:      types.BoolUnknown(),
	}

	req := resource.CreateRequest{Plan: newPlan(t, r, data)}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic rejecting an importlib-schema template, got none")
	}
	if len(resp.Diagnostics) != 1 {
		t.Fatalf("expected exactly 1 diagnostic, got %d: %v", len(resp.Diagnostics), resp.Diagnostics)
	}
	if got := resp.Diagnostics[0].Summary(); got != "Cannot create a built-in Cobbler Template" {
		t.Errorf("unexpected diagnostic summary: %q", got)
	}
	if !strings.Contains(resp.Diagnostics[0].Detail(), "importlib") {
		t.Errorf("expected detail to mention \"importlib\", got %q", resp.Diagnostics[0].Detail())
	}
}

func TestTemplateResource_Read_RefusesBuiltInTemplate(t *testing.T) {
	r := &TemplateResource{client: cobbler.NewClient(&cannedResponseHTTPClient{body: []byte(builtInTemplateGetResponse)}, cobbler.ClientConfig{})}

	data := templateResourceModel{
		Name:         types.StringValue("placeholder"),
		UID:          types.StringValue("00000000000000000000000000000001"),
		Comment:      types.StringValue(""),
		TemplateType: types.StringValue("jinja"),
		URI:          validURI("file"),
		Tags:         emptyTags(),
		Content:      types.StringValue(""),
		BuiltIn:      types.BoolValue(false),
	}

	req := resource.ReadRequest{State: newState(t, r, data)}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a built-in template, got none")
	}
	if got := resp.Diagnostics[0].Summary(); got != "Cobbler Template is built-in" {
		t.Errorf("unexpected diagnostic summary: %q", got)
	}
	if !strings.Contains(resp.Diagnostics[0].Detail(), "sample_end.ks") {
		t.Errorf("expected detail to mention the template name, got %q", resp.Diagnostics[0].Detail())
	}
}

func TestTemplateResource_Update_RefusesBuiltInTemplate(t *testing.T) {
	r := &TemplateResource{client: cobbler.NewClient(&tripwireHTTPClient{t: t}, cobbler.ClientConfig{})}

	data := templateResourceModel{
		Name:         types.StringValue("a-builtin-template"),
		UID:          types.StringValue("00000000000000000000000000000001"),
		Comment:      types.StringValue(""),
		TemplateType: types.StringValue("jinja"),
		URI:          validURI("importlib"),
		Tags:         emptyTags(),
		Content:      types.StringValue(""),
		BuiltIn:      types.BoolValue(true),
	}

	req := resource.UpdateRequest{Plan: newPlan(t, r, data)}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic rejecting an update to a built-in template, got none")
	}
	if got := resp.Diagnostics[0].Summary(); got != "Cobbler Template is built-in" {
		t.Errorf("unexpected diagnostic summary: %q", got)
	}
	if !strings.Contains(resp.Diagnostics[0].Detail(), "cannot be modified") {
		t.Errorf("expected detail to say the template cannot be modified, got %q", resp.Diagnostics[0].Detail())
	}
}

func TestTemplateResource_Delete_RefusesBuiltInTemplate(t *testing.T) {
	r := &TemplateResource{client: cobbler.NewClient(&tripwireHTTPClient{t: t}, cobbler.ClientConfig{})}

	data := templateResourceModel{
		Name:         types.StringValue("a-builtin-template"),
		UID:          types.StringValue("00000000000000000000000000000001"),
		Comment:      types.StringValue(""),
		TemplateType: types.StringValue("jinja"),
		URI:          validURI("importlib"),
		Tags:         emptyTags(),
		Content:      types.StringValue(""),
		BuiltIn:      types.BoolValue(true),
	}

	req := resource.DeleteRequest{State: newState(t, r, data)}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic rejecting deletion of a built-in template, got none")
	}
	if got := resp.Diagnostics[0].Summary(); got != "Cobbler Template is built-in" {
		t.Errorf("unexpected diagnostic summary: %q", got)
	}
	if !strings.Contains(resp.Diagnostics[0].Detail(), "cannot be deleted") {
		t.Errorf("expected detail to say the template cannot be deleted, got %q", resp.Diagnostics[0].Detail())
	}
}
