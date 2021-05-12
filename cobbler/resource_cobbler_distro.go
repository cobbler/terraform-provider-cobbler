package cobbler

import (
	"fmt"
	"log"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDistro() *schema.Resource {
	return &schema.Resource{
		Description: "`cobbler_distro` manages a distribution within Cobbler.",
		Create:      resourceDistroCreate,
		Read:        resourceDistroRead,
		Update:      resourceDistroUpdate,
		Delete:      resourceDistroDelete,

		Schema: map[string]*schema.Schema{
			"arch": {
				Description: "The architecture of the distro. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},

			"breed": {
				Description: "The \"breed\" of distribution. Valid options are: redhat, fedora, centos, scientific linux, suse, debian, and ubuntu. These choices may vary depending on the version of Cobbler in use.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"boot_files": {
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},

			"boot_loader": {
				Description: " Must be either `grub` or `pxelinux`.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},

			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"fetchable_files": {
				Description: "Templates for tftp or wget.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"kernel": {
				Description: "Absolute path to kernel on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"kernel_options": {
				Description: "Kernel options to use with the kernel.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"kernel_options_post": {
				Description: "Post install Kernel options to use with the kernel after installation.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"initrd": {
				Description: "Absolute path to initrd on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"mgmt_classes": {
				Description: "Management classes for external config management.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},

			"name": {
				Description: "A name for the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"os_version": {
				Description: "The version of the distro you are creating. This varies with the version of Cobbler you are using. An updated signature list may need to be obtained in order to support a newer version. Example: `bionic`.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"owners": {
				Description: "Owners list for authz_ownership.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},

			"redhat_management_key": {
				Description: "Obsolete - removed in Cobbler 3.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"template_files": {
				Description: "File mappings for built-in config management.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceDistroCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create a cobblerclient.Distro
	distro := buildDistro(d, config)

	// Attempte to create the Distro
	log.Printf("[DEBUG] Cobbler Distro: Create Options: %#v", distro)
	newDistro, err := config.cobblerClient.CreateDistro(distro)
	if err != nil {
		return fmt.Errorf("Cobbler Distro: Error Creating: %s", err)
	}

	d.SetId(newDistro.Name)

	return resourceDistroRead(d, meta)
}

func resourceDistroRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Retrieve the distro from cobbler
	distro, err := config.cobblerClient.GetDistro(d.Id())
	if err != nil {
		return fmt.Errorf("Cobbler Distro: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	d.Set("arch", distro.Arch)
	d.Set("breed", distro.Breed)
	d.Set("boot_files", distro.BootFiles)
	d.Set("boot_loader", distro.BootLoader)
	d.Set("comment", distro.Comment)
	d.Set("fetchable_files", distro.FetchableFiles)
	d.Set("kernel", distro.Kernel)
	d.Set("kernel_options", distro.KernelOptions)
	d.Set("kernel_options_post", distro.KernelOptionsPost)
	d.Set("initrd", distro.Initrd)
	d.Set("mgmt_classes", distro.MGMTClasses)
	d.Set("os_version", distro.OSVersion)
	d.Set("owners", distro.Owners)
	d.Set("redhat_management_key", distro.RedHatManagementKey)
	//d.Set("redhat_management_server", distro.RedHatManagementServer)         // Removed in Cobbler 3
	d.Set("template_files", distro.TemplateFiles)

	return nil
}

func resourceDistroUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// create a cobblerclient.Distro
	distro := buildDistro(d, config)

	// Attempt to updateh the distro with new information
	log.Printf("[DEBUG] Cobbler Distro: Updating Distro (%s) with options: %+v", d.Id(), distro)
	err := config.cobblerClient.UpdateDistro(&distro)
	if err != nil {
		return fmt.Errorf("Cobbler Distro: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceDistroRead(d, meta)
}

func resourceDistroDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Attempt to delete the distro
	if err := config.cobblerClient.DeleteDistro(d.Id()); err != nil {
		return fmt.Errorf("Cobbler Distro: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildDistro builds a cobbler.Distro from the Terraform attributes.
func buildDistro(d *schema.ResourceData, meta interface{}) cobbler.Distro {
	mgmtClasses := []string{}
	for _, i := range d.Get("mgmt_classes").([]interface{}) {
		mgmtClasses = append(mgmtClasses, i.(string))
	}

	owners := []string{}
	for _, i := range d.Get("owners").([]interface{}) {
		owners = append(owners, i.(string))
	}

	distro := cobbler.Distro{
		Arch:                d.Get("arch").(string),
		Breed:               d.Get("breed").(string),
		BootFiles:           d.Get("boot_files").(string),
		BootLoader:          d.Get("boot_loader").(string),
		Comment:             d.Get("comment").(string),
		FetchableFiles:      d.Get("fetchable_files").(string),
		Kernel:              d.Get("kernel").(string),
		KernelOptions:       d.Get("kernel_options").(string),
		KernelOptionsPost:   d.Get("kernel_options_post").(string),
		Initrd:              d.Get("initrd").(string),
		MGMTClasses:         mgmtClasses,
		Name:                d.Get("name").(string),
		OSVersion:           d.Get("os_version").(string),
		Owners:              owners,
		RedHatManagementKey: d.Get("redhat_management_key").(string),
		//RedHatManagementServer: d.Get("redhat_management_server").(string),  // Removed in Cobbler 3
		TemplateFiles: d.Get("template_files").(string),
	}

	return distro
}
