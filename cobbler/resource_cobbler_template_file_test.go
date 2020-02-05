package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	cobbler "github.com/wearespindle/cobblerclient"
)

func TestAccCobblerTemplateFile_basic(t *testing.T) {
	var ks cobbler.TemplateFile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCobblerPreCheck(t) },
		Providers:    testAccCobblerProviders,
		CheckDestroy: testAccCobblerCheckTemplateFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerTemplateFile_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckTemplateFileExists(t, "cobbler_template_file.foo", &ks),
				),
			},
		},
	})
}

func testAccCobblerCheckTemplateFileDestroy(s *terraform.State) error {
	config := testAccCobblerProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_template_file" {
			continue
		}

		if _, err := config.cobblerClient.GetTemplateFile(rs.Primary.ID); err == nil {
			return fmt.Errorf("Template File still exists")
		}
	}

	return nil
}

func testAccCobblerCheckTemplateFileExists(t *testing.T, n string, ks *cobbler.TemplateFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccCobblerProvider.Meta().(*Config)

		found, err := config.cobblerClient.GetTemplateFile(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Template File not found")
		}

		*ks = *found

		return nil
	}
}

var testAccCobblerTemplateFile_basic = `
	resource "cobbler_template_file" "foo" {
		name = "/var/lib/cobbler/templates/foo.ks"
		body = "I'm a Template file."
	}`
