package provider

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type terraformResourceData interface {
	id() types.String
	defectdojoResource() defectdojoResource
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

var typeOfTypesString = reflect.TypeOf(types.String{})
var typeOfTypesBool = reflect.TypeOf(types.Bool{})
var typeOfTypesInt64 = reflect.TypeOf(types.Int64{})
var typeOfStringSlice = reflect.TypeOf([]string{})
var typeOfInt64Slice = reflect.TypeOf([]int64{})
var typeOfTypesSet = reflect.TypeOf(types.Set{})

func (r terraformResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	data, diags := r.getData(ctx, req.Config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.createApiCall(ctx, r.provider)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 201 {
		populateResourceData(ctx, &diags, &data, ddResource)
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

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.updateApiCall(ctx, r.provider, idNumber)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if statusCode == 200 {
		populateResourceData(ctx, &diags, &data, ddResource)
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

	ddResource := data.defectdojoResource()

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

func populateDefectdojoResource(ctx context.Context, diags *diag.Diagnostics, resourceData terraformResourceData, ddResource *defectdojoResource) {
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
			ddFieldDescriptor, ok := ddVal.Type().FieldByName(ddFieldName)
			if !ok {
				diags.AddError("Error: No such field", fmt.Sprintf("A field named %s was specified to look sync data from the defectdojo client type, but no such field was found.", ddFieldName))
				continue
			}
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
					// the destination field is an int
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					zero := 0
					if !srcIsNull {
						srcVal := fieldValue.FieldByName("Value")
						strVal := srcVal.Interface().(string)
						intVal, err := strconv.Atoi(strVal)

						if err != nil {
							diags.AddError("Error converting value", fmt.Sprintf("Could not convert string value %s to *int: %e", strVal, err))
							continue
						}
						ddFieldValue.Set(reflect.ValueOf(zero))
						ddFieldValue.Set(reflect.ValueOf(intVal))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					// the destination field is a *int
					srcIsNull := fieldValue.FieldByName("Null").Bool()
					if !srcIsNull {
						destType := ddFieldDescriptor.Type.Elem()
						destVal := reflect.New(destType)
						str := fieldValue.FieldByName("Value").String()
						num, err := strconv.Atoi(str)
						if err != nil {
							diags.AddError("Error converting value", fmt.Sprintf("Could not convert string value %s to *int: %e", str, err))
							continue
						}
						destVal.Elem().Set(reflect.ValueOf(num))
						ddFieldValue.Set(destVal)
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
							diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &int64s, false)
							if len(diags_) > 0 {
								diags.Append(diags_...)
								continue
							}
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
							diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &strings, false)
							if len(diags_) > 0 {
								diags.Append(diags_...)
								continue
							}
							if strings == nil {
								strings = make([]string, 0)
							}
							destVal := reflect.New(ddFieldDescriptor.Type.Elem())
							destVal.Elem().Set(reflect.ValueOf(strings))
							ddFieldValue.Set(destVal)
						}
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
				}

			default:
				tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign anything (type was %s) to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
			}
		}
	}
}

func populateResourceData(ctx context.Context, diags *diag.Diagnostics, d *terraformResourceData, ddResource defectdojoResource) {
	tflog.Info(context.Background(), "populateResourceData")

	resourceVal := reflect.ValueOf(*d).Elem()
	resourceType := resourceVal.Type()
	// fmt.Printf("resourceVal: %s\n", resourceVal)
	// fmt.Printf("resourceType: %s\n", resourceType)

	ddVal := reflect.ValueOf(ddResource).Elem()

	for i := 0; i < resourceVal.NumField(); i++ {
		fieldDescriptor := resourceType.Field(i)
		// fmt.Printf("field: %s\n", fieldDescriptor.Name)
		tag := fieldDescriptor.Tag
		ddFieldName := tag.Get("ddField")
		if ddFieldName != "" {
			fieldValue := resourceVal.Field(i)

			ddFieldDescriptor, ok := ddVal.Type().FieldByName(ddFieldName)
			if !ok {
				diags.AddError("Error: No such field", fmt.Sprintf("A field named %s was specified to look sync data from the defectdojo client type, but no such field was found.", ddFieldName))
				continue
			}
			ddFieldValue := ddVal.FieldByName(ddFieldName)

			// fmt.Printf("ddFieldDescriptor: Kind = %s, Name = %s\n", ddFieldDescriptor.Type.Kind(), ddFieldDescriptor.Name)
			// fmt.Printf("fieldDescriptor: Kind = %s, Name = %s, type = %s\n", fieldDescriptor.Type.Kind(), fieldDescriptor.Name, fieldDescriptor.Type)

			switch fieldDescriptor.Type {

			case typeOfTypesString:
				if ddFieldDescriptor.Type.Kind() == reflect.String {
					// if the source field is a string, we can use it directly
					fieldValue.Set(reflect.ValueOf(types.String{Value: ddFieldValue.String()}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
					// if the source field is a pointer, make sure it's a pointer to a string, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.String{Value: ddFieldValue.Elem().String()}))
					} else {
						fieldValue.Set(reflect.ValueOf(types.String{Null: true}))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int {
					fieldValue.Set(reflect.ValueOf(types.String{Value: fmt.Sprint(ddFieldValue.Int())}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.String{Value: fmt.Sprint(ddFieldValue.Elem().Int())}))
					} else {
						fieldValue.Set(reflect.ValueOf(types.String{Null: true}))
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}

			case typeOfTypesBool:
				if ddFieldDescriptor.Type.Kind() == reflect.Bool {
					// if the source field is a bool, we can use it directly
					fieldValue.Set(reflect.ValueOf(types.Bool{Value: ddFieldValue.Bool()}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Bool {
					// if the source field is a pointer, make sure it's a pointer to a bool, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Bool{Value: ddFieldValue.Elem().Bool()}))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Bool{Null: true}))
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}

			case typeOfTypesInt64:
				if ddFieldDescriptor.Type.Kind() == reflect.Int64 || ddFieldDescriptor.Type.Kind() == reflect.Int {
					// if the source field is an int or int64, we can cast and use it directly
					fieldValue.Set(reflect.ValueOf(types.Int64{Value: (int64)(ddFieldValue.Int())}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && (ddFieldDescriptor.Type.Elem().Kind() == reflect.Int64 || ddFieldDescriptor.Type.Elem().Kind() == reflect.Int) {
					// if the source field is a pointer, make sure it's a pointer to an int64, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Int64{Value: (int64)(ddFieldValue.Elem().Int())}))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Int64{Null: true}))
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}

			case typeOfTypesSet:
				if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the source field is a pointer to a slice
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.FieldByName("Null").Bool()) {
							elems := []attr.Value{}
							for _, val := range ddFieldValue.Elem().Interface().([]int) {
								elems = append(elems, types.Int64{Value: (int64)(val)})
							}
							destVal := types.Set{
								ElemType: types.Int64Type,
								Elems:    elems,
							}
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.Set{
								ElemType: types.Int64Type,
								Null:     true,
							}
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.FieldByName("Null").Bool()) {
							elems := []attr.Value{}
							for _, val := range ddFieldValue.Elem().Interface().([]string) {
								elems = append(elems, types.String{Value: (string)(val)})
							}
							destVal := types.Set{
								ElemType: types.StringType,
								Elems:    elems,
							}
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.Set{
								ElemType: types.StringType,
								Null:     true,
							}
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}
			default:
				tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign anything (type was %s) to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
			}
		}
	}
}
