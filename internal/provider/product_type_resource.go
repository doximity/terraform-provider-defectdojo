package provider

import (
	"context"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type productTypeResourceType struct{}

func (t productTypeResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product Type",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the Product Type",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the Product Type",
				Optional:            true,
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					stringDefault(""),
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

func (t productTypeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productTypeResource{
		terraformResource: terraformResource{
			provider:     provider,
			dataProvider: productTypeDataProvider{},
		},
	}, diags
}

type productTypeResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type productTypeDefectdojoResource struct {
	dd.ProductType
}

func (ddr *productTypeDefectdojoResource) createApiCall(ctx context.Context, p provider) (int, []byte, error) {
	reqBody := dd.ProductTypesCreateJSONRequestBody(ddr.ProductType)
	apiResp, err := p.client.ProductTypesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ProductType = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) readApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.ProductTypesRetrieveWithResponse(ctx, idNumber, &dd.ProductTypesRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.ProductType = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) updateApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	reqBody := dd.ProductTypesUpdateJSONRequestBody(ddr.ProductType)
	apiResp, err := p.client.ProductTypesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ProductType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) deleteApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.ProductTypesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type productTypeResource struct {
	terraformResource
}

type productTypeDataProvider struct{}

func (r productTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productTypeResourceData) id() types.String {
	return d.Id
}

func (d *productTypeResourceData) defectdojoResource() defectdojoResource {
	return &productTypeDefectdojoResource{
		ProductType: dd.ProductType{},
	}
}
