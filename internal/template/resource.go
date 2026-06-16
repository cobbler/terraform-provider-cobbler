package template

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TemplateResource{}
var _ resource.ResourceWithImportState = &TemplateResource{}

type TemplateResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &TemplateResource{}
}

func (r *TemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

var uriAttrTypes = map[string]attr.Type{
	"schema": types.StringType,
	"path":   types.StringType,
}

func (r *TemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_template` manages an autoinstall template within Cobbler (4.0.0+). Replaces the legacy `cobbler_snippet` and `cobbler_template_file` resources.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the template. Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template_type": schema.StringAttribute{
				Description: "The template engine to use, e.g. `jinja2`, `cheetah`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("jinja2"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uri": schema.SingleNestedAttribute{
				Description: "Where the template's content lives.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"schema": schema.StringAttribute{
						Description: "Source type. One of: `file`, `importlib`, `environment`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("file"),
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("file", "importlib", "environment"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"path": schema.StringAttribute{
						Description: "Source-relative path to the template content.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"tags": schema.ListAttribute{
				Description: "Tags associated with the template.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The template body.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"built_in": schema.BoolAttribute{
				Description: "Whether the template is built into Cobbler (read-only).",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func parseTemplateSchema(s string) cobbler.TemplateSchema {
	switch s {
	case "importlib":
		return cobbler.TemplateSchemaImportlib
	case "environment":
		return cobbler.TemplateSchemaEnvironment
	default:
		return cobbler.TemplateSchemaFile
	}
}

func uriFromAPI(_ context.Context, u cobbler.URIOption, diags *diag.Diagnostics) types.Object {
	obj, d := types.ObjectValue(uriAttrTypes, map[string]attr.Value{
		"schema": types.StringValue(u.Schema.String()),
		"path":   types.StringValue(u.Path),
	})
	diags.Append(d...)
	return obj
}

func uriToAPI(_ context.Context, obj types.Object, _ *diag.Diagnostics) cobbler.URIOption {
	if obj.IsNull() || obj.IsUnknown() {
		return cobbler.URIOption{Schema: cobbler.TemplateSchemaFile}
	}
	attrs := obj.Attributes()
	sch, _ := attrs["schema"].(types.String)
	pth, _ := attrs["path"].(types.String)
	return cobbler.URIOption{
		Schema: parseTemplateSchema(sch.ValueString()),
		Path:   pth.ValueString(),
	}
}

func modelToTemplate(ctx context.Context, data templateResourceModel, diags *diag.Diagnostics) cobbler.Template {
	tpl := cobbler.NewTemplate()
	tpl.Name = data.Name.ValueString()
	tpl.Comment = data.Comment.ValueString()
	if v := data.TemplateType.ValueString(); v != "" {
		tpl.TemplateType = v
	}
	tpl.URI = uriToAPI(ctx, data.URI, diags)
	tpl.Content = data.Content.ValueString()

	var tags []string
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		diags.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if tags == nil {
		tags = []string{}
	}
	tpl.Tags = tags

	return tpl
}

func templateToModel(ctx context.Context, tpl cobbler.Template, data *templateResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(tpl.Name)
	data.Comment = types.StringValue(tpl.Comment)
	data.TemplateType = types.StringValue(tpl.TemplateType)
	data.URI = uriFromAPI(ctx, tpl.URI, diags)
	data.Content = types.StringValue(tpl.Content)
	data.BuiltIn = types.BoolValue(tpl.BuiltIn)

	tags := tpl.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsList, d := types.ListValueFrom(ctx, types.StringType, tags)
	diags.Append(d...)
	data.Tags = tagsList
}

func templateToDataSourceModel(ctx context.Context, tpl cobbler.Template, data *templateDataSourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(tpl.Name)
	data.Comment = types.StringValue(tpl.Comment)
	data.TemplateType = types.StringValue(tpl.TemplateType)
	data.URI = uriFromAPI(ctx, tpl.URI, diags)
	data.Content = types.StringValue(tpl.Content)
	data.BuiltIn = types.BoolValue(tpl.BuiltIn)

	tags := tpl.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsList, d := types.ListValueFrom(ctx, types.StringType, tags)
	diags.Append(d...)
	data.Tags = tagsList
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data templateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpl := modelToTemplate(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Template: Create", map[string]interface{}{"name": tpl.Name})

	created, err := r.client.CreateTemplate(tpl)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler Template", err.Error())
		return
	}

	templateToModel(ctx, *created, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data templateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpl, err := r.client.GetTemplate(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Cobbler Template", err.Error())
		return
	}

	templateToModel(ctx, *tpl, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data templateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpl := modelToTemplate(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Template: Update", map[string]interface{}{"name": tpl.Name})

	if err := r.client.UpdateTemplate(&tpl); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler Template", err.Error())
		return
	}

	updated, err := r.client.GetTemplate(tpl.Name, false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Template after update", err.Error())
		return
	}

	templateToModel(ctx, *updated, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data templateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Template: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteTemplate(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Template", err.Error())
	}
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
