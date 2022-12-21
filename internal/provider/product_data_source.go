package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type productDataSource struct {
	terraformDatasource
}

func (t productDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Product",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Product",
				Computed:            true,
			},
			"product_type_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"prod_numeric_grade": schema.Int64Attribute{
				MarkdownDescription: "The Numeric Grade of the Product",
				Optional:            true,
			},
			"business_criticality": schema.StringAttribute{
				MarkdownDescription: "The Business Criticality of the Product. Valid values are: 'very high', 'high', 'medium', 'low', 'very low', 'none'",
				Optional:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Computed:            true,
			},
			"life_cycle": schema.StringAttribute{
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Computed:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Computed:            true,
			},
			"user_records": schema.Int64Attribute{
				MarkdownDescription: "Estimate the number of user records within the application.",
				Computed:            true,
			},
			"revenue": schema.StringAttribute{
				MarkdownDescription: "Estimate the application's revenue.",
				Computed:            true,
			},
			"external_audience": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Computed:            true,
			},
			"internet_accessible": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Computed:            true,
			},
			"enable_skip_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Computed:            true,
			},
			"enable_full_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Computed:            true,
			},
			"product_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Computed:            true,
			},
			"technical_contact_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Computed:            true,
			},
			"team_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Computed:            true,
			},
			"regulation_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags to apply to the product",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d productDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &productDataSource{}

func NewProductDataSource() datasource.DataSource {
	return &productDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productDataProvider{},
		},
	}
}
