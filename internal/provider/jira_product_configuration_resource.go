package provider

import (
	"context"
	"fmt"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/doximity/terraform-provider-defectdojo/internal/ref"
)

type jiraProductConfigurationResourceType struct{}

func (t jiraProductConfigurationResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Jira Product Configuration is the connection between a Product and a Jira Instance. It defines the Product's settings for pushing Findings to Jira.",

		Attributes: map[string]tfsdk.Attribute{
			"project_key": {
				MarkdownDescription: "The Jira Project Key",
				Optional:            true,
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					stringDefault(""),
				},
			},

			"issue_template_dir": {
				MarkdownDescription: "The folder containing Django templates used to render the JIRA issue description. Leave empty to use the default jira_full templates.",
				Optional:            true,
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					stringDefault(""),
				},
			},

			"push_all_issues": {
				MarkdownDescription: "Automatically maintain parity with JIRA. Always create and update JIRA tickets for findings in this Product.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},

			"enable_engagement_epic_mapping": {
				MarkdownDescription: "Whether to map engagements to epics in Jira",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},

			"push_notes": {
				MarkdownDescription: "Whether to push notes to Jira",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},

			"product_jira_sla_notification": {
				MarkdownDescription: "Send SLA notifications as comments",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},

			"risk_acceptance_expiration_notification": {
				MarkdownDescription: "Send Risk Acceptance expiration notifications as comments",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},

			"jira_instance_id": {
				MarkdownDescription: "The ID of the Jira Instance to use for this Product",
				Optional:            true,
				Type:                types.StringType,
			},

			"product_id": {
				MarkdownDescription: "The ID of the Product to configure. Although optional, either the Product ID or the Engagement ID must be defined to create a Jira Product Configuration.",
				Optional:            true,
				Type:                types.StringType,
			},

			"engagement_id": {
				MarkdownDescription: "The ID of the Engagement. Although optional, either the Product ID or the Engagement ID must be defined to create a Jira Product Configuration.",
				Optional:            true,
				Type:                types.StringType,
			},

			"id": {
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType, // the id (for import purposes) MUST be a string
			},
		},
	}, nil
}

func (t jiraProductConfigurationResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraProductConfigurationResource{
		terraformResource: terraformResource{
			provider:     provider,
			dataProvider: jiraProductConfigurationDataProvider{},
		},
	}, diags
}

type jiraProductConfigurationResourceData struct {
	ProjectKey                           types.String `tfsdk:"project_key"`
	IssueTemplateDir                     types.String `tfsdk:"issue_template_dir"`
	PushAllIssues                        types.Bool   `tfsdk:"push_all_issues"`
	EnableEngagementEpicMapping          types.Bool   `tfsdk:"enable_engagement_epic_mapping"`
	PushNotes                            types.Bool   `tfsdk:"push_notes"`
	ProductJiraSlaNotification           types.Bool   `tfsdk:"product_jira_sla_notification"`
	RiskAcceptanceExpirationNotification types.Bool   `tfsdk:"risk_acceptance_expiration_notification"`
	JiraInstance                         types.String `tfsdk:"jira_instance_id"`
	Product                              types.String `tfsdk:"product_id"`
	Engagement                           types.String `tfsdk:"engagement_id"`
	Id                                   types.String `tfsdk:"id"`
}

type jiraProductConfigurationDefectdojoResource struct {
	dd.JIRAProject
}

func (ddr *jiraProductConfigurationDefectdojoResource) createApiCall(ctx context.Context, p provider) (int, []byte, error) {
	reqBody := dd.JiraProductConfigurationsCreateJSONRequestBody(ddr.JIRAProject)
	apiResp, err := p.client.JiraProductConfigurationsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.JIRAProject = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) readApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.JiraProductConfigurationsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.JIRAProject = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) updateApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	reqBody := dd.JiraProductConfigurationsUpdateJSONRequestBody(ddr.JIRAProject)
	apiResp, err := p.client.JiraProductConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.JIRAProject = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *jiraProductConfigurationDefectdojoResource) deleteApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.JiraProductConfigurationsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *jiraProductConfigurationResourceData) id() types.String {
	return d.Id
}

