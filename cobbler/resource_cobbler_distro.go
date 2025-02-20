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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arch": {
				Description: "The architecture of the distro. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"breed": {
				Description: "The \"breed\" of distribution. Valid options are: redhat, fedora, centos, scientific linux, suse, debian, and ubuntu. These choices may vary depending on the version of Cobbler in use.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"boot_files": {
				Description:   "Files copied into tftpboot beyond the kernel/initrd.",
				Type:          schema.TypeMap,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"boot_files_inherit"},
			},
			"boot_files_inherit": {
				Description:   "Signal that boot_files should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"boot_files"},
			},
			"boot_loaders": {
				Description: "Must be either 'grub', 'pxe', or 'ipxe'.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"boot_loaders_inherit": {
				Description:   "Signal that boot_loaders should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"boot_loaders"},
			},
			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
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
			"fetchable_files_inherit": {
				Description:   "Signal that fetchable_files should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"fetchable_files"},
			},
			"initrd": {
				Description: "Absolute path to initrd on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"remote_boot_initrd": {
				Description: "URL the bootloader directly retrieves and boots from",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"kernel": {
				Description: "Absolute path to kernel on filesystem. This must already exist prior to creating the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"remote_boot_kernel": {
				Description: "URL the bootloader directly retrieves and boots from",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"kernel_options": {
				Description: "Kernel options to use with the kernel.",
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
				Description: "Post install Kernel options to use with the kernel after installation.",
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
				Description: "Management classes for external config management.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"mgmt_classes_inherit": {
				Description:   "Signal that mgmt_classes should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"mgmt_classes"},
			},
			"name": {
				Description: "A name for the distro.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"os_version": {
				Description: "The version of the distro you are creating. This varies with the version of Cobbler you are using. An updated signature list may need to be obtained in order to support a newer version. Example: `focal`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"owners": {
				Description: "Owners list for authz_ownership.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"owners_inherit": {
				Description:   "Signal that owners should be set to inherit from its parent",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"owners"},
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
		},
	}
}

func resourceDistroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Create a cobblerclient.Distro
	distro, err := buildDistro(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

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
	distro, err := config.cobblerClient.GetDistro(d.Id(), false, false)
	if err != nil {
		return diag.Errorf("Cobbler Distro: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	err = d.Set("name", distro.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("arch", distro.Arch)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("breed", distro.Breed)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "boot_files", distro.BootFiles, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "boot_loaders", distro.BootLoaders, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("comment", distro.Comment)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "fetchable_files", distro.FetchableFiles, make(map[string]interface{}))
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
	err = d.Set("remote_boot_initrd", distro.RemoteBootInitrd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("remote_boot_kernel", distro.RemoteBootKernel)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "kernel_options", distro.KernelOptions, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "kernel_options_post", distro.KernelOptionsPost, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("initrd", distro.Initrd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "mgmt_classes", distro.MgmtClasses, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("os_version", distro.OSVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "owners", distro.Owners, make([]string, 0))
	if err != nil {
		return diag.FromErr(err)
	}
	err = SetInherit(d, "template_files", distro.TemplateFiles, make(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDistroUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// create a cobblerclient.Distro
	distro, err := buildDistro(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to updateh the distro with new information
	tflog.Debug(ctx, "Cobbler Distro: Updating Distro", map[string]interface{}{
		"distro":  d.Id(),
		"options": structs.Map(distro),
	})
	err = config.cobblerClient.UpdateDistro(&distro)
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
func buildDistro(d *schema.ResourceData, meta interface{}) (cobbler.Distro, error) { //nolint:unparam // We satisfy our own pattern here
	mgmtClasses, err := GetStringSlice(d, "mgmt_classes")
	if err != nil {
		return cobbler.Distro{}, err
	}
	owners, err := GetStringSlice(d, "owners")
	if err != nil {
		return cobbler.Distro{}, err
	}
	bootFiles, err := GetInterfaceMap(d, "boot_files")
	if err != nil {
		return cobbler.Distro{}, err
	}
	bootLoaders, err := GetStringSlice(d, "boot_loaders")
	if err != nil {
		return cobbler.Distro{}, err
	}
	fetchableFiles, err := GetInterfaceMap(d, "fetchable_files")
	if err != nil {
		return cobbler.Distro{}, err
	}
	kernelOptions, err := GetInterfaceMap(d, "kernel_options")
	if err != nil {
		return cobbler.Distro{}, err
	}
	kernelOptionsPost, err := GetInterfaceMap(d, "kernel_options_post")
	if err != nil {
		return cobbler.Distro{}, err
	}
	templateFiles, err := GetInterfaceMap(d, "template_files")
	if err != nil {
		return cobbler.Distro{}, err
	}

	distro := cobbler.NewDistro()
	distro.Arch = d.Get("arch").(string)
	distro.Breed = d.Get("breed").(string)
	distro.BootFiles = cobbler.Value[map[string]interface{}]{
		Data:        bootFiles,
		IsInherited: IsOptionInherited(d, "boot_files"),
	}
	distro.BootLoaders = cobbler.Value[[]string]{
		Data:        bootLoaders,
		IsInherited: IsOptionInherited(d, "boot_loaders"),
	}
	distro.Comment = d.Get("comment").(string)
	distro.FetchableFiles = cobbler.Value[map[string]interface{}]{
		Data:        fetchableFiles,
		IsInherited: IsOptionInherited(d, "fetchable_files"),
	}
	distro.Kernel = d.Get("kernel").(string)
	distro.KernelOptions = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptions,
		IsInherited: IsOptionInherited(d, "kernel_options"),
	}
	distro.KernelOptionsPost = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptionsPost,
		IsInherited: IsOptionInherited(d, "kernel_options_post"),
	}
	distro.Initrd = d.Get("initrd").(string)
	distro.MgmtClasses = cobbler.Value[[]string]{
		Data:        mgmtClasses,
		IsInherited: IsOptionInherited(d, "mgmt_classes"),
	}
	distro.Name = d.Get("name").(string)
	distro.OSVersion = d.Get("os_version").(string)
	distro.Owners = cobbler.Value[[]string]{
		Data:        owners,
		IsInherited: IsOptionInherited(d, "owners"),
	}
	distro.TemplateFiles = cobbler.Value[map[string]interface{}]{
		Data:        templateFiles,
		IsInherited: IsOptionInherited(d, "template_files"),
	}

	return distro, nil
}
