package cobbler

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSnippet() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_snippet` manages a snippet within Cobbler.",
		CreateContext: resourceSnippetCreate,
		ReadContext:   resourceSnippetRead,
		UpdateContext: resourceSnippetUpdate,
		DeleteContext: resourceSnippetDelete,

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

func resourceSnippetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	snippet := cobbler.Snippet{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler Snippet: Create Options: %#v", snippet)

	if err := config.cobblerClient.CreateSnippet(snippet); err != nil {
		//goland:noinspection GoErrorStringFormat
		return diag.Errorf("Cobbler Snippet: Error Creating: %s", err)
	}

	d.SetId(snippet.Name)

	return resourceSnippetRead(ctx, d, meta)
}

func resourceSnippetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since all attributes are required and not computed,
	// there's no reason to read.
	return nil
}

func resourceSnippetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	snippet := cobbler.Snippet{
		Name: d.Id(),
		Body: d.Get("body").(string),
	}

	log.Printf("[DEBUG] Cobbler Snippet: Updating Snippet (%s) with options: %+v", d.Id(), snippet)

	if err := config.cobblerClient.CreateSnippet(snippet); err != nil {
		//goland:noinspection GoErrorStringFormat
		return diag.Errorf("Cobbler Snippet: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceSnippetRead(ctx, d, meta)
}

func resourceSnippetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	if err := config.cobblerClient.DeleteSnippet(d.Id()); err != nil {
		//goland:noinspection GoErrorStringFormat
		return diag.Errorf("Cobbler Snippet: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}
