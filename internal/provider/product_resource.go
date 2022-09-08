package provider

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type productResourceType struct{}

func (t productResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the Product",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the Product",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A[^\s].*[^\s]\z`), "The description must not have leading or trailing whitespace"),
				},
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
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("very high", "high", "medium", "low", "very low", "none", ""),
				},
			},
			"platform": {
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("web service", "desktop", "iot", "mobile", "web", ""),
				},
			},
			"lifecycle": {
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("construction", "production", "retirement", ""),
				},
			},
			"origin": {
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("third party library", "purchased", "contractor", "internal", "open source", "outsourced", ""),
				},
			},
			"user_records": {
				MarkdownDescription: "Estimate the number of user records within the application.",
				Optional:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					int64validator.AtLeast(0),
				},
			},
			"revenue": {
				MarkdownDescription: "Estimate the application's revenue.",
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A-?\d{0,13}(?:\.\d{0,2})?\z`), `Must be a decimal number format, i.e. /^-?\d{0,13}(?:\.\d{0,2})?$/`),
				},
			},
			"external_audience": {
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"internet_accessible": {
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"enable_skip_risk_acceptance": {
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"enable_full_risk_acceptance": {
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Optional:            true,
				Type:                types.BoolType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					boolDefault(false),
				},
			},
			"product_manager_id": {
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"technical_contact_id": {
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"team_manager_id": {
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"regulation_ids": {
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"product_type_id": {
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
				Type:                types.Int64Type,
			},
			"tags": {
				MarkdownDescription: "Tags to apply to the product",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Validators: []tfsdk.AttributeValidator{
					setvalidator.ValuesAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`\A[a-z]+\z`), "Tags must be lower case values"),
					),
				},
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType, // the id (for import purposes) MUST be a string
			},
		},
	}, nil
}

func (t productResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return productResource{
		terraformResource: terraformResource{
			provider:     provider,
			dataProvider: productDataProvider{},
		},
	}, diags
}

type productResourceData struct {
	Name                       types.String `tfsdk:"name" ddField:"Name"`
	Description                types.String `tfsdk:"description" ddField:"Description"`
	ProductTypeId              types.Int64  `tfsdk:"product_type_id" ddField:"ProdType"`
	Id                         types.String `tfsdk:"id" ddField:"Id"`
	BusinessCriticality        types.String `tfsdk:"business_criticality" ddField:"BusinessCriticality"`
	EnableFullRiskAcceptance   types.Bool   `tfsdk:"enable_full_risk_acceptance" ddField:"EnableFullRiskAcceptance"`
	EnableSimpleRiskAcceptance types.Bool   `tfsdk:"enable_skip_risk_acceptance" ddField:"EnableSimpleRiskAcceptance"`
	ExternalAudience           types.Bool   `tfsdk:"external_audience" ddField:"ExternalAudience"`
	InternetAccessible         types.Bool   `tfsdk:"internet_accessible" ddField:"InternetAccessible"`
	Lifecycle                  types.String `tfsdk:"lifecycle" ddField:"Lifecycle"`
	Origin                     types.String `tfsdk:"origin" ddField:"Origin"`
	Platform                   types.String `tfsdk:"platform" ddField:"Platform"`
	ProdNumericGrade           types.Int64  `tfsdk:"prod_numeric_grade" ddField:"ProdNumericGrade"`
	ProductManagerId           types.Int64  `tfsdk:"product_manager_id" ddField:"ProductManager"`
	RegulationIds              types.Set    `tfsdk:"regulation_ids" ddField:"Regulations"`
	Revenue                    types.String `tfsdk:"revenue" ddField:"Revenue"`
	Tags                       types.Set    `tfsdk:"tags" ddField:"Tags"`
	TeamManagerId              types.Int64  `tfsdk:"team_manager_id" ddField:"TeamManager"`
	TechnicalContactId         types.Int64  `tfsdk:"technical_contact_id" ddField:"TechnicalContact"`
	UserRecords                types.Int64  `tfsdk:"user_records" ddField:"UserRecords"`
}

type productDefectdojoResource struct {
	dd.Product
}

func (ddr *productDefectdojoResource) createApiCall(ctx context.Context, p provider) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := dd.ProductsCreateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Product = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) readApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := p.client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) updateApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.ProductsUpdateJSONRequestBody(ddr.Product)
	apiResp, err := p.client.ProductsUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productDefectdojoResource) deleteApiCall(ctx context.Context, p provider, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := p.client.ProductsDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type productResource struct {
	terraformResource
}

type productDataProvider struct{}

func (r productDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productResourceData) id() types.String {
	return d.Id
}

var typeOfTypesString = reflect.TypeOf(types.String{})
var typeOfTypesBool = reflect.TypeOf(types.Bool{})
var typeOfTypesInt64 = reflect.TypeOf(types.Int64{})
var typeOfStringSlice = reflect.TypeOf([]string{})
var typeOfInt64Slice = reflect.TypeOf([]int64{})
var typeOfTypesSet = reflect.TypeOf(types.Set{})

func (d *productResourceData) populate(ddResource defectdojoResource) {
	tflog.Info(context.Background(), "populate")
	product := ddResource.(*productDefectdojoResource)

	resourceVal := reflect.ValueOf(d).Elem()
	resourceType := resourceVal.Type()
	// fmt.Printf("resourceVal: %s\n", resourceVal)
	// fmt.Printf("resourceType: %s\n", resourceType)

	ddVal := reflect.ValueOf(product).Elem()

	for i := 0; i < resourceVal.NumField(); i++ {
		fieldDescriptor := resourceType.Field(i)
		// fmt.Printf("field: %s\n", fieldDescriptor.Name)
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
				} else {
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
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
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
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
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
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
					fmt.Printf("WARN: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
				}
			default:
				fmt.Printf("WARN: Don't know how to assign anything (type was %s) to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type)
			}
		}
	}
}

func (d *productResourceData) defectdojoResource(diags *diag.Diagnostics) (defectdojoResource, error) {
	tflog.Info(context.Background(), "defectdojoResource")

	product := dd.Product{}

	resourceVal := reflect.ValueOf(d).Elem()
	resourceType := resourceVal.Type()
	// fmt.Printf("resourceVal: %s\n", resourceVal)
	// fmt.Printf("resourceType: %s\n", resourceType)

	ddVal := reflect.ValueOf(&product).Elem()

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

	return &productDefectdojoResource{
		Product: product,
	}, nil
}
