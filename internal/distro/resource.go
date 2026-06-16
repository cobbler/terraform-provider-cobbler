package distro

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DistroResource{}
var _ resource.ResourceWithImportState = &DistroResource{}

type DistroResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &DistroResource{}
}

func (r *DistroResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_distro"
}

func (r *DistroResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_distro` manages a distribution within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "A name for the distro.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the distro. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"breed": schema.StringAttribute{
				Description: "The \"breed\" of distribution. Valid options are: redhat, fedora, centos, scientific linux, suse, debian, and ubuntu.",
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
			"initrd": schema.StringAttribute{
				Description: "Absolute path to initrd on filesystem. This must already exist prior to creating the distro.",
				Required:    true,
			},
			"kernel": schema.StringAttribute{
				Description: "Absolute path to kernel on filesystem. This must already exist prior to creating the distro.",
				Required:    true,
			},
			"remote_boot_initrd": schema.StringAttribute{
				Description: "URL the bootloader directly retrieves and boots from.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"remote_boot_kernel": schema.StringAttribute{
				Description: "URL the bootloader directly retrieves and boots from.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_version": schema.StringAttribute{
				Description: "The version of the distro you are creating. Example: `focal`.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"boot_loaders": schema.SingleNestedAttribute{
				Description: "Must be either 'grub', 'pxe', or 'ipxe'.",
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
			"fetchable_files": schema.SingleNestedAttribute{
				Description: "Templates for tftp or wget.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
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
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options to use with the kernel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
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
			"kernel_options_post": schema.SingleNestedAttribute{
				Description: "Post install kernel options to use with the kernel after installation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
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
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
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
			"template_files": schema.MapAttribute{
				Description: "File mappings for built-in config management.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *DistroResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DistroResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data distroResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	distro := modelToDistro(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Distro: Create", map[string]interface{}{"name": distro.Name})

	newDistro, err := r.client.CreateDistro(distro)
	if err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error creating Cobbler Distro", err)
		return
	}

	distroToModel(ctx, *newDistro, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DistroResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data distroResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	distro, err := r.client.GetDistro(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Distro", err.Error())
		return
	}

	distroToModel(ctx, *distro, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DistroResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data distroResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	distro := modelToDistro(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Distro: Update", map[string]interface{}{"name": distro.Name})

	if err := r.client.UpdateDistro(&distro); err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error updating Cobbler Distro", err)
		return
	}

	updatedDistro, err := r.client.GetDistro(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Distro after update", err.Error())
		return
	}

	distroToModel(ctx, *updatedDistro, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DistroResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data distroResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Distro: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteDistro(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Distro", err.Error())
	}
}

func (r *DistroResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// modelToDistro converts a distroResourceModel to a cobbler.Distro.
func modelToDistro(ctx context.Context, data distroResourceModel, diags *diag.Diagnostics) cobbler.Distro {
	distro := cobbler.NewDistro()
	distro.Name = data.Name.ValueString()
	distro.Arch = data.Arch.ValueString()
	distro.Breed = data.Breed.ValueString()
	distro.Comment = data.Comment.ValueString()
	distro.Initrd = data.Initrd.ValueString()
	distro.Kernel = data.Kernel.ValueString()
	distro.RemoteBootInitrd = data.RemoteBootInitrd.ValueString()
	distro.RemoteBootKernel = data.RemoteBootKernel.ValueString()
	distro.OSVersion = data.OSVersion.ValueString()
	distro.BootFiles = inherit.StringMapTo(ctx, data.BootFiles, diags)
	distro.BootLoaders = inherit.StringListTo(ctx, data.BootLoaders, diags)
	distro.FetchableFiles = inherit.StringMapTo(ctx, data.FetchableFiles, diags)
	distro.KernelOptions = inherit.StringMapTo(ctx, data.KernelOptions, diags)
	distro.KernelOptionsPost = inherit.StringMapTo(ctx, data.KernelOptionsPost, diags)
	distro.Owners = inherit.StringListTo(ctx, data.Owners, diags)

	// ElementsAs fails on null/unknown values; guard to avoid plan-time errors
	// for Optional+Computed fields not set in the configuration.
	var templateFiles map[string]interface{}
	if !data.TemplateFiles.IsNull() && !data.TemplateFiles.IsUnknown() {
		diags.Append(data.TemplateFiles.ElementsAs(ctx, &templateFiles, false)...)
	}
	distro.TemplateFiles = cobbler.Value[map[string]interface{}]{Data: templateFiles, IsInherited: false}
	return distro
}

// distroToModel populates a distroResourceModel from a cobbler.Distro.
func distroToModel(ctx context.Context, distro cobbler.Distro, data *distroResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(distro.Name)
	data.Arch = types.StringValue(distro.Arch)
	data.Breed = types.StringValue(distro.Breed)
	data.Comment = types.StringValue(distro.Comment)
	data.Initrd = types.StringValue(distro.Initrd)
	data.Kernel = types.StringValue(distro.Kernel)
	data.RemoteBootInitrd = types.StringValue(distro.RemoteBootInitrd)
	data.RemoteBootKernel = types.StringValue(distro.RemoteBootKernel)
	data.OSVersion = types.StringValue(distro.OSVersion)
	data.BootFiles = inherit.StringMapFrom(ctx, distro.BootFiles, diags)
	data.BootLoaders = inherit.StringListFrom(ctx, distro.BootLoaders, diags)
	data.FetchableFiles = inherit.StringMapFrom(ctx, distro.FetchableFiles, diags)
	data.KernelOptions = inherit.StringMapFrom(ctx, distro.KernelOptions, diags)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, distro.KernelOptionsPost, diags)
	data.Owners = inherit.StringListFrom(ctx, distro.Owners, diags)
	templateFiles, d := types.MapValueFrom(ctx, types.StringType, distro.TemplateFiles.Data)
	diags.Append(d...)
	data.TemplateFiles = templateFiles
}
