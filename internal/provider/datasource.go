package provider

import (
	"context"
	"fmt"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type terraformDatasource struct {
	client *dd.ClientWithResponses
	dataProvider
}

func (r *terraformDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r terraformDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, diags := r.getData(ctx, req.Config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.id().IsNull() {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.id().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.readApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 200 {
		populateResourceData(ctx, &diags, &data, ddResource)
	} else if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	} else {
		resp.Diagnostics.AddError(
			"API Error Retrieving Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
