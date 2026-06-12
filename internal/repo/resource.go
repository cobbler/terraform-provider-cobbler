package repo

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RepoResource{}
var _ resource.ResourceWithImportState = &RepoResource{}

type RepoResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &RepoResource{}
}

func (r *RepoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repo"
}

func (r *RepoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_repo` manages a repo within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "A name for the repo.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the repo. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"breed": schema.StringAttribute{
				Description: "The \"breed\" of distribution. Valid options are: rsync, rhn, yum, apt, and wget.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"environment": schema.MapAttribute{
				Description: "Environment variables to use during repo command execution.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_updated": schema.BoolAttribute{
				Description: "Update the repo upon Cobbler sync. Valid values are true or false.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mirror": schema.StringAttribute{
				Description: "Address of the repo to mirror.",
				Required:    true,
			},
			"mirror_locally": schema.BoolAttribute{
				Description: "Whether to copy the files locally or just references to the external files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"apt_components": schema.ListAttribute{
				Description: "List of Apt components such as main, restricted, universe. Applicable to apt breeds only.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"apt_dists": schema.ListAttribute{
				Description: "List of Apt distribution names such as focal, focal-updates. Applicable to apt breeds only.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"rpm_list": schema.ListAttribute{
				Description: "List of specific RPMs to mirror.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"createrepo_flags": schema.SingleNestedAttribute{
				Description: "Flags to use with `createrepo`.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"inherited": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "List of Owners for authz_ownership.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"proxy": schema.SingleNestedAttribute{
				Description: "Proxy to use for downloading the repo.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"inherited": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *RepoResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RepoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data repoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo := modelToRepo(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Repo: Create", map[string]interface{}{"name": repo.Name})

	newRepo, err := r.client.CreateRepo(repo)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler Repo", err.Error())
		return
	}

	repoToModel(ctx, *newRepo, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data repoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo, err := r.client.GetRepo(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Repo", err.Error())
		return
	}

	repoToModel(ctx, *repo, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data repoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo := modelToRepo(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Repo: Update", map[string]interface{}{"name": repo.Name})

	if err := r.client.UpdateRepo(&repo); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler Repo", err.Error())
		return
	}

	updatedRepo, err := r.client.GetRepo(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Repo after update", err.Error())
		return
	}

	repoToModel(ctx, *updatedRepo, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data repoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Repo: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteRepo(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Repo", err.Error())
	}
}

func (r *RepoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// modelToRepo converts a repoResourceModel to a cobbler.Repo.
func modelToRepo(ctx context.Context, data repoResourceModel, diags *diag.Diagnostics) cobbler.Repo {
	repo := cobbler.NewRepo()
	repo.Name = data.Name.ValueString()
	repo.Arch = data.Arch.ValueString()
	repo.Breed = data.Breed.ValueString()
	repo.Comment = data.Comment.ValueString()
	repo.KeepUpdated = data.KeepUpdated.ValueBool()
	repo.Mirror = data.Mirror.ValueString()
	repo.MirrorLocally = data.MirrorLocally.ValueBool()
	repo.CreateRepoFlags = inherit.StringTo(ctx, data.CreateRepoFlags, diags)
	repo.Owners = inherit.StringListTo(ctx, data.Owners, diags)
	repo.Proxy = inherit.StringTo(ctx, data.Proxy, diags)

	// ElementsAs fails on null/unknown values; guard to avoid plan-time errors
	// for Optional+Computed fields not set in the configuration.
	var env map[string]string
	if !data.Environment.IsNull() && !data.Environment.IsUnknown() {
		diags.Append(data.Environment.ElementsAs(ctx, &env, false)...)
	}
	repo.Environment = env

	var aptComponents []string
	if !data.AptComponents.IsNull() && !data.AptComponents.IsUnknown() {
		diags.Append(data.AptComponents.ElementsAs(ctx, &aptComponents, false)...)
	}
	repo.AptComponents = aptComponents

	var aptDists []string
	if !data.AptDists.IsNull() && !data.AptDists.IsUnknown() {
		diags.Append(data.AptDists.ElementsAs(ctx, &aptDists, false)...)
	}
	repo.AptDists = aptDists

	var rpmList []string
	if !data.RpmList.IsNull() && !data.RpmList.IsUnknown() {
		diags.Append(data.RpmList.ElementsAs(ctx, &rpmList, false)...)
	}
	repo.RpmList = rpmList

	return repo
}

// repoToModel populates a repoResourceModel from a cobbler.Repo.
func repoToModel(ctx context.Context, repo cobbler.Repo, data *repoResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(repo.Name)
	data.Arch = types.StringValue(repo.Arch)
	data.Breed = types.StringValue(repo.Breed)
	data.Comment = types.StringValue(repo.Comment)
	data.KeepUpdated = types.BoolValue(repo.KeepUpdated)
	data.Mirror = types.StringValue(repo.Mirror)
	data.MirrorLocally = types.BoolValue(repo.MirrorLocally)
	data.CreateRepoFlags = inherit.StringFrom(ctx, repo.CreateRepoFlags, diags)
	data.Owners = inherit.StringListFrom(ctx, repo.Owners, diags)
	data.Proxy = inherit.StringFrom(ctx, repo.Proxy, diags)

	env, d := types.MapValueFrom(ctx, types.StringType, repo.Environment)
	diags.Append(d...)
	data.Environment = env

	aptComponents, d := types.ListValueFrom(ctx, types.StringType, repo.AptComponents)
	diags.Append(d...)
	data.AptComponents = aptComponents

	aptDists, d := types.ListValueFrom(ctx, types.StringType, repo.AptDists)
	diags.Append(d...)
	data.AptDists = aptDists

	rpmList, d := types.ListValueFrom(ctx, types.StringType, repo.RpmList)
	diags.Append(d...)
	data.RpmList = rpmList
}
