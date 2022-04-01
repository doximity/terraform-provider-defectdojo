package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductResource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "test"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "product_type_id", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", updatedName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProductResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
}
`, name)
}
