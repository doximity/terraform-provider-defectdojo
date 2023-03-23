package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/doximity/terraform-provider-defectdojo/internal/ref"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (t productTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Product Type. You can specify either the `id` or the `name` to look up the Product Type.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product Type",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Product Type",
				Computed:            true,
			},
			"critical_product": schema.BoolAttribute{
				MarkdownDescription: "Is this a critical Product Type",
				Computed:            true,
			},
			"key_product": schema.BoolAttribute{
				MarkdownDescription: "Is this a key Product Type",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
		},
	}
}

type productTypeDataSourceData struct {
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CriticalProduct types.Bool   `tfsdk:"critical_product"`
	KeyProduct      types.Bool   `tfsdk:"key_product"`
	Id              types.String `tfsdk:"id"`
}

type productTypeDataSource struct {
	client *dd.ClientWithResponses
}

func (d productTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type"
}

func NewProductTypeDataSource() datasource.DataSource {
	return &productTypeDataSource{}
}

func (r *productTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dd.ClientWithResponses)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected dd.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (d productTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productTypeDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var (
		params dd.ProductTypesListParams
	)
	if !data.Id.IsNull() {
		idNumber, err := strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"The id field could not be parsed into an integer")
			return
		} else {
			params.Id = &idNumber
		}
	}

	if !data.Name.IsNull() {
		params.Name = ref.Of(data.Name.ValueString())
	}

	apiResp, err := d.client.ProductTypesListWithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		var pt dd.ProductType
		if *apiResp.JSON200.Count == 0 {
			resp.Diagnostics.AddError(
				"Could not Retrieve Data Resource",
				"No Product Types matched the given parameters.")
			return
		} else if *apiResp.JSON200.Count > 1 {
			body, _ := ioutil.ReadAll(apiResp.HTTPResponse.Body)
			resp.Diagnostics.AddError(
				"Could not Retrieve Data Resource",
				fmt.Sprintf("%d Product Types matched the given parameters.\n\nResponse:\n\n%s", *apiResp.JSON200.Count, body))
			return
		} else {
			pt = (*apiResp.JSON200.Results)[0]

			data.Id = types.StringValue(fmt.Sprintf("%d", pt.Id))
			data.Name = types.StringValue(pt.Name)
			data.Description = types.StringValue(*pt.Description)
			data.CriticalProduct = types.BoolValue(*pt.CriticalProduct)
			data.KeyProduct = types.BoolValue(*pt.KeyProduct)
		}
	} else {
		body, _ := ioutil.ReadAll(apiResp.HTTPResponse.Body)

		resp.Diagnostics.AddError(
			"API Error Retrieving Data Source",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%+v", body),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
