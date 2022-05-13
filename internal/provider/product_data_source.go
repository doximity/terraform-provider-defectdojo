package provider

import (
	"context"
	"fmt"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type productDataSourceType struct{}

func (t productDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Product",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the Product",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the Product",
				Computed:            true,
				Type:                types.StringType,
			},
			"product_type_id": {
				MarkdownDescription: "The ID of the Product Type",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"id": {
				MarkdownDescription: "Identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t productDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productDataSource{
		provider: provider,
	}, diags
}

type productDataSourceData struct {
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ProductTypeId types.Int64  `tfsdk:"product_type_id"`
	Id            types.String `tfsdk:"id"`
}

type productDataSource struct {
	provider provider
}

func (d productDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data productDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			"The id field could not be parsed into an integer")
	}

	apiResp, err := d.provider.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	data.Name = types.String{Value: apiResp.JSON200.Name}
	data.Description = types.String{Value: apiResp.JSON200.Description}
	data.ProductTypeId = types.Int64{Value: int64(apiResp.JSON200.ProdType)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
