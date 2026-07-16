package image

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ImageDataSource{}

type ImageDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

func (d *ImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (d *ImageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler image.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the image.",
				Required:    true,
			},
			"uid": schema.StringAttribute{
				Description: "Server-assigned UID for this image. Use this as the value for `cobbler_system.image`.",
				Computed:    true,
			},
			"file": schema.StringAttribute{
				Description: "Path to the image media.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the image.",
				Computed:    true,
			},
			"autoinstall": schema.StringAttribute{
				Description: "Path to an autoinstall file.",
				Computed:    true,
			},
			"breed": schema.StringAttribute{
				Description: "The \"breed\" of distribution.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"image_type": schema.StringAttribute{
				Description: "Type of image (direct, iso, memdisk, virt-clone).",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "The OS version the image contains.",
				Computed:    true,
			},
			"boot_loaders": schema.ListAttribute{
				Description: "Boot loaders supported by the image.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"menu": schema.StringAttribute{
				Description: "The Cobbler UID of the parent menu.",
				Computed:    true,
			},
			"virt_auto_boot": schema.BoolAttribute{
				Description: "Whether to auto-boot the virtual machine.",
				Computed:    true,
			},
			"virt_bridge": schema.StringAttribute{
				Description: "Bridge for the virtual machine.",
				Computed:    true,
			},
			"virt_cpus": schema.Int64Attribute{
				Description: "Number of CPUs for the virtual machine.",
				Computed:    true,
			},
			"virt_disk_driver": schema.StringAttribute{
				Description: "Disk driver for the virtual machine.",
				Computed:    true,
			},
			"virt_file_size": schema.SingleNestedAttribute{
				Description: "Disk file size in GB for the virtual machine.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{
						Computed: true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"virt_path": schema.StringAttribute{
				Description: "Path on the virtualization host.",
				Computed:    true,
			},
			"virt_ram": schema.SingleNestedAttribute{
				Description: "RAM in MB for the virtual machine.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{
						Computed: true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"virt_type": schema.StringAttribute{
				Description: "Virtualization type.",
				Computed:    true,
			},
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options to use with the kernel.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"kernel_options_post": schema.SingleNestedAttribute{
				Description: "Post install kernel options.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "Owners list for authz_ownership.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
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

func (d *ImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*clientpkg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			"Expected *client.Config, got unexpected type.")
		return
	}
	d.client = cfg.CobblerClient
}

func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data imageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imagePtr, err := d.client.GetImage(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Image", err.Error())
		return
	}
	image := *imagePtr

	data.Name = types.StringValue(image.Name)
	data.UID = types.StringValue(image.Uid)
	data.File = types.StringValue(image.File)
	data.Arch = types.StringValue(image.Arch)
	data.Autoinstall = types.StringValue(image.Autoinstall)
	data.Breed = types.StringValue(image.Breed)
	data.Comment = types.StringValue(image.Comment)
	data.ImageType = types.StringValue(image.ImageType)
	data.OSVersion = types.StringValue(image.OsVersion)
	data.Menu = types.StringValue(image.Menu)
	data.VirtAutoBoot = types.BoolValue(image.Virt.AutoBoot.Data)
	data.VirtBridge = types.StringValue(image.VirtBridge)
	data.VirtCpus = types.Int64Value(int64(image.Virt.Cpus.Data))
	data.VirtDiskDriver = types.StringValue(image.Virt.DiskDriver)
	data.VirtFileSize = inherit.Float64From(ctx, image.Virt.FileSize, &resp.Diagnostics)
	data.VirtPath = types.StringValue(image.Virt.Path)
	data.VirtRam = inherit.IntFrom(ctx, image.Virt.Ram, &resp.Diagnostics)
	data.VirtType = types.StringValue(image.Virt.Type)
	data.KernelOptions = inherit.StringMapFrom(ctx, image.KernelOptions, &resp.Diagnostics)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, image.KernelOptionsPost, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, image.Owners, &resp.Diagnostics)

	bootLoaders, d2 := types.ListValueFrom(ctx, types.StringType, image.BootLoaders)
	resp.Diagnostics.Append(d2...)
	data.BootLoaders = bootLoaders

	templateFiles, d3 := types.MapValueFrom(ctx, types.StringType, image.TemplateFiles)
	resp.Diagnostics.Append(d3...)
	data.TemplateFiles = templateFiles

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
