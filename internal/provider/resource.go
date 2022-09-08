package provider

import (
	"context"
	"fmt"
	"reflect"
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
	populateDefectdojoResource(data, &ddResource)

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
	populateDefectdojoResource(data, &ddResource)

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
	populateDefectdojoResource(data, &ddResource)

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

func populateDefectdojoResource(resourceData terraformResourceData, ddResource *defectdojoResource) {

	resourceVal := reflect.ValueOf(resourceData).Elem()
	resourceType := resourceVal.Type()
	// fmt.Printf("resourceVal: %s\n", resourceVal)
	// fmt.Printf("resourceType: %s\n", resourceType)

	ddVal := reflect.ValueOf(*ddResource).Elem()

	for i := 0; i < resourceVal.NumField(); i++ {
		fieldDescriptor := resourceType.Field(i)
		tag := fieldDescriptor.Tag
		ddFieldName := tag.Get("ddField")
		if ddFieldName != "" {
			fieldValue := resourceVal.Field(i)
			ddFieldDescriptor, _ := ddVal.Type().FieldByName(ddFieldName)
			ddFieldValue := ddVal.FieldByName(ddFieldName)

			// fmt.Printf("ddFieldDescriptor: Kind = %s, Name = %s\n", ddFieldDescriptor.Type.Kind(), ddFieldDescriptor.Name)
			// fmt.Printf("fieldDescriptor: Kind = %s, Name = %s, type = %s\n", fieldDescriptor.Type.Kind(), fieldDescriptor.Name, fieldDescriptor.Type)

			switch fieldDescriptor.Type {

			case typeOfTypesString:
				if ddFieldDescriptor.Type.Kind() == reflect.String {
					// if the destination field is a string, we can grab the `Value` field and assign it directly
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					if !srcIsNull {
						ddFieldValue.Set(fieldValue.FieldByName("Value"))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
					// the destination field is a *string (or compatible/alias) so we have to set it to a pointer
					// if the source is Null:true, then we set to to a nil pointer, but we still have to make sure it
					// is a nil pointer of the correct type
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					if !srcIsNull {
						destType := ddFieldDescriptor.Type.Elem()
						destVal := reflect.New(destType)
						destVal.Elem().Set(fieldValue.FieldByName("Value").Convert(destType))
						ddFieldValue.Set(destVal)
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int {
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					zero := 0
					if !srcIsNull {
						srcVal := fieldValue.FieldByName("Value")
						strVal := srcVal.Interface().(string)
						intVal, err := strconv.Atoi(strVal)
						if err == nil {
							ddFieldValue.Set(reflect.ValueOf(zero))
						}
						ddFieldValue.Set(reflect.ValueOf(intVal))
					}
				} else {
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesBool:
				if ddFieldDescriptor.Type.Kind() == reflect.Bool {
					// if the destination field is a bool, we can grab the `Value` field and assign it directly
					ddFieldValue.Set(fieldValue.FieldByName("Value"))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Bool {
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					if !srcIsNull {
						destType := ddFieldDescriptor.Type.Elem()
						destVal := reflect.New(destType)
						destVal.Elem().Set(fieldValue.FieldByName("Value").Convert(destType))
						ddFieldValue.Set(destVal)
					} else {
						ddFieldValue.Set(reflect.New(ddFieldDescriptor.Type).Elem())
					}
				} else {
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesInt64:
				if ddFieldDescriptor.Type.Kind() == reflect.Int {
					// if the destination field is an int, we can grab the `Value` field and cast and assign it directly
					destVal := reflect.New(ddFieldDescriptor.Type)
					destVal.Elem().Set(fieldValue.FieldByName("Value").Convert(ddFieldDescriptor.Type))
					ddFieldValue.Set(destVal.Elem())
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					// the destination field is a *int so we have to set it to a pointer
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					if !srcIsNull {
						destType := ddFieldDescriptor.Type.Elem()
						destVal := reflect.New(destType)
						destVal.Elem().Set(fieldValue.FieldByName("Value").Convert(destType))
						ddFieldValue.Set(destVal)
					} else {
						ddFieldValue.Set(reflect.New(ddFieldDescriptor.Type).Elem())
					}
				} else {
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesSet:
				if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the source field is a pointer to a slice
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int

						if fieldValue.FieldByName("Null").Bool() {
							ints := make([]int, 0)
							destVal := reflect.New(ddFieldDescriptor.Type.Elem())
							destVal.Elem().Set(reflect.ValueOf(ints))
							ddFieldValue.Set(destVal)
						} else {
							int64s := []int64{}
							_ = fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &int64s, false)
							ints := []int{}
							for _, val := range int64s {
								ints = append(ints, (int)(val))
							}
							if ints == nil {
								ints = make([]int, 0)
							}
							destVal := reflect.New(ddFieldDescriptor.Type.Elem())
							destVal.Elem().Set(reflect.ValueOf(ints))
							ddFieldValue.Set(destVal)
						}
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string

						if fieldValue.FieldByName("Null").Bool() {
							strings := make([]string, 0)
							destVal := reflect.New(ddFieldDescriptor.Type.Elem())
							destVal.Elem().Set(reflect.ValueOf(strings))
							ddFieldValue.Set(destVal)
						} else {
							strings := []string{}
							_ = fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &strings, false)
							if strings == nil {
								strings = make([]string, 0)
							}
							destVal := reflect.New(ddFieldDescriptor.Type.Elem())
							destVal.Elem().Set(reflect.ValueOf(strings))
							ddFieldValue.Set(destVal)
						}
					}
				} else {
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
				}

			default:
				fmt.Printf("WARN: Don't know how to assign anything (type was %s) to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type)
			}
		}
	}
}
