package cobbler

import (
	"fmt"
	"log"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSnippet() *schema.Resource {
	return &schema.Resource{
		Description: "`cobbler_snippet` manages a snippet within Cobbler.",
		Create:      resourceSnippetCreate,
		Read:        resourceSnippetRead,
		Update:      resourceSnippetUpdate,
		Delete:      resourceSnippetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the snippet. This must be the name only, so without `/var/lib/cobbler/snippets`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"body": {
				Description: "The body of the snippet. May also point to a file: `body = file(\"my_snippet\")`.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceSnippetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	snippet := cobbler.Snippet{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler Snippet: Create Options: %#v", snippet)

	if err := config.cobblerClient.CreateSnippet(snippet); err != nil {
		return fmt.Errorf("Cobbler Snippet: Error Creating: %s", err)
	}

	d.SetId(snippet.Name)

	return resourceSnippetRead(d, meta)
}

func resourceSnippetRead(d *schema.ResourceData, meta interface{}) error {
	// Since all attributes are required and not computed,
	// there's no reason to read.
	return nil
}

func resourceSnippetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	snippet := cobbler.Snippet{
		Name: d.Id(),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler Snippet: Updating Snippet (%s) with options: %+v", d.Id(), snippet)

	if err := config.cobblerClient.CreateSnippet(snippet); err != nil {
		return fmt.Errorf("Cobbler Snippet: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceSnippetRead(d, meta)
}

func resourceSnippetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := config.cobblerClient.DeleteSnippet(d.Id()); err != nil {
		return fmt.Errorf("Cobbler Snippet: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}
