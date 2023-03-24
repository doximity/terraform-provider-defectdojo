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

func TestAccProductTypeBooleansDataSource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test default values of our booleans
			{
				Config: testAccProductTypeBooleanChecksDefaultConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "description", "test"),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "critical_product", "false"),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "key_product", "false"),
				),
			},
			// Test our booleans when defined as true
			{
				Config: testAccProductTypeBooleanChecksConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "description", "test"),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "critical_product", "true"),
					resource.TestCheckResourceAttr("data.defectdojo_product_type.test", "key_product", "true"),
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

func testAccProductTypeBooleanChecksDefaultConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
	name = %[1]q
	depends_on = [defectdojo_product_type.test]
}
`, name)
}

func testAccProductTypeBooleanChecksConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
	critical_product = true
	key_product = true
}
data "defectdojo_product_type" "test" {
	name = %[1]q
	depends_on = [defectdojo_product_type.test]
}
`, name)
}
