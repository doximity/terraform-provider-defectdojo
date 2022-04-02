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
		provider: provider,
	}, diags
}

type productResourceData struct {
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ProductTypeId types.Int64  `tfsdk:"product_type_id"`
	Id            types.String `tfsdk:"id"`
}

type productResource struct {
	provider provider
}

func (r productResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data productResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.provider.client.ProductsCreateWithResponse(ctx, dd.ProductsCreateJSONRequestBody{
		ProdType:    int(data.ProductTypeId.Value),
		Description: data.Description.Value,
		Name:        data.Name.Value,
	})

	if err != nil {
		panic(err)
	}

	if apiResp.StatusCode() == 201 {
		data.Id = types.String{Value: fmt.Sprint(apiResp.JSON201.Id)}
		data.Name = types.String{Value: apiResp.JSON201.Name}
		data.Description = types.String{Value: apiResp.JSON201.Description}
		data.ProductTypeId = types.Int64{Value: int64(apiResp.JSON201.ProdType)}
	} else {
		// ??
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a product")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r productResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data productResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		panic(err)
	}

	apiResp, err := r.provider.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	if err != nil {
		panic(err)
	}

	if apiResp.StatusCode() == 200 {
		data.Name = types.String{Value: apiResp.JSON200.Name}
		data.Description = types.String{Value: apiResp.JSON200.Description}
		data.ProductTypeId = types.Int64{Value: int64(apiResp.JSON200.ProdType)}
	} else {
		// ??
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r productResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data productResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		panic(err)
	}

	apiResp, err := r.provider.client.ProductsUpdateWithResponse(ctx, idNumber, dd.ProductsUpdateJSONRequestBody{
		ProdType:    int(data.ProductTypeId.Value),
		Description: data.Description.Value,
		Name:        data.Name.Value,
	})

	if err != nil {
		panic(err)
	}

	if apiResp.StatusCode() == 200 {
		data.Id = types.String{Value: fmt.Sprint(apiResp.JSON200.Id)}
		data.Name = types.String{Value: apiResp.JSON200.Name}
		data.Description = types.String{Value: apiResp.JSON200.Description}
		data.ProductTypeId = types.Int64{Value: int64(apiResp.JSON200.ProdType)}
	} else {
		// ??
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r productResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data productResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		panic(err)
	}

	apiResp, err := r.provider.client.ProductsDestroy(ctx, idNumber)
	if err != nil {
		panic(err)
	}

	if apiResp.StatusCode == 204 {
		// success
	} else {
		// ??
	}

	resp.State.RemoveResource(ctx)
}

func (r productResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
