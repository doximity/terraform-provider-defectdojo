package provider

import (
	"context"
	"fmt"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"product_type_id": {
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
				Type:                types.Int64Type,
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
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ProductTypeId types.Int64  `tfsdk:"product_type_id"`
	Id            types.String `tfsdk:"id"`
}

type productDefectdojoResource struct {
	dd.Product
}

func (ddr *productDefectdojoResource) createApiCall(ctx context.Context, p provider) (int, []byte, error) {
	reqBody := dd.ProductsCreateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.Product = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) readApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) updateApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	reqBody := dd.ProductsUpdateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) deleteApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	apiResp, err := p.client.ProductsDestroyWithResponse(ctx, idNumber)
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
	product := ddResource.(*productDefectdojoResource)
	d.Id = types.String{Value: fmt.Sprint(product.Id)}
	d.Name = types.String{Value: product.Name}
	d.Description = types.String{Value: product.Description}
	d.ProductTypeId = types.Int64{Value: int64(product.ProdType)}
}

func (d *productResourceData) defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error) {
	product := dd.Product{
		ProdType:    int(d.ProductTypeId.Value),
		Description: d.Description.Value,
		Name:        d.Name.Value,
	}
	return &productDefectdojoResource{
		Product: product,
	}, nil
}
