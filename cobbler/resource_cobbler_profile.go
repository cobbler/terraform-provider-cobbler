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
			"boot_files": {
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
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
			},
			"distro": {
				Description: "Parent distribution.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enable_gpxe": {
				Description: "Use gPXE instead of PXELINUX for advanced booting options.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enable_menu": {
				Description: "Enable a boot menu.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"fetchable_files": {
				Description: "Templates for tftp or wget.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options": {
				Description: "Kernel options for the profile.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_post": {
				Description: "Post install kernel options.",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"mgmt_classes": {
				Description: "For external configuration management.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mgmt_parameters": {
				Description: "Parameters which will be handed to your management application (Must be a valid YAML dictionary).",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
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
			"name_servers": {
				Description: "Name servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"virt_auto_boot": {
				Description: "Auto boot virtual machines.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
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
	err = d.Set("autoinstall_meta", profile.AutoinstallMeta.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("boot_files", profile.BootFiles.Data)
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
	// TODO: enable_ipxe
	err = d.Set("enable_gpxe", profile.EnableIPXE.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("enable_menu", profile.EnableMenu.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("fetchable_files", profile.FetchableFiles.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options", profile.KernelOptions.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options_post", profile.KernelOptionsPost.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mgmt_classes", profile.MgmtClasses.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	if profile.MgmtParameters.IsInherited {
		err = d.Set("mgmt_parameters", make(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = d.Set("mgmt_parameters", profile.MgmtParameters.Data)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = d.Set("name", profile.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name_servers_search", profile.NameServersSearch.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name_servers", profile.NameServers.Data)
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
	err = d.Set("owners", profile.Owners.Data)
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
	err = d.Set("template_files", profile.TemplateFiles.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_auto_boot", profile.VirtAutoBoot.Data)
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
	if !profile.VirtFileSize.IsInherited {
		err = d.Set("virt_file_size", profile.VirtFileSize.Data)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = d.Set("virt_path", profile.VirtPath)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_ram", profile.VirtRAM.Data)
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
		IsInherited: false,
	}
	profile.BootFiles = cobbler.Value[map[string]interface{}]{
		Data:        bootFiles,
		IsInherited: false,
	}
	profile.Comment = d.Get("comment").(string)
	profile.DHCPTag = d.Get("dhcp_tag").(string)
	profile.Distro = d.Get("distro").(string)
	// TODO: enable_ipxe
	profile.EnableIPXE = cobbler.Value[bool]{
		Data:        d.Get("enable_gpxe").(bool),
		IsInherited: false,
	}
	profile.EnableMenu = cobbler.Value[bool]{
		Data:        d.Get("enable_menu").(bool),
		IsInherited: false,
	}
	profile.FetchableFiles = cobbler.Value[map[string]interface{}]{
		Data:        fetchableFiles,
		IsInherited: false,
	}
	profile.KernelOptions = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptions,
		IsInherited: false,
	}
	profile.KernelOptionsPost = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptionsPost,
		IsInherited: false,
	}
	profile.MgmtClasses = cobbler.Value[[]string]{
		Data:        mgmtClasses,
		IsInherited: false,
	}
	profile.MgmtParameters = cobbler.Value[map[string]interface{}]{
		Data:        mgmtParameters,
		IsInherited: false,
	}
	profile.Name = d.Get("name").(string)
	profile.NameServersSearch = cobbler.Value[[]string]{
		Data:        nameServersSearch,
		IsInherited: false,
	}
	profile.NameServers = cobbler.Value[[]string]{
		Data:        nameServers,
		IsInherited: false,
	}
	profile.NextServerv4 = d.Get("next_server_v4").(string)
	profile.NextServerv6 = d.Get("next_server_v6").(string)
	profile.Owners = cobbler.Value[[]string]{
		Data:        owners,
		IsInherited: false,
	}
	profile.Proxy = d.Get("proxy").(string)
	profile.Repos = repos
	profile.Server = d.Get("server").(string)
	profile.TemplateFiles = cobbler.Value[map[string]interface{}]{
		Data:        templateFiles,
		IsInherited: false,
	}
	profile.VirtAutoBoot = cobbler.Value[bool]{
		Data:        d.Get("virt_auto_boot").(bool),
		IsInherited: false,
	}
	profile.VirtBridge = d.Get("virt_bridge").(string)
	profile.VirtCPUs = d.Get("virt_cpus").(int)
	profile.VirtDiskDriver = d.Get("virt_disk_driver").(string)
	profile.VirtFileSize = cobbler.Value[float64]{
		Data:        d.Get("virt_file_size").(float64),
		IsInherited: false,
	}
	profile.VirtPath = d.Get("virt_path").(string)
	profile.VirtRAM = cobbler.Value[int]{
		Data:        d.Get("virt_ram").(int),
		IsInherited: false,
	}
	profile.VirtType = d.Get("virt_type").(string)

	return profile, nil
}
