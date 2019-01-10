package cobbler

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	cobbler "github.com/jtopjian/cobblerclient"
)

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepoCreate,
		Read:   resourceRepoRead,
		Update: resourceRepoUpdate,
		Delete: resourceRepoDelete,

		Schema: map[string]*schema.Schema{
			"apt_components": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"apt_dists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"arch": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"breed": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"createrepo_flags": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"keep_updated": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"mirror": {
				Type:     schema.TypeString,
				Required: true,
			},

			"mirror_locally": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"owners": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"proxy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rpm_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			//"yumopts": &schema.Schema{
			//	Type:     schema.TypeMap,
			//	Optional: true,
			//	Computed: true,
			//},
		},
	}
}

func resourceRepoCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create a cobblerclient.Repo
	repo := buildRepo(d, config)

	// Attempte to create the Repo
	log.Printf("[DEBUG] Cobbler Repo: Create Options: %#v", repo)
	newRepo, err := config.cobblerClient.CreateRepo(repo)
	if err != nil {
		return fmt.Errorf("Cobbler Repo: Error Creating: %s", err)
	}

	d.SetId(newRepo.Name)

	return resourceRepoRead(d, meta)
}

func resourceRepoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Retrieve the repo from cobbler
	repo, err := config.cobblerClient.GetRepo(d.Id())
	if err != nil {
		return fmt.Errorf("Cobbler Repo: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	d.Set("arch", repo.Arch)
	d.Set("breed", repo.Breed)
	d.Set("comment", repo.Comment)
	d.Set("createrepo_flags", repo.CreateRepoFlags)
	d.Set("environment", repo.Environment)
	d.Set("keep_updated", repo.KeepUpdated)
	d.Set("mirror", repo.Mirror)
	d.Set("mirror_locally", repo.MirrorLocally)
	d.Set("name", repo.Name)
	d.Set("proxy", repo.Proxy)
	//d.Set("yumopts", repo.YumOpts)

	err = d.Set("apt_components", repo.AptComponents)
	if err != nil {
		log.Printf("[DEBUG] Unable to set apt_components: %s", err)
	}

	err = d.Set("apt_dists", repo.AptDists)
	if err != nil {
		log.Printf("[DEBUG] Unable to set apt_dists: %s", err)
	}

	err = d.Set("owners", repo.Owners)
	if err != nil {
		log.Printf("[DEBUG] Unable to set owners: %s", err)
	}

	err = d.Set("rpm_list", repo.RpmList)
	if err != nil {
		log.Printf("[DEBUG] Unable to set rpm_list: %s", err)
	}

	return nil
}

func resourceRepoUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// create a cobblerclient.Repo
	repo := buildRepo(d, config)

	// Attempt to updateh the repo with new information
	log.Printf("[DEBUG] Cobbler Repo: Updating Repo (%s) with options: %+v", d.Id(), repo)
	err := config.cobblerClient.UpdateRepo(&repo)
	if err != nil {
		return fmt.Errorf("Cobbler Repo: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceRepoRead(d, meta)
}

func resourceRepoDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Attempt to delete the repo
	if err := config.cobblerClient.DeleteRepo(d.Id()); err != nil {
		return fmt.Errorf("Cobbler Repo: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildRepo builds a cobbler.Repo from the Terraform attributes
func buildRepo(d *schema.ResourceData, meta interface{}) cobbler.Repo {
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

	//yumOpts := make(map[string]interface{})
	//y := d.Get("yum_opts")
	//if y != nil {
	//	m := y.(map[string]interface{})
	//	for k, v := range m {
	//		yumOpts[k] = v
	//	}
	//}

	repo := cobbler.Repo{
		AptComponents:   aptComponents,
		AptDists:        aptDists,
		Arch:            d.Get("arch").(string),
		Breed:           d.Get("breed").(string),
		Comment:         d.Get("comment").(string),
		CreateRepoFlags: d.Get("createrepo_flags").(string),
		Environment:     d.Get("environment").(string),
		KeepUpdated:     d.Get("keep_updated").(bool),
		Mirror:          d.Get("mirror").(string),
		MirrorLocally:   d.Get("mirror_locally").(bool),
		Name:            d.Get("name").(string),
		Owners:          owners,
		Proxy:           d.Get("proxy").(string),
		RpmList:         rpmList,
		//YumOpts:         yumOpts,
	}

	return repo
}
