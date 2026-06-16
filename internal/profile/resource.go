package profile

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ProfileResource{}
var _ resource.ResourceWithImportState = &ProfileResource{}

type ProfileResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &ProfileResource{}
}

func (r *ProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *ProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_profile` manages a profile within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the profile.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uid": schema.StringAttribute{
				Description: "Server-assigned UID for this profile. Use this as the value for `cobbler_profile.parent` or `cobbler_system.profile`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"autoinstall": schema.StringAttribute{
				Description: "Template remote kickstarts or preseeds.",
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
			"dhcp_tag": schema.StringAttribute{
				Description: "DHCP tag.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"distro": schema.StringAttribute{
				Description: "The Cobbler UID of the parent distribution. Use `cobbler_distro.foo.uid`.",
				Required:    true,
			},
			"next_server_v4": schema.StringAttribute{
				Description: "The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"next_server_v6": schema.StringAttribute{
				Description: "The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent": schema.StringAttribute{
				Description: "The Cobbler UID of the parent profile this profile inherits settings from. Use `cobbler_profile.foo.uid`.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.StringAttribute{
				Description: "The server-override for the profile.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_bridge": schema.StringAttribute{
				Description: "The bridge for virtual machines.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_cpus": schema.Int64Attribute{
				Description: "The number of virtual CPUs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"virt_disk_driver": schema.StringAttribute{
				Description: "The virtual machine disk driver.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_path": schema.StringAttribute{
				Description: "The virtual machine path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_type": schema.StringAttribute{
				Description: "The type of virtual machine. Valid options are: xenpv, xenfv, qemu, kvm, vmware, openvz.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"repos": schema.ListAttribute{
				Description: "Repos to auto-assign to this profile.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"autoinstall_meta": schema.SingleNestedAttribute{
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"enable_ipxe": schema.SingleNestedAttribute{
				Description: "Use iPXE instead of PXELINUX for advanced booting options.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"enable_menu": schema.SingleNestedAttribute{
				Description: "Enable a boot menu.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options for the profile.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"kernel_options_post": schema.SingleNestedAttribute{
				Description: "Post install kernel options.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"name_servers_search": schema.ListAttribute{
				Description: "Name server search settings. Not inheritable.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"name_servers": schema.SingleNestedAttribute{
				Description: "Name servers.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management. Not inheritable.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_auto_boot": schema.SingleNestedAttribute{
				Description: "Auto boot virtual machines.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.BoolAttribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"virt_file_size": schema.SingleNestedAttribute{
				Description: "The virtual machine file size.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"virt_ram": schema.SingleNestedAttribute{
				Description: "The amount of RAM for the virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{
						Description: "The value.",
						Optional:    true,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "If true, inherited from parent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *ProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data profileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := modelToProfile(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Profile: Create", map[string]interface{}{"name": profile.Name})

	newProfile, err := r.client.CreateProfile(profile)
	if err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error creating Cobbler Profile", err)
		return
	}

	profileToModel(ctx, *newProfile, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data profileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.GetProfile(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Profile", err.Error())
		return
	}

	profileToModel(ctx, *profile, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data profileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := modelToProfile(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Profile: Update", map[string]interface{}{"name": profile.Name})

	if err := r.client.UpdateProfile(&profile); err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error updating Cobbler Profile", err)
		return
	}

	updatedProfile, err := r.client.GetProfile(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Profile after update", err.Error())
		return
	}

	profileToModel(ctx, *updatedProfile, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data profileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Profile: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteProfile(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Profile", err.Error())
	}
}

func (r *ProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// stringOrInherit returns "<<inherit>>" when s is null, unknown, or empty.
// Cobbler rejects empty strings for enum-validated fields (e.g. virt_disk_driver, virt_type).
func stringOrInherit(s types.String) string {
	if v := s.ValueString(); v != "" {
		return v
	}
	return "<<inherit>>"
}

// modelToProfile converts a profileResourceModel to a cobbler.Profile.
func modelToProfile(ctx context.Context, data profileResourceModel, diags *diag.Diagnostics) cobbler.Profile {
	profile := cobbler.NewProfile()
	profile.Name = data.Name.ValueString()
	profile.Autoinstall = data.Autoinstall.ValueString()
	profile.Comment = data.Comment.ValueString()
	profile.DHCPTag = data.DHCPTag.ValueString()
	profile.Distro = data.Distro.ValueString()
	profile.TFTP.NextServerV4 = data.NextServerV4.ValueString()
	profile.TFTP.NextServerV6 = data.NextServerV6.ValueString()
	profile.Parent = data.Parent.ValueString()
	profile.Proxy = data.Proxy.ValueString()
	profile.Server = data.Server.ValueString()
	profile.VirtBridge = data.VirtBridge.ValueString()
	profile.Virt.Cpus = cobbler.Value[int]{Data: int(data.VirtCPUs.ValueInt64())}
	profile.Virt.DiskDriver = stringOrInherit(data.VirtDiskDriver)
	profile.Virt.Path = data.VirtPath.ValueString()
	profile.Virt.Type = stringOrInherit(data.VirtType)

	// ElementsAs fails on null/unknown; guard for Optional+Computed fields not set in config.
	var repos []string
	if !data.Repos.IsNull() && !data.Repos.IsUnknown() {
		diags.Append(data.Repos.ElementsAs(ctx, &repos, false)...)
	}
	profile.Repos = repos

	var nameServersSearch []string
	if !data.NameServersSearch.IsNull() && !data.NameServersSearch.IsUnknown() {
		diags.Append(data.NameServersSearch.ElementsAs(ctx, &nameServersSearch, false)...)
	}
	profile.DNS.NameServersSearch = nameServersSearch

	var templateFiles map[string]string
	if !data.TemplateFiles.IsNull() && !data.TemplateFiles.IsUnknown() {
		diags.Append(data.TemplateFiles.ElementsAs(ctx, &templateFiles, false)...)
	}
	profile.TemplateFiles = templateFiles

	profile.AutoinstallMeta = inherit.StringMapTo(ctx, data.AutoinstallMeta, diags)
	profile.EnableIPXE = inherit.BoolTo(ctx, data.EnableIPXE, diags)
	profile.EnableMenu = inherit.BoolTo(ctx, data.EnableMenu, diags)
	profile.KernelOptions = inherit.StringMapTo(ctx, data.KernelOptions, diags)
	profile.KernelOptionsPost = inherit.StringMapTo(ctx, data.KernelOptionsPost, diags)
	profile.DNS.NameServers = inherit.StringListTo(ctx, data.NameServers, diags)
	profile.Owners = inherit.StringListTo(ctx, data.Owners, diags)
	profile.Virt.AutoBoot = inherit.BoolTo(ctx, data.VirtAutoBoot, diags)
	profile.Virt.FileSize = inherit.Float64To(ctx, data.VirtFileSize, diags)
	profile.Virt.Ram = inherit.IntTo(ctx, data.VirtRAM, diags)

	return profile
}

// profileToModel populates a profileResourceModel from a cobbler.Profile.
func profileToModel(ctx context.Context, profile cobbler.Profile, data *profileResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(profile.Name)
	data.UID = types.StringValue(profile.Uid)
	data.Autoinstall = types.StringValue(profile.Autoinstall)
	data.Comment = types.StringValue(profile.Comment)
	data.DHCPTag = types.StringValue(profile.DHCPTag)
	data.Distro = types.StringValue(profile.Distro)
	data.NextServerV4 = types.StringValue(profile.TFTP.NextServerV4)
	data.NextServerV6 = types.StringValue(profile.TFTP.NextServerV6)
	data.Parent = types.StringValue(profile.Parent)
	data.Proxy = types.StringValue(profile.Proxy)
	data.Server = types.StringValue(profile.Server)
	data.VirtBridge = types.StringValue(profile.VirtBridge)
	data.VirtCPUs = types.Int64Value(int64(profile.Virt.Cpus.Data))
	data.VirtDiskDriver = types.StringValue(profile.Virt.DiskDriver)
	data.VirtPath = types.StringValue(profile.Virt.Path)
	data.VirtType = types.StringValue(profile.Virt.Type)

	repoList, d := types.ListValueFrom(ctx, types.StringType, profile.Repos)
	diags.Append(d...)
	data.Repos = repoList

	nameServersSearch, d := types.ListValueFrom(ctx, types.StringType, profile.DNS.NameServersSearch)
	diags.Append(d...)
	data.NameServersSearch = nameServersSearch

	templateFiles, d := types.MapValueFrom(ctx, types.StringType, profile.TemplateFiles)
	diags.Append(d...)
	data.TemplateFiles = templateFiles

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, profile.AutoinstallMeta, diags)
	data.EnableIPXE = inherit.BoolFrom(ctx, profile.EnableIPXE, diags)
	data.EnableMenu = inherit.BoolFrom(ctx, profile.EnableMenu, diags)
	data.KernelOptions = inherit.StringMapFrom(ctx, profile.KernelOptions, diags)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, profile.KernelOptionsPost, diags)
	data.NameServers = inherit.StringListFrom(ctx, profile.DNS.NameServers, diags)
	data.Owners = inherit.StringListFrom(ctx, profile.Owners, diags)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, profile.Virt.AutoBoot, diags)
	data.VirtFileSize = inherit.Float64From(ctx, profile.Virt.FileSize, diags)
	data.VirtRAM = inherit.IntFrom(ctx, profile.Virt.Ram, diags)
}
