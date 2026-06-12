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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
				Description: "Parent distribution.",
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
				Description: "The parent this profile inherits settings from.",
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
			"boot_files": schema.SingleNestedAttribute{
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
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
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
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
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
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
			"fetchable_files": schema.SingleNestedAttribute{
				Description: "Templates for tftp or wget.",
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
			"mgmt_classes": schema.SingleNestedAttribute{
				Description: "For external configuration management.",
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
			"mgmt_parameters": schema.SingleNestedAttribute{
				Description: "Parameters which will be handed to your management application (Must be a valid YAML dictionary).",
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
			"name_servers_search": schema.SingleNestedAttribute{
				Description: "Name server search settings.",
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
			"template_files": schema.SingleNestedAttribute{
				Description: "File mappings for built-in config management.",
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
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
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
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
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
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
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
		resp.Diagnostics.AddError("Error creating Cobbler Profile", err.Error())
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
		resp.Diagnostics.AddError("Error updating Cobbler Profile", err.Error())
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
	profile.NextServerv4 = data.NextServerV4.ValueString()
	profile.NextServerv6 = data.NextServerV6.ValueString()
	profile.Parent = data.Parent.ValueString()
	profile.Proxy = data.Proxy.ValueString()
	profile.Server = data.Server.ValueString()
	profile.VirtBridge = data.VirtBridge.ValueString()
	profile.VirtCPUs = int(data.VirtCPUs.ValueInt64())
	profile.VirtDiskDriver = stringOrInherit(data.VirtDiskDriver)
	profile.VirtPath = data.VirtPath.ValueString()
	profile.VirtType = stringOrInherit(data.VirtType)

	// ElementsAs fails on null/unknown; guard for Optional+Computed fields not set in config.
	var repos []string
	if !data.Repos.IsNull() && !data.Repos.IsUnknown() {
		diags.Append(data.Repos.ElementsAs(ctx, &repos, false)...)
	}
	profile.Repos = repos

	profile.AutoinstallMeta = inherit.StringMapTo(ctx, data.AutoinstallMeta, diags)
	profile.BootFiles = inherit.StringMapTo(ctx, data.BootFiles, diags)
	profile.EnableIPXE = inherit.BoolTo(ctx, data.EnableIPXE, diags)
	profile.EnableMenu = inherit.BoolTo(ctx, data.EnableMenu, diags)
	profile.FetchableFiles = inherit.StringMapTo(ctx, data.FetchableFiles, diags)
	profile.KernelOptions = inherit.StringMapTo(ctx, data.KernelOptions, diags)
	profile.KernelOptionsPost = inherit.StringMapTo(ctx, data.KernelOptionsPost, diags)
	profile.MgmtClasses = inherit.StringListTo(ctx, data.MgmtClasses, diags)
	profile.MgmtParameters = inherit.StringMapTo(ctx, data.MgmtParameters, diags)
	profile.NameServersSearch = inherit.StringListTo(ctx, data.NameServersSearch, diags)
	profile.NameServers = inherit.StringListTo(ctx, data.NameServers, diags)
	profile.Owners = inherit.StringListTo(ctx, data.Owners, diags)
	profile.TemplateFiles = inherit.StringMapTo(ctx, data.TemplateFiles, diags)
	profile.VirtAutoBoot = inherit.BoolTo(ctx, data.VirtAutoBoot, diags)
	profile.VirtFileSize = inherit.Float64To(ctx, data.VirtFileSize, diags)
	profile.VirtRAM = inherit.IntTo(ctx, data.VirtRAM, diags)

	return profile
}

// profileToModel populates a profileResourceModel from a cobbler.Profile.
func profileToModel(ctx context.Context, profile cobbler.Profile, data *profileResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(profile.Name)
	data.Autoinstall = types.StringValue(profile.Autoinstall)
	data.Comment = types.StringValue(profile.Comment)
	data.DHCPTag = types.StringValue(profile.DHCPTag)
	data.Distro = types.StringValue(profile.Distro)
	data.NextServerV4 = types.StringValue(profile.NextServerv4)
	data.NextServerV6 = types.StringValue(profile.NextServerv6)
	data.Parent = types.StringValue(profile.Parent)
	data.Proxy = types.StringValue(profile.Proxy)
	data.Server = types.StringValue(profile.Server)
	data.VirtBridge = types.StringValue(profile.VirtBridge)
	data.VirtCPUs = types.Int64Value(int64(profile.VirtCPUs))
	data.VirtDiskDriver = types.StringValue(profile.VirtDiskDriver)
	data.VirtPath = types.StringValue(profile.VirtPath)
	data.VirtType = types.StringValue(profile.VirtType)

	repoList, d := types.ListValueFrom(ctx, types.StringType, profile.Repos)
	diags.Append(d...)
	data.Repos = repoList

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, profile.AutoinstallMeta, diags)
	data.BootFiles = inherit.StringMapFrom(ctx, profile.BootFiles, diags)
	data.EnableIPXE = inherit.BoolFrom(ctx, profile.EnableIPXE, diags)
	data.EnableMenu = inherit.BoolFrom(ctx, profile.EnableMenu, diags)
	data.FetchableFiles = inherit.StringMapFrom(ctx, profile.FetchableFiles, diags)
	data.KernelOptions = inherit.StringMapFrom(ctx, profile.KernelOptions, diags)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, profile.KernelOptionsPost, diags)
	data.MgmtClasses = inherit.StringListFrom(ctx, profile.MgmtClasses, diags)
	data.MgmtParameters = inherit.StringMapFrom(ctx, profile.MgmtParameters, diags)
	data.NameServersSearch = inherit.StringListFrom(ctx, profile.NameServersSearch, diags)
	data.NameServers = inherit.StringListFrom(ctx, profile.NameServers, diags)
	data.Owners = inherit.StringListFrom(ctx, profile.Owners, diags)
	data.TemplateFiles = inherit.StringMapFrom(ctx, profile.TemplateFiles, diags)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, profile.VirtAutoBoot, diags)
	data.VirtFileSize = inherit.Float64From(ctx, profile.VirtFileSize, diags)
	data.VirtRAM = inherit.IntFrom(ctx, profile.VirtRAM, diags)
}
