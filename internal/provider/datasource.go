package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type terraformDatasource struct {
	provider provider
	dataProvider
}

func (r terraformDatasource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	data, diags := r.getData(ctx, req.Config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.id().Null {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.id().Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.readApiCall(ctx, r.provider, idNumber)
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
