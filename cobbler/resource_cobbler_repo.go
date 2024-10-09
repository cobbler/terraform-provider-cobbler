package cobbler

import (
	"context"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_repo` manages a repo within Cobbler.",
		CreateContext: resourceRepoCreate,
		ReadContext:   resourceRepoRead,
		UpdateContext: resourceRepoUpdate,
		DeleteContext: resourceRepoDelete,

		Schema: map[string]*schema.Schema{
			"apt_components": {
				Description: "List of Apt components such as main, restricted, universe. Applicable to apt breeds only.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"apt_dists": {
				Description: "List of Apt distribution names such as focal, focal-updates. Applicable to apt breeds only.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"arch": {
				Description: "The architecture of the repo. Valid options are: i386, x86_64, ia64, ppc, ppc64, s390, arm.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"breed": {
				Description: "The \"breed\" of distribution. Valid options are: rsync, rhn, yum, apt, and wget. These choices may vary depending on the version of Cobbler in use.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"createrepo_flags": {
				Description: "Flags to use with `createrepo`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"environment": {
				Description: "Environment variables to use during repo command execution.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"keep_updated": {
				Description: "Update the repo upon Cobbler sync. Valid values are true or false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"mirror": {
				Description: "Address of the repo to mirror.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"mirror_locally": {
				Description: "Whether to copy the files locally or just references to the external files. Valid values are true or false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "A name for the repo.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"owners": {
				Description: "List of Owners for authz_ownership.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"proxy": {
				Description: "Proxy to use for downloading the repo. This argument does not work on older versions of Cobbler.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"rpm_list": {
				Description: "List of specific RPMs to mirror.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func resourceRepoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Create a cobblerclient.Repo
	repo := buildRepo(d, config)

	// Attempt to create the Repo
	tflog.Debug(ctx, "Cobbler Repo: Create Options", map[string]interface{}{
		"options": structs.Map(repo),
	})
	newRepo, err := config.cobblerClient.CreateRepo(repo)
	if err != nil {
		return diag.Errorf("Cobbler Repo: Error Creating: %s", err)
	}

	d.SetId(newRepo.Name)

	return resourceRepoRead(ctx, d, meta)
}

func resourceRepoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Retrieve the repo from cobbler
	repo, err := config.cobblerClient.GetRepo(d.Id())
	if err != nil {
		return diag.Errorf("Cobbler Repo: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	err = d.Set("arch", repo.Arch)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("breed", repo.Breed)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("comment", repo.Comment)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("createrepo_flags", repo.CreateRepoFlags)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("keep_updated", repo.KeepUpdated)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mirror", repo.Mirror)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mirror_locally", repo.MirrorLocally)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name", repo.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("proxy", repo.Proxy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("apt_components", repo.AptComponents)
	if err != nil {
		tflog.Debug(ctx, "Unable to set apt_components", map[string]interface{}{
			"error": err,
		})
	}

	err = d.Set("apt_dists", repo.AptDists)
	if err != nil {
		tflog.Debug(ctx, "Unable to set apt_dists", map[string]interface{}{
			"error": err,
		})
	}

	err = d.Set("owners", repo.Owners)
	if err != nil {
		tflog.Debug(ctx, "Unable to set owners", map[string]interface{}{
			"error": err,
		})
	}

	err = d.Set("rpm_list", repo.RpmList)
	if err != nil {
		tflog.Debug(ctx, "Unable to set rpm_list", map[string]interface{}{
			"error": err,
		})
	}

	return nil
}

func resourceRepoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// create a cobblerclient.Repo
	repo := buildRepo(d, config)

	// Attempt to updateh the repo with new information
	tflog.Debug(ctx, "Cobbler Repo: Updating Repo with options", map[string]interface{}{
		"repo":    d.Id(),
		"options": structs.Map(repo),
	})
	err := config.cobblerClient.UpdateRepo(&repo)
	if err != nil {
		return diag.Errorf("Cobbler Repo: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceRepoRead(ctx, d, meta)
}

func resourceRepoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Attempt to delete the repo
	if err := config.cobblerClient.DeleteRepo(d.Id()); err != nil {
		return diag.Errorf("Cobbler Repo: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildRepo builds a cobbler.Repo from the Terraform attributes.
func buildRepo(d *schema.ResourceData, meta interface{}) cobbler.Repo { //nolint:unparam // We satisfy our own pattern here
	aptComponents := []string{}
	for _, i := range d.Get("apt_components").([]interface{}) {
		aptComponents = append(aptComponents, i.(string))
	}

	aptDists := []string{}
	for _, i := range d.Get("apt_dists").([]interface{}) {
		aptDists = append(aptDists, i.(string))
	}

	owners := []string{}
	for _, i := range d.Get("owners").([]interface{}) {
		owners = append(owners, i.(string))
	}

	rpmList := []string{}
	for _, i := range d.Get("rpm_list").([]interface{}) {
		rpmList = append(rpmList, i.(string))
	}

	repo := cobbler.Repo{
		AptComponents:   aptComponents,
		AptDists:        aptDists,
		Arch:            d.Get("arch").(string),
		Breed:           d.Get("breed").(string),
		Comment:         d.Get("comment").(string),
		CreateRepoFlags: d.Get("createrepo_flags").(string),
		KeepUpdated:     d.Get("keep_updated").(bool),
		Mirror:          d.Get("mirror").(string),
		MirrorLocally:   d.Get("mirror_locally").(bool),
		Name:            d.Get("name").(string),
		Owners:          owners,
		Proxy:           d.Get("proxy").(string),
		RpmList:         rpmList,
	}

	return repo
}
