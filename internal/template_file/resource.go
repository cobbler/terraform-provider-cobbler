package template_file

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

var _ resource.Resource = &TemplateFileResource{}
var _ resource.ResourceWithImportState = &TemplateFileResource{}

type TemplateFileResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &TemplateFileResource{}
}

func (r *TemplateFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_file"
}

func (r *TemplateFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_template_file` manages a template file within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the template file. This must be the name only, so without /var/lib/cobbler/templates.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Description: "The body of the template file. May also point to a file: body = file(\"my_template.ks\").",
				Required:    true,
			},
		},
	}
}

func (r *TemplateFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TemplateFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data templateFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateFile := cobbler.TemplateFile{
		Name: data.Name.ValueString(),
		Body: data.Body.ValueString(),
	}

	tflog.Debug(ctx, "Cobbler TemplateFile: Create", map[string]interface{}{"name": templateFile.Name})

	if err := r.client.CreateTemplateFile(templateFile); err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler TemplateFile", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data templateFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateFile, err := r.client.GetTemplateFile(data.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading TemplateFile", err.Error())
		return
	}

	data.Name = types.StringValue(templateFile.Name)
	data.Body = types.StringValue(templateFile.Body)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data templateFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateFile := cobbler.TemplateFile{
		Name: data.Name.ValueString(),
		Body: data.Body.ValueString(),
	}

	tflog.Debug(ctx, "Cobbler TemplateFile: Update", map[string]interface{}{"name": templateFile.Name})

	// The Cobbler API uses CreateTemplateFile for both create and update
	if err := r.client.CreateTemplateFile(templateFile); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler TemplateFile", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data templateFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler TemplateFile: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteTemplateFile(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler TemplateFile", err.Error())
		return
	}
}

func (r *TemplateFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
