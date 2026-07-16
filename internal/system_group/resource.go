package system_group

import (
	"context"
	"strings"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SystemGroupResource{}
var _ resource.ResourceWithImportState = &SystemGroupResource{}

type SystemGroupResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &SystemGroupResource{}
}

func (r *SystemGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_group"
}

func (r *SystemGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_system_group` manages a Cobbler 4.0.0+ system group (a named collection of distros for bulk operations).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the group. Changing this forces a new resource.",
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
			"items": schema.ListAttribute{
				Description: "Names of the distros belonging to this group.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SystemGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func modelToGroup(ctx context.Context, data systemGroupResourceModel, diags *diag.Diagnostics) cobbler.SystemGroup {
	g := cobbler.NewSystemGroup()
	g.Name = data.Name.ValueString()
	g.Comment = data.Comment.ValueString()

	var items []string
	if !data.Items.IsNull() && !data.Items.IsUnknown() {
		diags.Append(data.Items.ElementsAs(ctx, &items, false)...)
	}
	if items == nil {
		items = []string{}
	}
	g.Members = items
	return g
}

func groupToModel(ctx context.Context, g cobbler.SystemGroup, data *systemGroupResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(g.Name)
	data.Comment = types.StringValue(g.Comment)

	items := g.Members
	if items == nil {
		items = []string{}
	}
	l, d := types.ListValueFrom(ctx, types.StringType, items)
	diags.Append(d...)
	data.Items = l
}

func groupToDataSourceModel(ctx context.Context, g cobbler.SystemGroup, data *systemGroupDataSourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(g.Name)
	data.Comment = types.StringValue(g.Comment)

	items := g.Members
	if items == nil {
		items = []string{}
	}
	l, d := types.ListValueFrom(ctx, types.StringType, items)
	diags.Append(d...)
	data.Items = l
}

func (r *SystemGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data systemGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	g := modelToGroup(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler SystemGroup: Create", map[string]interface{}{"name": g.Name})

	created, err := r.client.CreateSystemGroup(g)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler SystemGroup", err.Error())
		return
	}

	groupToModel(ctx, *created, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data systemGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	g, err := r.client.GetSystemGroup(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Cobbler SystemGroup", err.Error())
		return
	}

	groupToModel(ctx, *g, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data systemGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	g := modelToGroup(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler SystemGroup: Update", map[string]interface{}{"name": g.Name})

	if err := r.client.UpdateSystemGroup(&g); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler SystemGroup", err.Error())
		return
	}

	updated, err := r.client.GetSystemGroup(g.Name, false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler SystemGroup after update", err.Error())
		return
	}

	groupToModel(ctx, *updated, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data systemGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler SystemGroup: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteSystemGroup(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler SystemGroup", err.Error())
	}
}

func (r *SystemGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
