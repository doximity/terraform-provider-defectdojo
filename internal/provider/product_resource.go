package provider

import (
	"context"
	"fmt"
	"regexp"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/doximity/terraform-provider-defectdojo/internal/ref"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
					stringvalidator.RegexMatches(regexp.MustCompile(`^-?\d{0,13}(?:\.\d{0,2})?$`), `Must be a decimal number format, i.e. /^-?\d{0,13}(?:\.\d{0,2})?$/`),
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
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	ProductTypeId              types.Int64  `tfsdk:"product_type_id"`
	Id                         types.String `tfsdk:"id"`
	BusinessCriticality        types.String `tfsdk:"business_criticality"`
	EnableFullRiskAcceptance   types.Bool   `tfsdk:"enable_full_risk_acceptance"`
	EnableSimpleRiskAcceptance types.Bool   `tfsdk:"enable_skip_risk_acceptance"`
	ExternalAudience           types.Bool   `tfsdk:"external_audience"`
	InternetAccessible         types.Bool   `tfsdk:"internet_accessible"`
	Lifecycle                  types.String `tfsdk:"lifecycle"`
	Origin                     types.String `tfsdk:"origin"`
	Platform                   types.String `tfsdk:"platform"`
	ProdNumericGrade           types.Int64  `tfsdk:"prod_numeric_grade"`
	ProductManagerId           types.Int64  `tfsdk:"product_manager_id"`
	RegulationIds              []int64      `tfsdk:"regulation_ids"`
	Revenue                    types.String `tfsdk:"revenue"`
	Tags                       []string     `tfsdk:"tags"`
	TeamManagerId              types.Int64  `tfsdk:"team_manager_id"`
	TechnicalContactId         types.Int64  `tfsdk:"technical_contact_id"`
	UserRecords                types.Int64  `tfsdk:"user_records"`
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

func (d *productResourceData) populate(ddResource defectdojoResource) {
	tflog.Info(context.Background(), "populate")
	product := ddResource.(*productDefectdojoResource)

	d.Id = types.String{Value: fmt.Sprint(product.Id)}
	d.Name = types.String{Value: product.Name}
	d.Description = types.String{Value: product.Description}
	d.ProductTypeId = types.Int64{Value: int64(product.ProdType)}
	if product.ProdNumericGrade != nil {
		d.ProdNumericGrade = types.Int64{Value: int64(*product.ProdNumericGrade)}
	}
	if product.BusinessCriticality != nil {
		d.BusinessCriticality = types.String{Value: string(*product.BusinessCriticality)}
	}
	if product.Platform != nil {
		d.Platform = types.String{Value: string(*product.Platform)}
	}
	if product.Lifecycle != nil {
		d.Lifecycle = types.String{Value: string(*product.Lifecycle)}
	}
	if product.Origin != nil {
		d.Origin = types.String{Value: string(*product.Origin)}
	}
	if product.UserRecords != nil {
		d.UserRecords = types.Int64{Value: int64(*product.UserRecords)}
	}
	if product.Revenue != nil {
		d.Revenue = types.String{Value: string(*product.Revenue)}
	}
	if product.ExternalAudience != nil {
		d.ExternalAudience = types.Bool{Value: bool(*product.ExternalAudience)}
	}
	if product.InternetAccessible != nil {
		d.InternetAccessible = types.Bool{Value: bool(*product.InternetAccessible)}
	}
	if product.EnableSimpleRiskAcceptance != nil {
		d.EnableSimpleRiskAcceptance = types.Bool{Value: bool(*product.EnableSimpleRiskAcceptance)}
	}
	if product.EnableFullRiskAcceptance != nil {
		d.EnableFullRiskAcceptance = types.Bool{Value: bool(*product.EnableFullRiskAcceptance)}
	}
	if product.ProductManager != nil {
		d.ProductManagerId = types.Int64{Value: int64(*product.ProductManager)}
	}
	if product.TechnicalContact != nil {
		d.TechnicalContactId = types.Int64{Value: int64(*product.TechnicalContact)}
	}
	if product.TeamManager != nil {
		d.TeamManagerId = types.Int64{Value: int64(*product.TeamManager)}
	}
	if product.Regulations != nil && len(*product.Regulations) > 0 {
		var ids []int64
		for _, r := range *product.Regulations {
			ids = append(ids, int64(r))
		}
		d.RegulationIds = ids
		// must set to empty [] by default because
		// the API does
		if len(ids) == 0 {
			d.RegulationIds = make([]int64, 0)
		}
	}
	if product.Tags != nil {
		var tags []string
		for _, t := range *product.Tags {
			tags = append(tags, string(t))
		}
		d.Tags = tags
		// don't set to empty [] by default because
		// the API doesn't
	}
}

func (d *productResourceData) defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error) {
	tflog.Info(context.Background(), "defectdojoResource")
	product := dd.Product{
		ProdType:    int(d.ProductTypeId.Value),
		Description: d.Description.Value,
		Name:        d.Name.Value,
	}

	if !d.BusinessCriticality.IsNull() {
		product.BusinessCriticality = (*dd.ProductBusinessCriticality)(&d.BusinessCriticality.Value)
	}
	if !d.EnableFullRiskAcceptance.IsNull() {
		product.EnableFullRiskAcceptance = &d.EnableFullRiskAcceptance.Value
	}
	if !d.EnableSimpleRiskAcceptance.IsNull() {
		product.EnableSimpleRiskAcceptance = &d.EnableSimpleRiskAcceptance.Value
	}
	if !d.ExternalAudience.IsNull() {
		product.ExternalAudience = &d.ExternalAudience.Value
	}
	if !d.InternetAccessible.IsNull() {
		product.InternetAccessible = &d.InternetAccessible.Value
	}
	if !d.Lifecycle.IsNull() {
		product.Lifecycle = (*dd.ProductLifecycle)(&d.Lifecycle.Value)
	}
	if !d.Origin.IsNull() {
		product.Origin = (*dd.ProductOrigin)(&d.Origin.Value)
	}
	if !d.Platform.IsNull() {
		product.Platform = (*dd.ProductPlatform)(&d.Platform.Value)
	}
	if !d.ProdNumericGrade.IsNull() {
		product.ProdNumericGrade = ref.Of(int(d.ProdNumericGrade.Value))
	}
	if !d.ProductManagerId.IsNull() {
		product.ProductManager = ref.Of(int(d.ProductManagerId.Value))
	}
	if !d.Revenue.IsNull() {
		product.Revenue = &d.Revenue.Value
	}
	if !d.TeamManagerId.IsNull() {
		product.TeamManager = ref.Of(int(d.TeamManagerId.Value))
	}
	if !d.TechnicalContactId.IsNull() {
		product.TechnicalContact = ref.Of(int(d.TechnicalContactId.Value))
	}
	if !d.UserRecords.IsNull() {
		product.UserRecords = ref.Of(int(d.UserRecords.Value))
	}
	if len(d.RegulationIds) != 0 {
		var ids []int
		for _, id := range d.RegulationIds {
			ids = append(ids, int(id))
		}
		product.Regulations = &ids
	}
	if len(d.Tags) != 0 {
		var tags []string
		for _, tag := range d.Tags {
			tags = append(tags, tag)
		}
		product.Tags = &tags
	}

	return &productDefectdojoResource{
		Product: product,
	}, nil
}
