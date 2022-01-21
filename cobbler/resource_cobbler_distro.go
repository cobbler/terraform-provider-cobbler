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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			//			"boot_loader": {
			//				Description: "Must be either 'grub', 'pxe', or 'ipxe'.",
			//				Type:        schema.TypeString,
			//				Optional:    true,
			//				ForceNew:    true,
			//				Computed:    true,
			//			},

			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
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
			"initrd": {
				Description: "Absolute path to initrd on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kernel": {
				Description: "Absolute path to kernel on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kernel_options": {
				Description: "Kernel options to use with the kernel.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_post": {
				Description: "Post install Kernel options to use with the kernel after installation.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
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
				Description: "The version of the distro you are creating. This varies with the version of Cobbler you are using. An updated signature list may need to be obtained in order to support a newer version. Example: `focal`.",
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
			"template_files": {
				Description: "File mappings for built-in config management.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	//d.Set("boot_loader", distro.BootLoader)
	d.Set("comment", distro.Comment)
	d.Set("fetchable_files", distro.FetchableFiles)
	d.Set("initrd", distro.Initrd)
	d.Set("kernel", distro.Kernel)
	d.Set("kernel_options", distro.KernelOptions)
	d.Set("kernel_options_post", distro.KernelOptionsPost)
	d.Set("initrd", distro.Initrd)
	d.Set("mgmt_classes", distro.MGMTClasses)
	d.Set("os_version", distro.OSVersion)
	d.Set("owners", distro.Owners)
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

	distro := cobbler.Distro{
		Arch:      d.Get("arch").(string),
		Breed:     d.Get("breed").(string),
		BootFiles: bootFiles,
		//BootLoader:          d.Get("boot_loader").(string),
		Comment:           d.Get("comment").(string),
		FetchableFiles:    fetchableFiles,
		Kernel:            d.Get("kernel").(string),
		KernelOptions:     kernelOptions,
		KernelOptionsPost: kernelOptionsPost,
		Initrd:            d.Get("initrd").(string),
		MGMTClasses:       mgmtClasses,
		Name:              d.Get("name").(string),
		OSVersion:         d.Get("os_version").(string),
		Owners:            owners,
		TemplateFiles:     templateFiles,
	}

	return distro
}
