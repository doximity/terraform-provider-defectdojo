package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type terraformResourceData interface {
	id() types.String
	populate(defectdojoResource)
	defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error)
}

type defectdojoResource interface {
	createApiCall(context.Context, provider) (int, []byte, error)
	readApiCall(context.Context, provider, int) (int, []byte, error)
	updateApiCall(context.Context, provider, int) (int, []byte, error)
	deleteApiCall(context.Context, provider, int) (int, []byte, error)
}
type dataProvider interface {
	getData(context.Context, dataGetter) (terraformResourceData, diag.Diagnostics)
}

type terraformResource struct {
	provider provider
	dataProvider
}

type dataGetter interface {
	Get(context.Context, interface{}) diag.Diagnostics
}

func (r terraformResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	data, diags := r.getData(ctx, req.Config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ddResource, err := data.defectdojoResource(&resp.Diagnostics)
	if err != nil {
		return
	}

	statusCode, body, err := ddResource.createApiCall(ctx, r.provider)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 201 {
		data.populate(ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Creating Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%s", string(body)),
		)
		return
	}

	tflog.Trace(ctx, "created a JiraProductConfiguration")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r terraformResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	data, diags := r.getData(ctx, req.State)
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

	ddResource, err := data.defectdojoResource(&resp.Diagnostics)
	if err != nil {
		return
	}

	statusCode, body, err := ddResource.readApiCall(ctx, r.provider, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 200 {
		data.populate(ddResource)
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

func (r terraformResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	data, diags := r.getData(ctx, req.Plan)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.id().Null {
		resp.Diagnostics.AddError(
			"Could not Update Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.id().Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Update Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	ddResource, err := data.defectdojoResource(&resp.Diagnostics)
	if err != nil {
		return
	}

	statusCode, body, err := ddResource.updateApiCall(ctx, r.provider, idNumber)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 200 {
		data.populate(ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Updating Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r terraformResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	data, diags := r.getData(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.id().Null {
		resp.Diagnostics.AddError(
			"Could not Delete Resource",
			"The Id field was null but it is required to retrieve the product")
		return
	}

	idNumber, err := strconv.Atoi(data.id().Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Delete Resource",
			fmt.Sprintf("Error while parsing the Product ID from state: %s", err))
		return
	}

	ddResource, err := data.defectdojoResource(&resp.Diagnostics)
	if err != nil {
		return
	}

	statusCode, body, err := ddResource.deleteApiCall(ctx, r.provider, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode != 204 {
		resp.Diagnostics.AddError(
			"API Error Deleting Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r terraformResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
