package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type productTypeDataSourceType struct{}

func (t productTypeDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Product Type",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the Product Type",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the Product Type",
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t productTypeDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productTypeDataSource{
		provider: provider,
	}, diags
}

type productTypeDataSourceData struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

type productTypeDataSource struct {
	provider provider
}

func (d productTypeDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data productTypeDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var (
		params dd.ProductTypesListParams
	)
	if !data.Id.Null {
		idNumber, err := strconv.Atoi(data.Id.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"The id field could not be parsed into an integer")
			return
		} else {
			params.Id = &idNumber
		}
	}

	if !data.Name.Null {
		params.Name = &data.Name.Value
	}

	apiResp, err := d.provider.client.ProductTypesListWithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	var pt dd.ProductType
	if *apiResp.JSON200.Count == 0 {
		resp.Diagnostics.AddError(
			"Could not Retrieve Data Resource",
			"No Product Types matched the given parameters.")
	} else if *apiResp.JSON200.Count > 1 {
		body, _ := ioutil.ReadAll(apiResp.HTTPResponse.Body)
		resp.Diagnostics.AddError(
			"Could not Retrieve Data Resource",
			fmt.Sprintf("%d Product Types matched the given parameters.\n\nResponse:\n\n%s", *apiResp.JSON200.Count, body))
	} else {
		pt = (*apiResp.JSON200.Results)[0]

		data.Id = types.String{Value: fmt.Sprintf("%d", pt.Id)}
		data.Name = types.String{Value: pt.Name}
		data.Description = types.String{Value: *pt.Description}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
