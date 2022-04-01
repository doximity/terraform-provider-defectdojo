package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("dox-test-repo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", "dox-test-repo"),
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
			// TODO: truemilk/go-defectdojo does not support Products.Update(...) yet
			// {
			// 	Config: testAccExampleResourceConfig("dox-new-name"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("defectdojo_product.test", "name", "dox-new-name"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {
  base_url = "https://demo.defectdojo.org"
}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
}
`, name)
}
