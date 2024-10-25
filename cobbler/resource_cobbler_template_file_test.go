package cobbler

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	cobbler "github.com/cobbler/cobblerclient"
)

func TestAccCobblerTemplateFile_basic(t *testing.T) {
	var ks cobbler.TemplateFile

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCobblerPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCobblerCheckTemplateFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCobblerTemplateFileBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCobblerCheckTemplateFileExists("cobbler_template_file.foo", &ks),
				),
			},
			{
				ResourceName:      "cobbler_template_file.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCobblerCheckTemplateFileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cobbler_template_file" {
			continue
		}

		if _, err := cobblerApiClient.GetTemplateFile(rs.Primary.ID); err == nil {
			return fmt.Errorf("template file still exists")
		}
	}

	return nil
}

func testAccCobblerCheckTemplateFileExists(n string, ks *cobbler.TemplateFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		found, err := cobblerApiClient.GetTemplateFile(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("template file not found")
		}

		*ks = *found

		return nil
	}
}

var testAccCobblerTemplateFileBasic = `
	resource "cobbler_template_file" "foo" {
		name = "foo.seed"
		body = "I'm a Template file."
	}`
