package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraProductConfigurationResource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	jirakey := fmt.Sprintf("APPSEC%s", resource.UniqueId())
	newjirakey := fmt.Sprintf("APPSEC%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "project_key", jirakey),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "issue_template_dir", ""),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "push_all_issues", "false"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "enable_engagement_epic_mapping", "false"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "push_notes", "false"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "product_jira_sla_notification", "false"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "risk_acceptance_expiration_notification", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_jira_product_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccJiraProductConfigurationResourceUpdateConfig(name, newjirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "project_key", newjirakey),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "issue_template_dir", "some/dir"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "push_all_issues", "true"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "enable_engagement_epic_mapping", "true"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "push_notes", "true"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "product_jira_sla_notification", "true"),
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "risk_acceptance_expiration_notification", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccJiraProductConfigurationResourceDeleteDrift(t *testing.T) {
	name := fmt.Sprintf("dox-delete-%s", resource.UniqueId())
	jirakey := fmt.Sprintf("APPSEC%s", resource.UniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "project_key", jirakey),
				),
			},
			// Delete the underlying resource and see that it detects it has been deleted
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccJiraProductConfigurationResourceConfig(name, jirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDeleteResourceOutsideTerraform("defectdojo_jira_product_configuration.test"),
				),
			},
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_product_configuration.test", "project_key", jirakey),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccJiraProductConfigurationResourceInvalid(t *testing.T) {
	name := fmt.Sprintf("dox-delete-%s", resource.UniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ExpectError: regexp.MustCompile(`Invalid\s+Resource`),
				Config:      testAccInvalidJiraProductConfigurationResourceConfig(name),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccInvalidJiraProductConfigurationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_jira_product_configuration" "test" {
  project_key = %[1]q
}
`, name)
}

func testAccJiraProductConfigurationResourceConfig(productname string, name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
	description = "test"
  product_type_id = 1
}
resource "defectdojo_jira_product_configuration" "test" {
  product_id = defectdojo_product.test.id
  project_key = %[2]q
}
`, productname, name)
}

func testAccJiraProductConfigurationResourceUpdateConfig(productname string, name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
	description = "test"
  product_type_id = 1
}
resource "defectdojo_jira_product_configuration" "test" {
  product_id = defectdojo_product.test.id
  project_key = %[2]q
  issue_template_dir = "some/dir"
  push_all_issues = true
  enable_engagement_epic_mapping = true
  push_notes = true
  product_jira_sla_notification = true
  risk_acceptance_expiration_notification = true
}
`, productname, name)
}
