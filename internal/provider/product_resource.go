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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (t productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Product",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A[^\s].*[^\s]\z`), "The description must not have leading or trailing whitespace"),
				},
			},
			"prod_numeric_grade": schema.Int64Attribute{
				MarkdownDescription: "The Numeric Grade of the Product",
				Optional:            true,
			},
			"business_criticality": schema.StringAttribute{
				MarkdownDescription: "The Business Criticality of the Product. Valid values are: 'very high', 'high', 'medium', 'low', 'very low', 'none'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("very high", "high", "medium", "low", "very low", "none", ""),
				},
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("web service", "desktop", "iot", "mobile", "web", ""),
				},
			},
			"life_cycle": schema.StringAttribute{
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("construction", "production", "retirement", ""),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("third party library", "purchased", "contractor", "internal", "open source", "outsourced", ""),
				},
			},
			"user_records": schema.Int64Attribute{
				MarkdownDescription: "Estimate the number of user records within the application.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"revenue": schema.StringAttribute{
				MarkdownDescription: "Estimate the application's revenue.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A-?\d{0,13}(?:\.\d{0,2})?\z`), `Must be a decimal number format, i.e. /^-?\d{0,13}(?:\.\d{0,2})?$/`),
				},
			},
			"external_audience": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
			},
			"internet_accessible": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
			},
			"enable_skip_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
			},
			"enable_full_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
			},
			"product_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Optional:            true,
			},
			"technical_contact_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Optional:            true,
			},
			"team_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Optional:            true,
			},
			"regulation_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"product_type_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags to apply to the product",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`\A[a-z]+\z`), "Tags must be lower case values"),
					),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
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
	Lifecycle                  types.String `tfsdk:"life_cycle" ddField:"Lifecycle"`
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

func (ddr *productDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := dd.ProductsCreateJSONRequestBody(ddr.Product)
	apiResp, err := client.ProductsCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Product = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.ProductsUpdateJSONRequestBody(ddr.Product)
	apiResp, err := client.ProductsUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.ProductsDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type productResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &productResource{}
var _ resource.ResourceWithImportState = &productResource{}
var _ resource.ResourceWithConfigure = &productResource{}

func NewProductResource() resource.Resource {
	return &productResource{
		terraformResource: terraformResource{
			dataProvider: productDataProvider{},
		},
	}
}

func (r productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
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

func (d *productResourceData) defectdojoResource() defectdojoResource {
	return &productDefectdojoResource{
		Product: dd.Product{},
	}
}
