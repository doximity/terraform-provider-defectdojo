package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductTypeResource(t *testing.T) {
	name := fmt.Sprintf("dox-test-pt-%s", resource.UniqueId())
	desc := fmt.Sprintf("dox test pt description %s", resource.UniqueId())
	updatedName := fmt.Sprintf("dox-new-pt-name-%s", resource.UniqueId())
	updatedDesc := fmt.Sprintf("updated description %s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductTypeResourceConfig(name, desc),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", desc),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", "true"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductTypeResourceConfig(updatedName, updatedDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", updatedDesc),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", "false"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProductTypeResourceConfig(name string, desc string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
  name = %[1]q
  description = %[2]q
  critical_product = true
  key_product = true
}
`, name, desc)
}
