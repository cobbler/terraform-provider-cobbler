package cobbler

import (
	"fmt"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`cobbler_profile` manages a profile within Cobbler.",
		Create:      resourceProfileCreate,
		Read:        resourceProfileRead,
		Update:      resourceProfileUpdate,
		Delete:      resourceProfileDelete,

		Schema: map[string]*schema.Schema{
			"autoinstall": {
				Description: "Template remote kickstarts or preseeds.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"autoinstall_meta": {
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"boot_files": {
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Type:        schema.TypeList,
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options": {
				Description: "Kernel options for the profile.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_post": {
				Description: "Post install kernel options.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Type:        schema.TypeString,
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"virt_auto_boot": {
				Description: "Auto boot virtual machines.",
				Type:        schema.TypeString,
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
				Type:        schema.TypeString,
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
				Type:        schema.TypeString,
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
				Type:        schema.TypeString,
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

func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create a cobblerclient.Profile struct
	profile := buildProfile(d, config)

	// Attempt to create the Profile
	log.Printf("[DEBUG] Cobbler Profile: Create Options: %#v", profile)
	newProfile, err := config.cobblerClient.CreateProfile(profile)
	if err != nil {
		return fmt.Errorf("Cobbler Profile: Error Creating: %s", err)
	}

	d.SetId(newProfile.Name)

	return resourceProfileRead(d, meta)
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Retrieve the profile entry from Cobbler
	profile, err := config.cobblerClient.GetProfile(d.Id())
	log.Printf("[INFO] HELLO WORLD: %#v", profile)
	if err != nil {
		return fmt.Errorf("Cobbler Profile: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	d.Set("autoinstall", profile.Autoinstall)
	d.Set("autoinstall_meta", profile.AutoinstallMeta)
	d.Set("boot_files", profile.BootFiles)
	d.Set("comment", profile.Comment)
	d.Set("dhcp_tag", profile.DHCPTag)
	d.Set("distro", profile.Distro)
	d.Set("enable_gpxe", profile.EnableGPXE)
	d.Set("enable_menu", profile.EnableMenu)
	d.Set("fetchable_files", profile.FetchableFiles)
	d.Set("kernel_options", profile.KernelOptions)
	d.Set("kernel_options_post", profile.KernelOptionsPost)
	d.Set("mgmt_classes", profile.MGMTClasses)
	d.Set("mgmt_parameters", profile.MGMTParameters)
	d.Set("name", profile.Name)
	d.Set("name_servers_search", profile.NameServersSearch)
	d.Set("name_servers", profile.NameServers)
	d.Set("next_server_v4", profile.NextServerv4)
	d.Set("next_server_v6", profile.NextServerv6)
	d.Set("owners", profile.Owners)
	d.Set("proxy", profile.Proxy)
	d.Set("repos", profile.Repos)
	d.Set("template_files", profile.TemplateFiles)
	d.Set("virt_auto_boot", profile.VirtAutoBoot)
	d.Set("virt_bridge", profile.VirtBridge)
	d.Set("virt_cpus", profile.VirtCPUs)
	d.Set("virt_disk_driver", profile.VirtDiskDriver)
	d.Set("virt_file_size", profile.VirtFileSize)
	d.Set("virt_path", profile.VirtPath)
	d.Set("virt_ram", profile.VirtRAM)
	d.Set("virt_type", profile.VirtType)

	return nil
}

func resourceProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create a cobblerclient.Profile struct
	profile := buildProfile(d, config)

	// Attempt to update the profile with new information
	log.Printf("[INFO] HELLO WORLD: %#v", profile)
	log.Printf("[DEBUG] Cobbler Profile: Updating Profile (%s) with options: %+v", d.Id(), profile)
	err := config.cobblerClient.UpdateProfile(&profile)
	if err != nil {
		return fmt.Errorf("Cobbler Profile: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceProfileRead(d, meta)
}

func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Attempt to delete the profile
	if err := config.cobblerClient.DeleteProfile(d.Id()); err != nil {
		return fmt.Errorf("Cobbler Profile: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildProfile builds a cobblerclient.Profile out of the Terraform attributes.
func buildProfile(d *schema.ResourceData, meta interface{}) cobbler.Profile {
	mgmtClasses := []string{}
	for _, i := range d.Get("mgmt_classes").([]interface{}) {
		mgmtClasses = append(mgmtClasses, i.(string))
	}
	nameServersSearch := []string{}
	for _, i := range d.Get("name_servers_search").([]interface{}) {
		nameServersSearch = append(nameServersSearch, i.(string))
	}
	nameServers := []string{}
	for _, i := range d.Get("name_servers").([]interface{}) {
		nameServers = append(nameServers, i.(string))
	}
	owners := []string{}
	for _, i := range d.Get("owners").([]interface{}) {
		owners = append(owners, i.(string))
	}
	repos := []string{}
	for _, i := range d.Get("repos").([]interface{}) {
		repos = append(repos, i.(string))
	}
	bootFiles := []string{}
	for _, i := range d.Get("boot_files").([]interface{}) {
		bootFiles = append(bootFiles, i.(string))
	}
	fetchableFiles := []string{}
	for _, i := range d.Get("fetchable_files").([]interface{}) {
		fetchableFiles = append(fetchableFiles, i.(string))
	}
	kernelOptions := []string{}
	for _, i := range d.Get("kernel_options").([]interface{}) {
		kernelOptions = append(owners, i.(string))
	}
	kernelOptionsPost := []string{}
	for _, i := range d.Get("kernel_options_post").([]interface{}) {
		kernelOptionsPost = append(owners, i.(string))
	}
	templateFiles := []string{}
	for _, i := range d.Get("template_files").([]interface{}) {
		templateFiles = append(templateFiles, i.(string))
	}
	autoinstallMeta := []string{}
	for _, i := range d.Get("autoinstall_meta").([]interface{}) {
		autoinstallMeta = append(autoinstallMeta, i.(string))
	}

	profile := cobbler.Profile{
		Autoinstall:       d.Get("autoinstall").(string),
		AutoinstallMeta:   autoinstallMeta,
		BootFiles:         bootFiles,
		Comment:           d.Get("comment").(string),
		DHCPTag:           d.Get("dhcp_tag").(string),
		Distro:            d.Get("distro").(string),
		EnableGPXE:        d.Get("enable_gpxe").(bool),
		EnableMenu:        d.Get("enable_menu").(bool),
		FetchableFiles:    fetchableFiles,
		KernelOptions:     kernelOptions,
		KernelOptionsPost: kernelOptionsPost,
		MGMTClasses:       mgmtClasses,
		MGMTParameters:    d.Get("mgmt_parameters").(string),
		Name:              d.Get("name").(string),
		NameServersSearch: nameServersSearch,
		NameServers:       nameServers,
		NextServerv4:      d.Get("next_server_v4").(string),
		NextServerv6:      d.Get("next_server_v6").(string),
		Owners:            owners,
		Proxy:             d.Get("proxy").(string),
		Repos:             repos,
		Server:            d.Get("server").(string),
		TemplateFiles:     templateFiles,
		VirtAutoBoot:      d.Get("virt_auto_boot").(string),
		VirtBridge:        d.Get("virt_bridge").(string),
		VirtCPUs:          d.Get("virt_cpus").(string),
		VirtDiskDriver:    d.Get("virt_disk_driver").(string),
		VirtFileSize:      d.Get("virt_file_size").(string),
		VirtPath:          d.Get("virt_path").(string),
		VirtRAM:           d.Get("virt_ram").(string),
		VirtType:          d.Get("virt_type").(string),
	}

	return profile
}
