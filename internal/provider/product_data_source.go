package provider

import (
	"context"
	"log"
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
		MarkdownDescription: "Example data source",

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

	log.Printf("got here")

	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("got here")

	idNumber, err := strconv.Atoi(data.Id.Value)
	if err != nil {
		panic(err)
	}

	apiResp, err := d.provider.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})

	if err != nil {
		panic(err)
	}

	data.Name = types.String{Value: apiResp.JSON200.Name}
	data.Description = types.String{Value: apiResp.JSON200.Description}
	data.ProductTypeId = types.Int64{Value: int64(apiResp.JSON200.ProdType)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
