package provider

import (
	"context"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (t productTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product Type",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product Type",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Product Type",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringDefault(""),
				},
			},
			"critical_product": schema.BoolAttribute{
				MarkdownDescription: "Is this a critical Product Type",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
			},
			"key_product": schema.BoolAttribute{
				MarkdownDescription: "Is this a key Product Type",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
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

type productTypeResourceData struct {
	Name            types.String `tfsdk:"name" ddField:"Name"`
	Description     types.String `tfsdk:"description" ddField:"Description"`
	CriticalProduct types.Bool   `tfsdk:"critical_product" ddField:"Critical Product"`
	KeyProduct      types.Bool   `tfsdk:"key_product" ddField:"Key Product"`
	Id              types.String `tfsdk:"id" ddField:"Id"`
}

type productTypeDefectdojoResource struct {
	dd.ProductType
}

func (ddr *productTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := dd.ProductTypesCreateJSONRequestBody(ddr.ProductType)
	apiResp, err := client.ProductTypesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ProductType = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypesRetrieveWithResponse(ctx, idNumber, &dd.ProductTypesRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.ProductType = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := dd.ProductTypesUpdateJSONRequestBody(ddr.ProductType)
	apiResp, err := client.ProductTypesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ProductType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type productTypeResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &productTypeResource{}
var _ resource.ResourceWithImportState = &productTypeResource{}

func NewProductTypeResource() resource.Resource {
	return &productTypeResource{
		terraformResource: terraformResource{
			dataProvider: productTypeDataProvider{},
		},
	}
}

func (r productTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type"
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
