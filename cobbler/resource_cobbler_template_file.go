package cobbler

import (
	"context"
	cobbler "github.com/cobbler/cobblerclient"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTemplateFile() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_template_file` manages a template file within Cobbler.",
		CreateContext: resourceTemplateFileCreate,
		ReadContext:   resourceTemplateFileRead,
		UpdateContext: resourceTemplateFileUpdate,
		DeleteContext: resourceTemplateFileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the template file. This must be the name only, so without `/var/lib/cobbler/templates`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"body": {
				Description: "The body of the template file. May also point to a file: `body = file(\"my_template.ks\")`.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceTemplateFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	ks := cobbler.TemplateFile{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	tflog.Debug(ctx, "Cobbler TemplateFile: Create Options", map[string]interface{}{
		"options": structs.Map(ks),
	})

	if err := config.cobblerClient.CreateTemplateFile(ks); err != nil {
		return diag.Errorf("Cobbler TemplateFile: Error Creating: %s", err)
	}

	d.SetId(ks.Name)

	return resourceTemplateFileRead(ctx, d, meta)
}

func resourceTemplateFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since all attributes are required and not computed, there's no reason to read.
	return nil
}

func resourceTemplateFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	ks := cobbler.TemplateFile{
		Name: d.Id(),
		Body: d.Get("body").(string),
	}

	tflog.Debug(ctx, "Cobbler TemplateFile: Updating Template with options", map[string]interface{}{
		"template": d.Id(),
		"options":  structs.Map(ks),
	})

	if err := config.cobblerClient.CreateTemplateFile(ks); err != nil {
		return diag.Errorf("Cobbler TemplateFile: Error Updating (%s): %s", d.Id(), err)
	}

	return resourceTemplateFileRead(ctx, d, meta)
}

func resourceTemplateFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	if err := config.cobblerClient.DeleteTemplateFile(d.Id()); err != nil {
		return diag.Errorf("Cobbler TemplateFile: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}
