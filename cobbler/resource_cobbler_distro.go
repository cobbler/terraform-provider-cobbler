package cobbler

import (
	"context"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDistro() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_distro` manages a distribution within Cobbler.",
		CreateContext: resourceDistroCreate,
		ReadContext:   resourceDistroRead,
		UpdateContext: resourceDistroUpdate,
		DeleteContext: resourceDistroDelete,

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
			"boot_loaders": {
				Description: "Must be either 'grub', 'pxe', or 'ipxe'.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceDistroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Create a cobblerclient.Distro
	distro := buildDistro(d, config)

	// Attempt to create the Distro
	tflog.Debug(ctx, "Cobbler Distro Create Options", map[string]interface{}{
		"options": structs.Map(distro),
	})
	newDistro, err := config.cobblerClient.CreateDistro(distro)
	if err != nil {
		return diag.Errorf("Cobbler Distro: Error Creating: %s", err)
	}

	d.SetId(newDistro.Name)

	return resourceDistroRead(ctx, d, meta)
}

func resourceDistroRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Retrieve the distro from cobbler
	distro, err := config.cobblerClient.GetDistro(d.Id())
	if err != nil {
		return diag.Errorf("Cobbler Distro: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	err = d.Set("arch", distro.Arch)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("breed", distro.Breed)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("boot_files", distro.BootFiles)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("boot_loaders", distro.BootLoaders)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("comment", distro.Comment)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("fetchable_files", distro.FetchableFiles)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("initrd", distro.Initrd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel", distro.Kernel)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options", distro.KernelOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options_post", distro.KernelOptionsPost)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("initrd", distro.Initrd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mgmt_classes", distro.MGMTClasses)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("os_version", distro.OSVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("owners", distro.Owners)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("template_files", distro.TemplateFiles)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDistroUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// create a cobblerclient.Distro
	distro := buildDistro(d, config)

	// Attempt to updateh the distro with new information
	tflog.Debug(ctx, "Cobbler Distro: Updating Distro", map[string]interface{}{
		"distro":  d.Id(),
		"options": structs.Map(distro),
	})
	err := config.cobblerClient.UpdateDistro(&distro)
	if err != nil {
		return diag.Errorf("Cobbler Distro: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceDistroRead(ctx, d, meta)
}

func resourceDistroDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Attempt to delete the distro
	if err := config.cobblerClient.DeleteDistro(d.Id()); err != nil {
		return diag.Errorf("Cobbler Distro: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildDistro builds a cobbler.Distro from the Terraform attributes.
func buildDistro(d *schema.ResourceData, meta interface{}) cobbler.Distro { //nolint:unparam // We satisfy our own pattern here
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
	bootLoaders := []string{}
	for _, i := range d.Get("boot_loaders").([]interface{}) {
		bootLoaders = append(bootLoaders, i.(string))
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
		Arch:              d.Get("arch").(string),
		Breed:             d.Get("breed").(string),
		BootFiles:         bootFiles,
		BootLoaders:       bootLoaders,
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
