package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductBaseResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("defectdojo_product.test", "tags.0", "bar"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "tags.1", "foo"),

					resource.TestCheckResourceAttr("defectdojo_product.test", "business_criticality", "high"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "enable_full_risk_acceptance", "true"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "enable_skip_risk_acceptance", "true"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "external_audience", "true"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "internet_accessible", "true"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "life_cycle", "production"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "origin", "internal"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "platform", "web"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "prod_numeric_grade", "100"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "regulation_ids.#", "0"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "revenue", "100.00"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "user_records", "1000000"),
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
				Config: testAccProductResourceMinimalConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// required attrs
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "test"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "product_type_id", "1"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "regulation_ids.#", "0"),

					resource.TestCheckNoResourceAttr("defectdojo_product.test", "business_criticality"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "life_cycle"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "origin"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "platform"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "prod_numeric_grade"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "revenue"),
					resource.TestCheckNoResourceAttr("defectdojo_product.test", "user_records"),

					// these booleans are reset to default value when not specified, they never return null
					resource.TestCheckResourceAttr("defectdojo_product.test", "enable_full_risk_acceptance", "false"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "enable_skip_risk_acceptance", "false"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "external_audience", "false"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "internet_accessible", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccProductResourceNoTags(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceNoTagsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "test"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "product_type_id", "1"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "tags.#", "0"),
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
				Config: testAccProductResourceNoTagsConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", updatedName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceEmptyTags(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceEmptyTagsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "test"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "product_type_id", "1"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "tags.#", "0"),
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
				Config: testAccProductResourceEmptyTagsConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", updatedName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceDeleteDrift(t *testing.T) {
	name := fmt.Sprintf("dox-delete-%s", resource.UniqueId())

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
			// Delete the underlying resource and see that it detects it has been deleted
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccProductResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDeleteResourceOutsideTerraform("defectdojo_product.test"),
				),
			},
			{
				Config: testAccProductResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "test"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "product_type_id", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceInvalid(t *testing.T) {
	name := fmt.Sprintf("dox-invalid-%s", resource.UniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(`.*Invalid\s+Attribute.*`),
				Config:      testAccProductResourceInvalidConfig(name),
			},
		},
	})
}

func testAccProductResourceNoTagsConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1

  business_criticality = "high"
  enable_full_risk_acceptance = false
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
}
`, name)
}

func testAccProductResourceEmptyTagsConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1

  business_criticality = "high"
  enable_full_risk_acceptance = false
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
	tags = []
}
`, name)
}

func testAccProductResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1
  tags = ["foo", "bar"]

  business_criticality = "high"
  enable_full_risk_acceptance = true
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
}
`, name)
}
func testAccProductResourceMinimalConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1
}
`, name)
}

func testAccProductResourceInvalidConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
  tags = ["foo", "BAR"]

  business_criticality = "something else"
  revenue = "a lot of money"
}
`, name)
}