func (d *jiraProductConfigurationResourceData) populate(ddResource defectdojoResource) {
	jiraProject := ddResource.(*jiraProductConfigurationDefectdojoResource)

	d.Id = types.String{Value: fmt.Sprint(jiraProject.Id)}

	if jiraProject.Product != nil {
		d.Product = types.String{Value: fmt.Sprint(*jiraProject.Product)}
	}
	if jiraProject.Engagement != nil {
		d.Engagement = types.String{Value: fmt.Sprint(*jiraProject.Engagement)}
	}
	if jiraProject.JiraInstance != nil {
		d.JiraInstance = types.String{Value: fmt.Sprint(*jiraProject.JiraInstance)}
	}
	if jiraProject.RiskAcceptanceExpirationNotification != nil {
		d.RiskAcceptanceExpirationNotification = types.Bool{Value: *jiraProject.RiskAcceptanceExpirationNotification}
	}
	if jiraProject.ProductJiraSlaNotification != nil {
		d.ProductJiraSlaNotification = types.Bool{Value: *jiraProject.ProductJiraSlaNotification}
	}
	if jiraProject.PushNotes != nil {
		d.PushNotes = types.Bool{Value: *jiraProject.PushNotes}
	}
	if jiraProject.EnableEngagementEpicMapping != nil {
		d.EnableEngagementEpicMapping = types.Bool{Value: *jiraProject.EnableEngagementEpicMapping}
	}
	if jiraProject.PushAllIssues != nil {
		d.PushAllIssues = types.Bool{Value: *jiraProject.PushAllIssues}
	}
	if jiraProject.IssueTemplateDir != nil {
		d.IssueTemplateDir = types.String{Value: *jiraProject.IssueTemplateDir}
	}
	if jiraProject.ProjectKey != nil {
		d.ProjectKey = types.String{Value: *jiraProject.ProjectKey}
	}
}

func (d *jiraProductConfigurationResourceData) defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error) {
	var productIdNumber, engagementIdNumber, jiraInstanceIdNumber int
	var err error

	if !d.Product.IsNull() {
		productIdNumber, err = strconv.Atoi(d.Product.Value)
		if err != nil {
			diags.AddError(
				"Could not Create Resource",
				fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
			return nil, err
		}
	}
	if !d.Engagement.IsNull() {
		engagementIdNumber, err = strconv.Atoi(d.Engagement.Value)
		if err != nil {
			diags.AddError(
				"Could not Create Resource",
				fmt.Sprintf("Error while parsing the Engagement ID from state: %s", err))
			return nil, err
		}
	}
	if !d.JiraInstance.IsNull() {
		jiraInstanceIdNumber, err = strconv.Atoi(d.JiraInstance.Value)
		if err != nil {
			diags.AddError(
				"Could not Create Resource",
				fmt.Sprintf("Error while parsing the Jira Instance ID from state: %s", err))
			return nil, err
		}
	}

	ret := dd.JIRAProject{
		RiskAcceptanceExpirationNotification: ref.Of(d.RiskAcceptanceExpirationNotification.Value),
		ProductJiraSlaNotification:           ref.Of(d.ProductJiraSlaNotification.Value),
		PushNotes:                            ref.Of(d.PushNotes.Value),
		EnableEngagementEpicMapping:          ref.Of(d.EnableEngagementEpicMapping.Value),
		PushAllIssues:                        ref.Of(d.PushAllIssues.Value),
		IssueTemplateDir:                     ref.Of(d.IssueTemplateDir.Value),
		ProjectKey:                           ref.Of(d.ProjectKey.Value),
	}

	if productIdNumber != 0 {
		ret.Product = ref.Of(productIdNumber)
	}
	if engagementIdNumber != 0 {
		ret.Engagement = ref.Of(engagementIdNumber)
	}
	if jiraInstanceIdNumber != 0 {
		ret.JiraInstance = ref.Of(jiraInstanceIdNumber)
	}

	ddResource := jiraProductConfigurationDefectdojoResource{
		JIRAProject: ret,
	}

	return &ddResource, nil
}

type jiraProductConfigurationResource struct {
	terraformResource
}

func (r jiraProductConfigurationResource) ValidateConfig(ctx context.Context, req tfsdk.ValidateResourceConfigRequest, resp *tfsdk.ValidateResourceConfigResponse) {
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
