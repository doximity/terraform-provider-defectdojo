package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "name", "dox-data-test"),
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "description", "test"),
				),
			},
		},
	})
}

const testAccProductDataSourceConfig = `
provider "defectdojo" {
  base_url = "https://demo.defectdojo.org"
}
resource "defectdojo_product" "test" {
  name = "dox-data-test"
	description = "test"
  product_type_id = 1
}
data "defectdojo_product" "test" {
  id = defectdojo_product.test.id
  depends_on = [defectdojo_product.test]
}
`