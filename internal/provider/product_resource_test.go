package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func testAccDeleteResourceOutsideTerraform(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// retrieve the client from the test provider
		client, err := newClient(context.Background(), os.Getenv("DEFECTDOJO_BASEURL"), os.Getenv("DEFECTDOJO_APIKEY"), os.Getenv("DEFECTDOJO_USERNAME"), os.Getenv("DEFECTDOJO_PASSWORD"))
		if err != nil {
			return err
		}

		i, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.ProductsDestroy(context.Background(), i)
		if err != nil {
			return err
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("bad status code deleting the resource: %d", resp.StatusCode)
		}

		tflog.Debug(context.Background(), fmt.Sprintf("Deleted resource, ID: %d", i))

		return nil
	}
}
