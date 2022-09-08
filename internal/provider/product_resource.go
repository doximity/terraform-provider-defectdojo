package provider

import (
	"context"
	"fmt"
	"regexp"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type productResourceType struct{}

func (t productResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the Product",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the Product",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A[^\s].*[^\s]\z`), "The description must not have leading or trailing whitespace"),
				},
			},
			"prod_numeric_grade": {
				MarkdownDescription: "The Numeric Grade of the Product",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"business_criticality": {
				MarkdownDescription: "The Business Criticality of the Product. Valid values are: 'very high', 'high', 'medium', 'low', 'very low', 'none'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("very high", "high", "medium", "low", "very low", "none", ""),
				},
			},
			"platform": {
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("web service", "desktop", "iot", "mobile", "web", ""),
				},
			},
			"lifecycle": {
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("construction", "production", "retirement", ""),
				},
			},
			"origin": {
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("third party library", "purchased", "contractor", "internal", "open source", "outsourced", ""),
				},
			},
			"user_records": {
				MarkdownDescription: "Estimate the number of user records within the application.",
				Optional:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					int64validator.AtLeast(0),
				},
			},
			"revenue": {
				MarkdownDescription: "Estimate the application's revenue.",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A-?\d{0,13}(?:\.\d{0,2})?\z`), `Must be a decimal number format, i.e. /^-?\d{0,13}(?:\.\d{0,2})?$/`),
				},
			},
			"external_audience": {
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"internet_accessible": {
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"enable_skip_risk_acceptance": {
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"enable_full_risk_acceptance": {
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"product_manager_id": {
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"technical_contact_id": {
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"team_manager_id": {
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"regulation_ids": {
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"product_type_id": {
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
				Type:                types.Int64Type,
			},
			"tags": {
				MarkdownDescription: "Tags to apply to the product",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Validators: []tfsdk.AttributeValidator{
					setvalidator.ValuesAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`\A[a-z]+\z`), "Tags must be lower case values"),
					),
				},
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

func (t productResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productResource{
		terraformResource: terraformResource{
			provider:     provider,
			dataProvider: productDataProvider{},
		},
	}, diags
}

type productResourceData struct {
	Name                       types.String `tfsdk:"name" ddField:"Name"`
	Description                types.String `tfsdk:"description" ddField:"Description"`
	ProductTypeId              types.Int64  `tfsdk:"product_type_id" ddField:"ProdType"`
	Id                         types.String `tfsdk:"id" ddField:"Id"`
	BusinessCriticality        types.String `tfsdk:"business_criticality" ddField:"BusinessCriticality"`
	EnableFullRiskAcceptance   types.Bool   `tfsdk:"enable_full_risk_acceptance" ddField:"EnableFullRiskAcceptance"`
	EnableSimpleRiskAcceptance types.Bool   `tfsdk:"enable_skip_risk_acceptance" ddField:"EnableSimpleRiskAcceptance"`
	ExternalAudience           types.Bool   `tfsdk:"external_audience" ddField:"ExternalAudience"`
	InternetAccessible         types.Bool   `tfsdk:"internet_accessible" ddField:"InternetAccessible"`
	Lifecycle                  types.String `tfsdk:"lifecycle" ddField:"Lifecycle"`
	Origin                     types.String `tfsdk:"origin" ddField:"Origin"`
	Platform                   types.String `tfsdk:"platform" ddField:"Platform"`
	ProdNumericGrade           types.Int64  `tfsdk:"prod_numeric_grade" ddField:"ProdNumericGrade"`
	ProductManagerId           types.Int64  `tfsdk:"product_manager_id" ddField:"ProductManager"`
	RegulationIds              types.Set    `tfsdk:"regulation_ids" ddField:"Regulations"`
	Revenue                    types.String `tfsdk:"revenue" ddField:"Revenue"`
	Tags                       types.Set    `tfsdk:"tags" ddField:"Tags"`
	TeamManagerId              types.Int64  `tfsdk:"team_manager_id" ddField:"TeamManager"`
	TechnicalContactId         types.Int64  `tfsdk:"technical_contact_id" ddField:"TechnicalContact"`
	UserRecords                types.Int64  `tfsdk:"user_records" ddField:"UserRecords"`
}

type productDefectdojoResource struct {
	dd.Product
}

func (ddr *productDefectdojoResource) createApiCall(ctx context.Context, p provider) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := dd.ProductsCreateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Product = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) readApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := p.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) updateApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.ProductsUpdateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) deleteApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := p.client.ProductsDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type productResource struct {
	terraformResource
}

type productDataProvider struct{}

func (r productDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productResourceData) id() types.String {
	return d.Id
}

func (d *productResourceData) defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error) {
	return &productDefectdojoResource{
		Product: dd.Product{},
	}, nil
}
