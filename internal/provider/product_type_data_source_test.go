package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductTypeIdDataSource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductTypeDataSourceIdConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "description", "test"),
				),
			},
		},
	})
}

func TestAccProductTypeNameDataSource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductTypeDataSourceNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "description", "test"),
				),
			},
		},
	})
}

func testAccProductTypeDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
  name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
  id = defectdojo_product_type.test.id
  depends_on = [defectdojo_product_type.test]
}
`, name)
}

func testAccProductTypeDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
  name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
  name = %[2]q
  depends_on = [defectdojo_product_type.test]
}
`, name, name)
}
