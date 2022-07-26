package provider

import (
	"context"
	"fmt"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

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
				MarkdownDescription: "The ID of the Product to configure",
				Required:            true,
				Type:                types.StringType,
			},

			"engagement_id": {
				MarkdownDescription: "The ID of the Engagement. Can be empty",
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
		provider: provider,
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

func (d *jiraProductConfigurationResourceData) populate(jiraProject *dd.JIRAProject) {
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

func (d *jiraProductConfigurationResourceData) jiraProject(diags *diag.Diagnostics) (*dd.JIRAProject, error) {
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

	return &ret, nil
}

type jiraProductConfigurationResource struct {
	provider provider
}

func (r jiraProductConfigurationResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data jiraProductConfigurationResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	jiraProject, err := data.jiraProject(&resp.Diagnostics)
	if err != nil {
		return
	}
	reqBody := dd.JiraProductConfigurationsCreateJSONRequestBody(*jiraProject)
	apiResp, err := r.provider.client.JiraProductConfigurationsCreateWithResponse(ctx, reqBody)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 201 {

		data.populate(apiResp.JSON201)
	} else {
		resp.Diagnostics.AddError(
			"API Error Creating Resource",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%s", string(apiResp.Body)),
		)
		return
	}

	tflog.Trace(ctx, "created a JiraProductConfiguration")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r jiraProductConfigurationResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data jiraProductConfigurationResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.Null {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	apiResp, err := r.provider.client.JiraProductConfigurationsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		data.populate(apiResp.JSON200)
	} else if apiResp.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	} else {
		resp.Diagnostics.AddError(
			"API Error Retrieving Resource",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(apiResp.Body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r jiraProductConfigurationResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data jiraProductConfigurationResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.Null {
		resp.Diagnostics.AddError(
			"Could not Update Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Update Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	jiraProject, err := data.jiraProject(&resp.Diagnostics)
	reqBody := dd.JiraProductConfigurationsUpdateJSONRequestBody(*jiraProject)
	apiResp, err := r.provider.client.JiraProductConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		data.populate(apiResp.JSON200)
	} else {
		resp.Diagnostics.AddError(
			"API Error Updating Resource",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(apiResp.Body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r jiraProductConfigurationResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data jiraProductConfigurationResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.Null {
		resp.Diagnostics.AddError(
			"Could not Delete Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Delete Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	apiResp, err := r.provider.client.JiraProductConfigurationsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"API Error Deleting Resource",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(apiResp.Body)),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r jiraProductConfigurationResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
