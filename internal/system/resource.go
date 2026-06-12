package system

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SystemResource{}
var _ resource.ResourceWithImportState = &SystemResource{}

type SystemResource struct {
	client cobbler.Client
}

func NewResource() resource.Resource {
	return &SystemResource{}
}

func (r *SystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (r *SystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`cobbler_system` manages a system within Cobbler.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the system.",
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
			"gateway": schema.StringAttribute{
				Description: "Network gateway.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname of the system.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image": schema.StringAttribute{
				Description: "Parent image (if no profile is used).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv6_default_device": schema.StringAttribute{
				Description: "IPv6 default device.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name_servers": schema.ListAttribute{
				Description: "Name servers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"name_servers_search": schema.ListAttribute{
				Description: "Name server search settings.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"netboot_enabled": schema.BoolAttribute{
				Description: "(Re)install this machine at next boot.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
			"power_address": schema.StringAttribute{
				Description: "Power management address.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"power_id": schema.StringAttribute{
				Description: "Usually a plug number or blade name if power type requires it.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"power_pass": schema.StringAttribute{
				Description: "Power management password.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"power_type": schema.StringAttribute{
				Description: "Power management type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"power_user": schema.StringAttribute{
				Description: "Power management user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "Parent profile.",
				Required:    true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "System status (development, testing, acceptance, production).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"virt_pxe_boot": schema.BoolAttribute{
				Description: "Use PXE to build this virtual machine.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
			"interface": InterfaceMapAttribute(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
			"boot_loaders": schema.SingleNestedAttribute{
				Description: "Must be either `grub`, `pxe`, or `ipxe`.",
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
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
			"kernel_options": schema.SingleNestedAttribute{
				Description: "Kernel options for the system.",
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
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
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
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
			"virt_cpus": schema.SingleNestedAttribute{
				Description: "The number of virtual CPUs.",
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

func (r *SystemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	var data systemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	system := modelToSystem(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler System: Create", map[string]interface{}{"name": system.Name})

	newSystem, err := r.client.CreateSystem(system)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Cobbler System", err.Error())
		return
	}

	// Build and attach interfaces
	planIfacesAPI := InterfaceMapToAPI(ctx, data.Interface, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	for name, iface := range planIfacesAPI {
		if err := newSystem.CreateInterface(name, iface); err != nil {
			resp.Diagnostics.AddError("Error creating interface",
				"Error adding interface "+name+" to system "+newSystem.Name+": "+err.Error())
			return
		}
	}

	tflog.Debug(ctx, "Cobbler System: syncing system")
	if err := r.client.Sync(); err != nil {
		resp.Diagnostics.AddError("Error syncing Cobbler", err.Error())
		return
	}

	// Read back the system to get computed values
	readSystem, err := r.client.GetSystem(newSystem.Name, false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler System after create", err.Error())
		return
	}

	ifaces, err := readSystem.GetInterfaces()
	if err != nil {
		resp.Diagnostics.AddError("Error getting interfaces after create", err.Error())
		return
	}

	systemToModel(ctx, *readSystem, ifaces, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data systemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := r.client.GetSystem(data.Name.ValueString(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading System", err.Error())
		return
	}

	ifaces, err := system.GetInterfaces()
	if err != nil {
		resp.Diagnostics.AddError("Error getting interfaces", err.Error())
		return
	}

	systemToModel(ctx, *system, ifaces, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	var plan systemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state systemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the existing system to perform interface operations on it
	system, err := r.client.GetSystem(plan.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler System", err.Error())
		return
	}

	newSystem := modelToSystem(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler System: Update", map[string]interface{}{"name": newSystem.Name})

	if err := r.client.UpdateSystem(&newSystem); err != nil {
		resp.Diagnostics.AddError("Error updating Cobbler System", err.Error())
		return
	}

	// Interface diff: delete removed interfaces, create/update all plan interfaces
	var planIfaces map[string]types.Object
	resp.Diagnostics.Append(plan.Interface.ElementsAs(ctx, &planIfaces, false)...)

	var stateIfaces map[string]types.Object
	resp.Diagnostics.Append(state.Interface.ElementsAs(ctx, &stateIfaces, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete interfaces that exist in state but not in plan
	for name := range stateIfaces {
		if _, exists := planIfaces[name]; !exists {
			if err := system.DeleteInterface(name); err != nil {
				resp.Diagnostics.AddError("Error deleting interface", err.Error())
				return
			}
		}
	}

	// Create/update all interfaces from plan
	planIfacesAPI := InterfaceMapToAPI(ctx, plan.Interface, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	for name, iface := range planIfacesAPI {
		if err := system.CreateInterface(name, iface); err != nil {
			resp.Diagnostics.AddError("Error creating interface", err.Error())
			return
		}
	}

	tflog.Debug(ctx, "Cobbler System: syncing system")
	if err := r.client.Sync(); err != nil {
		resp.Diagnostics.AddError("Error syncing Cobbler", err.Error())
		return
	}

	// Read back updated system
	readSystem, err := r.client.GetSystem(plan.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler System after update", err.Error())
		return
	}

	updatedIfaces, err := readSystem.GetInterfaces()
	if err != nil {
		resp.Diagnostics.AddError("Error getting interfaces after update", err.Error())
		return
	}

	systemToModel(ctx, *readSystem, updatedIfaces, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data systemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cobbler System: Delete", map[string]interface{}{"name": data.Name.ValueString()})

	if err := r.client.DeleteSystem(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Cobbler System", err.Error())
	}
}

func (r *SystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// systemStringOrInherit returns "<<inherit>>" when s is null, unknown, or empty.
// Cobbler rejects empty strings for enum-validated fields (e.g. virt_disk_driver, virt_type).
func systemStringOrInherit(s types.String) string {
	if v := s.ValueString(); v != "" {
		return v
	}
	return "<<inherit>>"
}

// modelToSystem converts a systemResourceModel to a cobbler.System.
func modelToSystem(ctx context.Context, data systemResourceModel, diags *diag.Diagnostics) cobbler.System {
	system := cobbler.NewSystem()
	system.Name = data.Name.ValueString()
	system.Autoinstall = data.Autoinstall.ValueString()
	system.Comment = data.Comment.ValueString()
	system.Gateway = data.Gateway.ValueString()
	system.Hostname = data.Hostname.ValueString()
	system.Image = data.Image.ValueString()
	system.IPv6DefaultDevice = data.IPv6DefaultDevice.ValueString()
	system.NetbootEnabled = data.NetbootEnabled.ValueBool()
	system.NextServerv4 = data.NextServerV4.ValueString()
	system.NextServerv6 = data.NextServerV6.ValueString()
	system.PowerAddress = data.PowerAddress.ValueString()
	system.PowerID = data.PowerID.ValueString()
	system.PowerPass = data.PowerPass.ValueString()
	system.PowerType = data.PowerType.ValueString()
	system.PowerUser = data.PowerUser.ValueString()
	system.Profile = data.Profile.ValueString()
	system.Proxy = data.Proxy.ValueString()
	system.Status = data.Status.ValueString()
	system.VirtDiskDriver = systemStringOrInherit(data.VirtDiskDriver)
	system.VirtPath = data.VirtPath.ValueString()
	system.VirtPXEBoot = data.VirtPXEBoot.ValueBool()
	system.VirtType = systemStringOrInherit(data.VirtType)

	var nameServers []string
	if !data.NameServers.IsNull() && !data.NameServers.IsUnknown() {
		diags.Append(data.NameServers.ElementsAs(ctx, &nameServers, false)...)
	}
	system.NameServers = nameServers

	var nameServersSearch []string
	if !data.NameServersSearch.IsNull() && !data.NameServersSearch.IsUnknown() {
		diags.Append(data.NameServersSearch.ElementsAs(ctx, &nameServersSearch, false)...)
	}
	system.NameServersSearch = nameServersSearch

	system.AutoinstallMeta = inherit.StringMapTo(ctx, data.AutoinstallMeta, diags)
	system.BootFiles = inherit.StringMapTo(ctx, data.BootFiles, diags)
	system.BootLoaders = inherit.StringListTo(ctx, data.BootLoaders, diags)
	system.EnableIPXE = inherit.BoolTo(ctx, data.EnableIPXE, diags)
	system.FetchableFiles = inherit.StringMapTo(ctx, data.FetchableFiles, diags)
	system.KernelOptions = inherit.StringMapTo(ctx, data.KernelOptions, diags)
	system.KernelOptionsPost = inherit.StringMapTo(ctx, data.KernelOptionsPost, diags)
	system.MgmtClasses = inherit.StringListTo(ctx, data.MgmtClasses, diags)
	system.MgmtParameters = inherit.StringMapTo(ctx, data.MgmtParameters, diags)
	system.Owners = inherit.StringListTo(ctx, data.Owners, diags)
	system.TemplateFiles = inherit.StringMapTo(ctx, data.TemplateFiles, diags)
	system.VirtAutoBoot = inherit.BoolTo(ctx, data.VirtAutoBoot, diags)
	system.VirtCPUs = inherit.IntTo(ctx, data.VirtCPUs, diags)
	system.VirtFileSize = inherit.Float64To(ctx, data.VirtFileSize, diags)
	system.VirtRAM = inherit.IntTo(ctx, data.VirtRAM, diags)

	return system
}

// systemToModel populates a systemResourceModel from a cobbler.System and interfaces.
func systemToModel(ctx context.Context, system cobbler.System, ifaces cobbler.Interfaces, data *systemResourceModel, diags *diag.Diagnostics) {
	data.Name = types.StringValue(system.Name)
	data.Autoinstall = types.StringValue(system.Autoinstall)
	data.Comment = types.StringValue(system.Comment)
	data.Gateway = types.StringValue(system.Gateway)
	data.Hostname = types.StringValue(system.Hostname)
	data.Image = types.StringValue(system.Image)
	data.IPv6DefaultDevice = types.StringValue(system.IPv6DefaultDevice)
	data.NetbootEnabled = types.BoolValue(system.NetbootEnabled)
	data.NextServerV4 = types.StringValue(system.NextServerv4)
	data.NextServerV6 = types.StringValue(system.NextServerv6)
	data.PowerAddress = types.StringValue(system.PowerAddress)
	data.PowerID = types.StringValue(system.PowerID)
	data.PowerPass = types.StringValue(system.PowerPass)
	data.PowerType = types.StringValue(system.PowerType)
	data.PowerUser = types.StringValue(system.PowerUser)
	data.Profile = types.StringValue(system.Profile)
	data.Proxy = types.StringValue(system.Proxy)
	data.Status = types.StringValue(system.Status)
	data.VirtDiskDriver = types.StringValue(system.VirtDiskDriver)
	data.VirtPath = types.StringValue(system.VirtPath)
	data.VirtPXEBoot = types.BoolValue(system.VirtPXEBoot)
	data.VirtType = types.StringValue(system.VirtType)

	nameServersList, d := types.ListValueFrom(ctx, types.StringType, system.NameServers)
	diags.Append(d...)
	data.NameServers = nameServersList

	nameServersSearchList, d := types.ListValueFrom(ctx, types.StringType, system.NameServersSearch)
	diags.Append(d...)
	data.NameServersSearch = nameServersSearchList

	data.Interface = InterfaceMapFromAPI(ctx, ifaces, diags)

	data.AutoinstallMeta = inherit.StringMapFrom(ctx, system.AutoinstallMeta, diags)
	data.BootFiles = inherit.StringMapFrom(ctx, system.BootFiles, diags)
	data.BootLoaders = inherit.StringListFrom(ctx, system.BootLoaders, diags)
	data.EnableIPXE = inherit.BoolFrom(ctx, system.EnableIPXE, diags)
	data.FetchableFiles = inherit.StringMapFrom(ctx, system.FetchableFiles, diags)
	data.KernelOptions = inherit.StringMapFrom(ctx, system.KernelOptions, diags)
	data.KernelOptionsPost = inherit.StringMapFrom(ctx, system.KernelOptionsPost, diags)
	data.MgmtClasses = inherit.StringListFrom(ctx, system.MgmtClasses, diags)
	data.MgmtParameters = inherit.StringMapFrom(ctx, system.MgmtParameters, diags)
	data.Owners = inherit.StringListFrom(ctx, system.Owners, diags)
	data.TemplateFiles = inherit.StringMapFrom(ctx, system.TemplateFiles, diags)
	data.VirtAutoBoot = inherit.BoolFrom(ctx, system.VirtAutoBoot, diags)
	data.VirtCPUs = inherit.IntFrom(ctx, system.VirtCPUs, diags)
	data.VirtFileSize = inherit.Float64From(ctx, system.VirtFileSize, diags)
	data.VirtRAM = inherit.IntFrom(ctx, system.VirtRAM, diags)
}
