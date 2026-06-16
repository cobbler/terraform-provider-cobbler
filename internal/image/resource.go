package image

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ImageResource{}
var _ resource.ResourceWithImportState = &ImageResource{}

type ImageResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &ImageResource{}
}

func (r *ImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (r *ImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_image` manages an image within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "A name for the image.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file": schema.StringAttribute{
				Description: "Path to the image media. Format depends on `image_type`.",
				Required:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the image. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"autoinstall": schema.StringAttribute{
				Description: "Path to an autoinstall file (e.g. kickstart, preseed). Leave empty to inherit from defaults.",
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
			"image_type": schema.StringAttribute{
				Description: "Type of image. Valid options are: direct, iso, memdisk, virt-clone.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_version": schema.StringAttribute{
				Description: "The OS version the image contains. Example: `focal`.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"boot_loaders": schema.ListAttribute{
				Description: "Boot loaders supported by the image. Must be subset of: 'grub', 'pxe', 'ipxe'.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"menu": schema.StringAttribute{
				Description: "Name of the parent Cobbler menu that this image appears under.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_auto_boot": schema.BoolAttribute{
				Description: "Whether to auto-boot the virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_bridge": schema.StringAttribute{
				Description: "Bridge for the virtual machine to attach to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_cpus": schema.Int64Attribute{
				Description: "Number of CPUs to allocate to the virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"virt_disk_driver": schema.StringAttribute{
				Description: "Disk driver for the virtual machine. Valid options are: raw, qcow2, qed, vdi, vdmk. Leave empty to inherit.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_file_size": schema.SingleNestedAttribute{
				Description: "Disk file size in GB for the virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{
						Optional: true,
						Computed: true,
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
			"virt_path": schema.StringAttribute{
				Description: "Path on the virtualization host where the image is stored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virt_ram": schema.SingleNestedAttribute{
				Description: "RAM in MB to allocate to the virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{
						Optional: true,
						Computed: true,
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
			"virt_type": schema.StringAttribute{
				Description: "Virtualization type. Valid options are: qemu, kvm, xenpv, xenfv, vmware, vmwarew, openvz, auto. Leave empty to inherit.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"mgmt_classes": schema.SingleNestedAttribute{
				Description: "Management classes for external config management.",
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

func (r *ImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data imageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	image := modelToImage(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Image: Create", map[string]interface{}{"name": image.Name})

	newImage, err := r.client.CreateImage(image)
	if err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error creating Cobbler Image", err)
		return
	}

	imageToModel(ctx, *newImage, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	image, err := r.client.GetImage(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Image", err.Error())
		return
	}

	imageToModel(ctx, *image, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data imageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	image := modelToImage(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Image: Update", map[string]interface{}{"name": image.Name})

	if err := r.client.UpdateImage(&image); err != nil {
		clientpkg.AddClientError(&resp.Diagnostics, "Error updating Cobbler Image", err)
		return
	}

	updatedImage, err := r.client.GetImage(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Image after update", err.Error())
		return
	}

	imageToModel(ctx, *updatedImage, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler Image: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteImage(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler Image", err.Error())
	}
}

func (r *ImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// stringOrInherit converts a Terraform string to its Cobbler equivalent, mapping
// empty / null / unknown to the Cobbler "<<inherit>>" sentinel. Cobbler enum
// validators (autoinstall, virt_type, virt_disk_driver) reject the empty string,
// and the terraform-plugin-framework marks Optional+Computed strings null in the
// config as Unknown in the plan, so we must translate before sending.
func stringOrInherit(s types.String) string {
	if s.IsNull() || s.IsUnknown() {
		return "<<inherit>>"
	}
	if v := s.ValueString(); v != "" {
		return v
	}
	return "<<inherit>>"
}

// modelToImage converts an imageResourceModel to a cobbler.Image.
func modelToImage(ctx context.Context, data imageResourceModel, diags *diag.Diagnostics) cobbler.Image {
	image := cobbler.NewImage()
	image.Name = data.Name.ValueString()
	image.File = data.File.ValueString()
	image.Arch = data.Arch.ValueString()
	image.Autoinstall = stringOrInherit(data.Autoinstall)
	image.Breed = data.Breed.ValueString()
	image.Comment = data.Comment.ValueString()
	image.ImageType = data.ImageType.ValueString()
	image.OsVersion = data.OSVersion.ValueString()
	image.Menu = data.Menu.ValueString()
	image.VirtAutoBoot = data.VirtAutoBoot.ValueBool()
	image.VirtBridge = data.VirtBridge.ValueString()
	image.VirtCpus = int(data.VirtCpus.ValueInt64())
	image.VirtDiskDriver = stringOrInherit(data.VirtDiskDriver)
	image.VirtFileSize = inherit.Float64To(ctx, data.VirtFileSize, diags)
	image.VirtPath = data.VirtPath.ValueString()
	image.VirtRam = inherit.IntTo(ctx, data.VirtRam, diags)
	image.VirtType = stringOrInherit(data.VirtType)
	image.KernelOptions = inherit.StringMapTo(ctx, data.KernelOptions, diags)
	image.KernelOptionsPost = inherit.StringMapTo(ctx, data.KernelOptionsPost, diags)
	image.FetchableFiles = inherit.StringMapTo(ctx, data.FetchableFiles, diags)
	image.BootFiles = inherit.StringMapTo(ctx, data.BootFiles, diags)
	image.MgmtClasses = inherit.StringListTo(ctx, data.MgmtClasses, diags)
	image.Owners = inherit.StringListTo(ctx, data.Owners, diags)

	// ElementsAs fails on null/unknown values; guard to avoid plan-time errors
	// for Optional+Computed fields not set in the configuration.
	var bootLoaders []string
	if !data.BootLoaders.IsNull() && !data.BootLoaders.IsUnknown() {
		diags.Append(data.BootLoaders.ElementsAs(ctx, &bootLoaders, false)...)
	}
	image.BootLoaders = bootLoaders

	var templateFiles map[string]interface{}
	if !data.TemplateFiles.IsNull() && !data.TemplateFiles.IsUnknown() {
		diags.Append(data.TemplateFiles.ElementsAs(ctx, &templateFiles, false)...)
	}
	image.TemplateFiles = cobbler.Value[map[string]interface{}]{Data: templateFiles, IsInherited: false}
	return image
}

// imageToModel populates an imageResourceModel from a cobbler.Image.
func imageToModel(ctx context.Context, image cobbler.Image, data *imageResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(image.Name)
	data.File = types.StringValue(image.File)
	data.Arch = types.StringValue(image.Arch)
	data.Autoinstall = types.StringValue(image.Autoinstall)
	data.Breed = types.StringValue(image.Breed)
	data.Comment = types.StringValue(image.Comment)
	data.ImageType = types.StringValue(image.ImageType)
	data.OSVersion = types.StringValue(image.OsVersion)
	data.Menu = types.StringValue(image.Menu)
	data.VirtAutoBoot = types.BoolValue(image.VirtAutoBoot)
	data.VirtBridge = types.StringValue(image.VirtBridge)
	data.VirtCpus = types.Int64Value(int64(image.VirtCpus))
	data.VirtDiskDriver = types.StringValue(image.VirtDiskDriver)
	data.VirtFileSize = inherit.Float64From(ctx, image.VirtFileSize, diags)
	data.VirtPath = types.StringValue(image.VirtPath)
	data.VirtRam = inherit.IntFrom(ctx, image.VirtRam, diags)
	data.VirtType = types.StringValue(image.VirtType)
	data.KernelOptions = inherit.StringMapFrom(ctx, image.KernelOptions, diags)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, image.KernelOptionsPost, diags)
	data.FetchableFiles = inherit.StringMapFrom(ctx, image.FetchableFiles, diags)
	data.BootFiles = inherit.StringMapFrom(ctx, image.BootFiles, diags)
	data.MgmtClasses = inherit.StringListFrom(ctx, image.MgmtClasses, diags)
	data.Owners = inherit.StringListFrom(ctx, image.Owners, diags)

	bootLoaders, d := types.ListValueFrom(ctx, types.StringType, image.BootLoaders)
	diags.Append(d...)
	data.BootLoaders = bootLoaders

	templateFiles, d := types.MapValueFrom(ctx, types.StringType, image.TemplateFiles.Data)
	diags.Append(d...)
	data.TemplateFiles = templateFiles
}
