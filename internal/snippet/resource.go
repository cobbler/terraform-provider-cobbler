package snippet

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SnippetResource{}
var _ resource.ResourceWithImportState = &SnippetResource{}

type SnippetResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &SnippetResource{}
}

func (r *SnippetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snippet"
}

func (r *SnippetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_snippet` manages a snippet within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the snippet. This must be the name only, so without `/var/lib/cobbler/snippets`.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Description: "The body of the snippet. May also point to a file: `body = file(\"my_snippet\")`.",
				Required:    true,
			},
		},
	}
}

func (r *SnippetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*clientpkg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			"Expected *client.Config, got unexpected type.")
		return
	}
	r.client = cfg.CobblerClient
}

func (r *SnippetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data snippetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snippet := cobbler.Snippet{
		Name: data.Name.ValueString(),
		Body: data.Body.ValueString(),
	}

	tflog.Debug(ctx, "Cobbler Snippet: Create", map[string]interface{}{"name": snippet.Name})

	if err := r.client.CreateSnippet(snippet); err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler Snippet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnippetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data snippetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snippet, err := r.client.GetSnippet(data.Name.ValueString())
	if err != nil {
		// GetSnippet passes server errors through directly; Cobbler returns "not found"
		// in the error message for missing snippets (via read_autoinstall_snippet RPC).
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Snippet", err.Error())
		return
	}

	data.Name = types.StringValue(snippet.Name)
	data.Body = types.StringValue(snippet.Body)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnippetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data snippetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snippet := cobbler.Snippet{
		Name: data.Name.ValueString(),
		Body: data.Body.ValueString(),
	}

	tflog.Debug(ctx, "Cobbler Snippet: Update", map[string]interface{}{"name": snippet.Name})

	// The Cobbler API uses CreateSnippet for both create and update
	if err := r.client.CreateSnippet(snippet); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler Snippet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnippetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data snippetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Snippet: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteSnippet(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Snippet", err.Error())
		return
	}
}

func (r *SnippetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
