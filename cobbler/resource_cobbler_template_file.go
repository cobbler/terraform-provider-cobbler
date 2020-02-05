package cobbler

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	cobbler "github.com/wearespindle/cobblerclient"
)

func resourceTemplateFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateFileCreate,
		Read:   resourceTemplateFileRead,
		Update: resourceTemplateFileUpdate,
		Delete: resourceTemplateFileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"body": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTemplateFileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ks := cobbler.TemplateFile{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler TemplateFile: Create Options: %#v", ks)

	if err := config.cobblerClient.CreateTemplateFile(ks); err != nil {
		return fmt.Errorf("Cobbler TemplateFile: Error Creating: %s", err)
	}

	d.SetId(ks.Name)

	return resourceTemplateFileRead(d, meta)
}

func resourceTemplateFileRead(d *schema.ResourceData, meta interface{}) error {
	// Since all attributes are required and not computed,
	// there's no reason to read.
	return nil
}

func resourceTemplateFileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ks := cobbler.TemplateFile{
		Name: d.Id(),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler TemplateFile: Updating Template (%s) with options: %+v", d.Id(), ks)

	if err := config.cobblerClient.CreateTemplateFile(ks); err != nil {
		return fmt.Errorf("Cobbler TemplateFile: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceTemplateFileRead(d, meta)
}

func resourceTemplateFileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := config.cobblerClient.DeleteTemplateFile(d.Id()); err != nil {
		return fmt.Errorf("Cobbler TemplateFile: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}
