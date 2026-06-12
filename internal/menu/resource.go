package menu

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &MenuResource{}
var _ resource.ResourceWithImportState = &MenuResource{}

type MenuResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &MenuResource{}
}

func (r *MenuResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_menu"
}

func (r *MenuResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_menu` manages a boot menu within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "A name for the menu. Changing this forces a new resource.",
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
			"parent": schema.StringAttribute{
				Description: "The name of the parent menu. Used for hierarchical menus.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name shown in the boot menu.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"autoinstall_meta": schema.SingleNestedAttribute{
				Description: "Autoinstall template metadata.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedMapAttrs(),
			},
			"fetchable_files": schema.SingleNestedAttribute{
				Description: "Templates for tftp or wget.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedMapAttrs(),
			},
			"boot_files": schema.SingleNestedAttribute{
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedMapAttrs(),
			},
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"mgmt_classes": schema.SingleNestedAttribute{
				Description: "Management classes for external config management.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedListAttrs(),
			},
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedListAttrs(),
			},
			"mgmt_parameters": schema.SingleNestedAttribute{
				Description: "Parameters for external management systems.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: inheritedMapAttrs(),
			},
		},
	}
}

// inheritedMapAttrs returns the sub-attributes for a SingleNestedAttribute
// that wraps a Cobbler Value[map[string]interface{}] (inherited or explicit map).
func inheritedMapAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Map{
				mapplanmodifier.UseStateForUnknown(),
			},
		},
		"inherited": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

// inheritedListAttrs returns the sub-attributes for a SingleNestedAttribute
// that wraps a Cobbler Value[[]string] (inherited or explicit list).
func inheritedListAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"inherited": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

func (r *MenuResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MenuResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data menuResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	menu := modelToMenu(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Menu: Create", map[string]interface{}{"name": menu.Name})

	newMenu, err := r.client.CreateMenu(menu)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler Menu", err.Error())
		return
	}

	menuToModel(ctx, *newMenu, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MenuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data menuResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	menu, err := r.client.GetMenu(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Cobbler Menu", err.Error())
		return
	}

	menuToModel(ctx, *menu, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MenuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data menuResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	menu := modelToMenu(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Menu: Update", map[string]interface{}{"name": menu.Name})

	if err := r.client.UpdateMenu(&menu); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler Menu", err.Error())
		return
	}

	updatedMenu, err := r.client.GetMenu(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Menu after update", err.Error())
		return
	}

	menuToModel(ctx, *updatedMenu, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MenuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data menuResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Menu: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteMenu(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Menu", err.Error())
	}
}

func (r *MenuResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// modelToMenu converts a menuResourceModel to a cobbler.Menu.
func modelToMenu(ctx context.Context, data menuResourceModel, diags *diag.Diagnostics) cobbler.Menu {
	menu := cobbler.NewMenu()
	menu.Name = data.Name.ValueString()
	menu.Comment = data.Comment.ValueString()
	menu.Parent = data.Parent.ValueString()
	menu.DisplayName = data.DisplayName.ValueString()
	menu.AutoinstallMeta = inherit.StringMapTo(ctx, data.AutoinstallMeta, diags)
	menu.FetchableFiles = inherit.StringMapTo(ctx, data.FetchableFiles, diags)
	menu.BootFiles = inherit.StringMapTo(ctx, data.BootFiles, diags)
	menu.MgmtClasses = inherit.StringListTo(ctx, data.MgmtClasses, diags)
	menu.Owners = inherit.StringListTo(ctx, data.Owners, diags)
	menu.MgmtParameters = inherit.StringMapTo(ctx, data.MgmtParameters, diags)

	var templateFiles map[string]interface{}
	if !data.TemplateFiles.IsNull() && !data.TemplateFiles.IsUnknown() {
		diags.Append(data.TemplateFiles.ElementsAs(ctx, &templateFiles, false)...)
	}
	menu.TemplateFiles = cobbler.Value[map[string]interface{}]{Data: templateFiles, IsInherited: false}

	return menu
}

// menuToModel populates a menuResourceModel from a cobbler.Menu.
func menuToModel(ctx context.Context, menu cobbler.Menu, data *menuResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(menu.Name)
	data.Comment = types.StringValue(menu.Comment)
	data.Parent = types.StringValue(menu.Parent)
	data.DisplayName = types.StringValue(menu.DisplayName)
	data.AutoinstallMeta = inherit.StringMapFrom(ctx, menu.AutoinstallMeta, diags)
	data.FetchableFiles = inherit.StringMapFrom(ctx, menu.FetchableFiles, diags)
	data.BootFiles = inherit.StringMapFrom(ctx, menu.BootFiles, diags)
	data.MgmtClasses = inherit.StringListFrom(ctx, menu.MgmtClasses, diags)
	data.Owners = inherit.StringListFrom(ctx, menu.Owners, diags)
	data.MgmtParameters = inherit.StringMapFrom(ctx, menu.MgmtParameters, diags)

	templateFiles, d := types.MapValueFrom(ctx, types.StringType, menu.TemplateFiles.Data)
	diags.Append(d...)
	data.TemplateFiles = templateFiles
}
