package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductTypeResource(t *testing.T) {
	name := fmt.Sprintf("dox-test-pt-%s", resource.UniqueId())
	desc := fmt.Sprintf("dox test pt description %s", resource.UniqueId())
	criticalProduct := "true"
	keyProduct := "true"
	updatedName := fmt.Sprintf("dox-new-pt-name-%s", resource.UniqueId())
	updatedDesc := fmt.Sprintf("updated description %s", resource.UniqueId())
	updatedCriticalProduct := "false"
	updatedKeyProduct := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductTypeResourceConfig(name, desc, criticalProduct, keyProduct),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", desc),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", criticalProduct),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", keyProduct),
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
				Config: testAccProductTypeResourceConfig(updatedName, updatedDesc, updatedCriticalProduct, updatedKeyProduct),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", updatedDesc),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", updatedCriticalProduct),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", updatedKeyProduct),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProductTypeResourceConfig(name string, desc string, criticalProduct string, keyProduct string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
  name = %[1]q
  description = %[2]q
  critical_product = %[3]q
  key_product = %[4]q
}
`, name, desc, criticalProduct, keyProduct)
}
