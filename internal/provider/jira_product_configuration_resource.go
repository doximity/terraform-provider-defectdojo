package provider

import (
	"context"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//type jiraProductConfigurationResourceType struct{}

func (t *jiraProductConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Jira Product Configuration is the connection between a Product and a Jira Instance. It defines the Product's settings for pushing Findings to Jira.",

		Attributes: map[string]schema.Attribute{
			"project_key": schema.StringAttribute{
				MarkdownDescription: "The Jira Project Key",
				Optional:            true,
				Computed:            true,
			},

			"issue_template_dir": schema.StringAttribute{
				MarkdownDescription: "The folder containing Django templates used to render the JIRA issue description. Leave empty to use the default jira_full templates.",
				Optional:            true,
				Computed:            true,
			},

			"push_all_issues": schema.BoolAttribute{
				MarkdownDescription: "Automatically maintain parity with JIRA. Always create and update JIRA tickets for findings in this Product.",
				Optional:            true,
				Computed:            true,
			},

			"enable_engagement_epic_mapping": schema.BoolAttribute{
				MarkdownDescription: "Whether to map engagements to epics in Jira",
				Optional:            true,
				Computed:            true,
			},

			"push_notes": schema.BoolAttribute{
				MarkdownDescription: "Whether to push notes to Jira",
				Optional:            true,
				Computed:            true,
			},

			"product_jira_sla_notification": schema.BoolAttribute{
				MarkdownDescription: "Send SLA notifications as comments",
				Optional:            true,
				Computed:            true,
			},

			"risk_acceptance_expiration_notification": schema.BoolAttribute{
				MarkdownDescription: "Send Risk Acceptance expiration notifications as comments",
				Optional:            true,
				Computed:            true,
			},

			"jira_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Jira Instance to use for this Product",
				Optional:            true,
			},

			"product_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Product to configure. Although optional, either the Product ID or the Engagement ID must be defined to create a Jira Product Configuration.",
				Optional:            true,
			},

			"engagement_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Engagement. Although optional, either the Product ID or the Engagement ID must be defined to create a Jira Product Configuration.",
				Optional:            true,
			},

			"id": schema.StringAttribute{ // the id (for import purposes) MUST be a string
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type jiraProductConfigurationResourceData struct {
	ProjectKey                           types.String `tfsdk:"project_key" ddField:"ProjectKey"`
	IssueTemplateDir                     types.String `tfsdk:"issue_template_dir" ddField:"IssueTemplateDir"`
	PushAllIssues                        types.Bool   `tfsdk:"push_all_issues" ddField:"PushAllIssues"`
	EnableEngagementEpicMapping          types.Bool   `tfsdk:"enable_engagement_epic_mapping" ddField:"EnableEngagementEpicMapping"`
	PushNotes                            types.Bool   `tfsdk:"push_notes" ddField:"PushNotes"`
	ProductJiraSlaNotification           types.Bool   `tfsdk:"product_jira_sla_notification" ddField:"ProductJiraSlaNotification"`
	RiskAcceptanceExpirationNotification types.Bool   `tfsdk:"risk_acceptance_expiration_notification" ddField:"RiskAcceptanceExpirationNotification"`
	JiraInstance                         types.String `tfsdk:"jira_instance_id" ddField:"JiraInstance"`
	Product                              types.String `tfsdk:"product_id" ddField:"Product"`
	Engagement                           types.String `tfsdk:"engagement_id" ddField:"Engagement"`
	Id                                   types.String `tfsdk:"id" ddField:"Id"`
}

type jiraProductConfigurationDefectdojoResource struct {
	dd.JIRAProject
}

func (ddr *jiraProductConfigurationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := dd.JiraProductConfigurationsCreateJSONRequestBody(ddr.JIRAProject)
	apiResp, err := client.JiraProductConfigurationsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.JIRAProject = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.JiraProductConfigurationsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.JIRAProject = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := dd.JiraProductConfigurationsUpdateJSONRequestBody(ddr.JIRAProject)
	apiResp, err := client.JiraProductConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.JIRAProject = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.JiraProductConfigurationsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *jiraProductConfigurationResourceData) id() types.String {
	return d.Id
}

func (d *jiraProductConfigurationResourceData) defectdojoResource() defectdojoResource {
	return &jiraProductConfigurationDefectdojoResource{
		JIRAProject: dd.JIRAProject{},
	}
}

type jiraProductConfigurationResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &jiraProductConfigurationResource{}
var _ resource.ResourceWithImportState = &jiraProductConfigurationResource{}

func NewJiraProductConfigurationResource() resource.Resource {
	return &jiraProductConfigurationResource{
		terraformResource: terraformResource{
			dataProvider: jiraProductConfigurationDataProvider{},
		},
	}
}

func (r jiraProductConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_product_configuration"
}

func (r jiraProductConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var productId types.String
	req.Config.GetAttribute(ctx, path.Root("product_id"), &productId)
	var engagementId types.String
	req.Config.GetAttribute(ctx, path.Root("engagement_id"), &engagementId)
	if productId.IsNull() && engagementId.IsNull() {
		resp.Diagnostics.AddError("Invalid Resource", "The jira_product_configuration resource is invalid. Either product_id or engagement_id must be set.")
	}
}

type jiraProductConfigurationDataProvider struct{}

func (r jiraProductConfigurationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data jiraProductConfigurationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
