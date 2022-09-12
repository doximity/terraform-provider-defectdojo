package provider

import (
	"context"

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
				Optional:            true,
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
			},
			"platform": {
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Computed:            true,
				Type:                types.StringType,
			},
			"lifecycle": {
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Computed:            true,
				Type:                types.StringType,
			},
			"origin": {
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Computed:            true,
				Type:                types.StringType,
			},
			"user_records": {
				MarkdownDescription: "Estimate the number of user records within the application.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"revenue": {
				MarkdownDescription: "Estimate the application's revenue.",
				Computed:            true,
				Type:                types.StringType,
			},
			"external_audience": {
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"internet_accessible": {
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"enable_skip_risk_acceptance": {
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"enable_full_risk_acceptance": {
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"product_manager_id": {
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"technical_contact_id": {
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"team_manager_id": {
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"regulation_ids": {
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"tags": {
				MarkdownDescription: "Tags to apply to the product",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

func (t productDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productDataSource{
		terraformDatasource: terraformDatasource{
			provider:     provider,
			dataProvider: productDataProvider{},
		},
	}, diags
}

type productDataSource struct {
	terraformDatasource
}
