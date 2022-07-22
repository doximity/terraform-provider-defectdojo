package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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
		var resp *http.Response
		if match, err := regexp.MatchString(`^defectdojo_product\.`, resourceName); err == nil && match {
			resp, err = client.ProductsDestroy(context.Background(), i)
			if err != nil {
				return err
			}
		} else if match, err := regexp.MatchString(`^defectdojo_jira_product_configuration\.`, resourceName); err == nil && match {
			resp, err = client.JiraProductConfigurationsDestroy(context.Background(), i)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unknown type: %s, %s", resourceName, err)
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("bad status code deleting the resource: %d", resp.StatusCode)
		}

		tflog.Debug(context.Background(), fmt.Sprintf("Deleted resource, ID: %d", i))

		return nil
	}
}
