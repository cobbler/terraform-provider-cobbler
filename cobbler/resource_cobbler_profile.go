package cobbler

import (
	"context"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_profile` manages a profile within Cobbler.",
		CreateContext: resourceProfileCreate,
		ReadContext:   resourceProfileRead,
		UpdateContext: resourceProfileUpdate,
		DeleteContext: resourceProfileDelete,

		Schema: map[string]*schema.Schema{
			"autoinstall": {
				Description: "Template remote kickstarts or preseeds.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"autoinstall_meta": {
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"autoinstall_meta_inherit": {
				Description:   "Signal that autoinstall_meta should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"autoinstall_meta"},
			},
			"boot_files": {
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"boot_files_inherit": {
				Description:   "Signal that boot_files should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"boot_files"},
			},
			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"dhcp_tag": {
				Description: "DHCP tag.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"distro": {
				Description: "Parent distribution.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enable_ipxe": {
				Description: "Use iPXE instead of PXELINUX for advanced booting options.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enable_ipxe_inherit": {
				Description:   "Signal that enable_ipxe should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enable_ipxe"},
			},
			"enable_menu": {
				Description: "Enable a boot menu.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enable_menu_inherit": {
				Description:   "Signal that enable_menu should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enable_menu"},
			},
			"fetchable_files": {
				Description: "Templates for tftp or wget.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"fetchable_files_inherit": {
				Description:   "Signal that fetchable_files should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"fetchable_files"},
			},
			"kernel_options": {
				Description: "Kernel options for the profile.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_inherit": {
				Description:   "Signal that kernel_options should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"kernel_options"},
			},
			"kernel_options_post": {
				Description: "Post install kernel options.",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_post_inherit": {
				Description:   "Signal that kernel_options_post should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"kernel_options_post"},
			},
			"mgmt_classes": {
				Description: "For external configuration management.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mgmt_classes_inherit": {
				Description:   "Signal that mgmt_classes should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"mgmt_classes"},
			},
			"mgmt_parameters": {
				Description: "Parameters which will be handed to your management application (Must be a valid YAML dictionary).",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"mgmt_parameters_inherit": {
				Description:   "Signal that mgmt_parameters should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"mgmt_parameters"},
			},
			"name": {
				Description: "The name of the profile.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name_servers_search": {
				Description: "Name server search settings.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name_servers_search_inherit": {
				Description:   "Signal that name_servers_search should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name_servers_search"},
			},
			"name_servers": {
				Description: "Name servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name_servers_inherit": {
				Description:   "Signal that name_servers should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name_servers"},
			},
			"next_server_v4": {
				Description: "The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"next_server_v6": {
				Description: "The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"owners": {
				Description: "Owners list for authz_ownership.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"owners_inherit": {
				Description:   "Signal that owners should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"owners"},
			},
			"parent": {
				Description: "The parent this profile inherits settings from.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"proxy": {
				Description: "Proxy URL.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"repos": {
				Description: "Repos to auto-assign to this profile.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"server": {
				Description: "The server-override for the profile.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"template_files": {
				Description: "File mappings for built-in config management.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"template_files_inherit": {
				Description:   "Signal that template_files should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"template_files"},
			},
			"virt_auto_boot": {
				Description: "Auto boot virtual machines.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"virt_auto_boot_inherit": {
				Description:   "Signal that virt_auto_boot should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"virt_auto_boot"},
			},
			"virt_bridge": {
				Description: "The bridge for virtual machines.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"virt_cpus": {
				Description: "The number of virtual CPUs",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"virt_disk_driver": {
				Description: "The virtual machine disk driver.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"virt_file_size": {
				Description: "The virtual machine file size.",
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
			},
			"virt_file_size_inherit": {
				Description:   "Signal that virt_file_size should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"virt_file_size"},
			},
			"virt_path": {
				Description: "The virtual machine path.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"virt_ram": {
				Description: "The amount of RAM for the virtual machine.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"virt_ram_inherit": {
				Description:   "Signal that virt_ram should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"virt_ram"},
			},
			"virt_type": {
				Description: "The type of virtual machine. Valid options are: xenpv, xenfv, qemu, kvm, vmware, openvz.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Create a cobblerclient.Profile struct
	profile, err := buildProfile(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to create the Profile
	tflog.Debug(ctx, "Cobbler Profile: Create Options", map[string]interface{}{
		"options": structs.Map(profile),
	})
	newProfile, err := config.cobblerClient.CreateProfile(profile)
	if err != nil {
		return diag.Errorf("Cobbler Profile: Error Creating: %s", err)
	}

	d.SetId(newProfile.Name)

	return resourceProfileRead(ctx, d, meta)
}

func resourceProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Retrieve the profile entry from Cobbler
	profile, err := config.cobblerClient.GetProfile(d.Id(), false, false)
	if err != nil {
		return diag.Errorf("Cobbler Profile: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	err = d.Set("autoinstall", profile.Autoinstall)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "autoinstall_meta", profile.AutoinstallMeta, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "boot_files", profile.BootFiles, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("comment", profile.Comment)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("dhcp_tag", profile.DHCPTag)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("distro", profile.Distro)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "enable_ipxe", profile.EnableIPXE, false)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "enable_menu", profile.EnableMenu, false)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "fetchable_files", profile.FetchableFiles, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "kernel_options", profile.KernelOptions, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "kernel_options_post", profile.KernelOptionsPost, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "mgmt_classes", profile.MgmtClasses, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "mgmt_parameters", profile.MgmtParameters, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name", profile.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "name_servers_search", profile.NameServersSearch, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "name_servers", profile.NameServers, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("next_server_v4", profile.NextServerv4)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("next_server_v6", profile.NextServerv6)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "owners", profile.Owners, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("proxy", profile.Proxy)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("repos", profile.Repos)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "template_files", profile.TemplateFiles, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "virt_auto_boot", profile.VirtAutoBoot, false)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_bridge", profile.VirtBridge)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_cpus", profile.VirtCPUs)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_disk_driver", profile.VirtDiskDriver)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "virt_file_size", profile.VirtFileSize, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_path", profile.VirtPath)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "virt_ram", profile.VirtRAM, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_type", profile.VirtType)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Create a cobblerclient.Profile struct
	profile, err := buildProfile(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to update the profile with new information
	tflog.Debug(ctx, "Cobbler Profile: Updating Profile with options", map[string]interface{}{
		"profile": d.Id(),
		"options": structs.Map(profile),
	})
	err = config.cobblerClient.UpdateProfile(&profile)
	if err != nil {
		return diag.Errorf("error updating Cobbler Profile: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceProfileRead(ctx, d, meta)
}

func resourceProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Attempt to delete the profile
	if err := config.cobblerClient.DeleteProfile(d.Id()); err != nil {
		return diag.Errorf("Cobbler Profile: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildProfile builds a cobblerclient.Profile out of the Terraform attributes.
func buildProfile(d *schema.ResourceData, meta interface{}) (cobbler.Profile, error) { //nolint:unparam // We satisfy our own pattern here
	mgmtClasses, err := GetStringSlice(d, "mgmt_classes")
	if err != nil {
		return cobbler.Profile{}, err
	}
	mgmtParameters, err := GetInterfaceMap(d, "mgmt_parameters")
	if err != nil {
		return cobbler.Profile{}, err
	}
	nameServersSearch, err := GetStringSlice(d, "name_servers_search")
	if err != nil {
		return cobbler.Profile{}, err
	}
	nameServers, err := GetStringSlice(d, "name_servers")
	if err != nil {
		return cobbler.Profile{}, err
	}
	owners, err := GetStringSlice(d, "owners")
	if err != nil {
		return cobbler.Profile{}, err
	}
	repos, err := GetStringSlice(d, "repos")
	if err != nil {
		return cobbler.Profile{}, err
	}
	bootFiles, err := GetInterfaceMap(d, "boot_files")
	if err != nil {
		return cobbler.Profile{}, err
	}
	fetchableFiles, err := GetInterfaceMap(d, "fetchable_files")
	if err != nil {
		return cobbler.Profile{}, err
	}
	kernelOptions, err := GetInterfaceMap(d, "kernel_options")
	if err != nil {
		return cobbler.Profile{}, err
	}
	kernelOptionsPost, err := GetInterfaceMap(d, "kernel_options_post")
	if err != nil {
		return cobbler.Profile{}, err
	}
	templateFiles, err := GetInterfaceMap(d, "template_files")
	if err != nil {
		return cobbler.Profile{}, err
	}
	autoinstallMeta, err := GetInterfaceMap(d, "autoinstall_meta")
	if err != nil {
		return cobbler.Profile{}, err
	}

	profile := cobbler.NewProfile()
	profile.Autoinstall = d.Get("autoinstall").(string)
	profile.AutoinstallMeta = cobbler.Value[map[string]interface{}]{
		Data:        autoinstallMeta,
		IsInherited: IsOptionInherited(d, "autoinstall_meta"),
	}
	profile.BootFiles = cobbler.Value[map[string]interface{}]{
		Data:        bootFiles,
		IsInherited: IsOptionInherited(d, "boot_files"),
	}
	profile.Comment = d.Get("comment").(string)
	profile.DHCPTag = d.Get("dhcp_tag").(string)
	profile.Distro = d.Get("distro").(string)
	profile.EnableIPXE = cobbler.Value[bool]{
		Data:        d.Get("enable_ipxe").(bool),
		IsInherited: IsOptionInherited(d, "enable_ipxe"),
	}
	profile.EnableMenu = cobbler.Value[bool]{
		Data:        d.Get("enable_menu").(bool),
		IsInherited: IsOptionInherited(d, "enable_menu"),
	}
	profile.FetchableFiles = cobbler.Value[map[string]interface{}]{
		Data:        fetchableFiles,
		IsInherited: IsOptionInherited(d, "fetchable_files"),
	}
	profile.KernelOptions = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptions,
		IsInherited: IsOptionInherited(d, "kernel_options"),
	}
	profile.KernelOptionsPost = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptionsPost,
		IsInherited: IsOptionInherited(d, "kernel_options_post"),
	}
	profile.MgmtClasses = cobbler.Value[[]string]{
		Data:        mgmtClasses,
		IsInherited: IsOptionInherited(d, "mgmt_classes"),
	}
	profile.MgmtParameters = cobbler.Value[map[string]interface{}]{
		Data:        mgmtParameters,
		IsInherited: IsOptionInherited(d, "mgmt_parameters"),
	}
	profile.Name = d.Get("name").(string)
	profile.NameServersSearch = cobbler.Value[[]string]{
		Data:        nameServersSearch,
		IsInherited: IsOptionInherited(d, "name_servers_search"),
	}
	profile.NameServers = cobbler.Value[[]string]{
		Data:        nameServers,
		IsInherited: IsOptionInherited(d, "name_servers"),
	}
	profile.NextServerv4 = d.Get("next_server_v4").(string)
	profile.NextServerv6 = d.Get("next_server_v6").(string)
	profile.Owners = cobbler.Value[[]string]{
		Data:        owners,
		IsInherited: IsOptionInherited(d, "owners"),
	}
	profile.Proxy = d.Get("proxy").(string)
	profile.Repos = repos
	profile.Server = d.Get("server").(string)
	profile.TemplateFiles = cobbler.Value[map[string]interface{}]{
		Data:        templateFiles,
		IsInherited: IsOptionInherited(d, "template_files"),
	}
	profile.VirtAutoBoot = cobbler.Value[bool]{
		Data:        d.Get("virt_auto_boot").(bool),
		IsInherited: IsOptionInherited(d, "virt_auto_boot"),
	}
	profile.VirtBridge = d.Get("virt_bridge").(string)
	profile.VirtCPUs = d.Get("virt_cpus").(int)
	profile.VirtDiskDriver = d.Get("virt_disk_driver").(string)
	profile.VirtFileSize = cobbler.Value[float64]{
		Data:        d.Get("virt_file_size").(float64),
		IsInherited: IsOptionInherited(d, "virt_file_size"),
	}
	profile.VirtPath = d.Get("virt_path").(string)
	profile.VirtRAM = cobbler.Value[int]{
		Data:        d.Get("virt_ram").(int),
		IsInherited: IsOptionInherited(d, "virt_ram"),
	}
	profile.VirtType = d.Get("virt_type").(string)

	return profile, nil
}
