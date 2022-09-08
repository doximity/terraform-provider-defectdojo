package provider

import (
	"context"
	"fmt"
	"testing"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestProductResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedDescription := "A Description"
	expectedName := "A Name"
	expectedBusinessCriticality := "medium"
	expectedPlatform := "web"
	expectedLifecycle := "construction"
	expectedOrigin := "internal"
	expectedRevenue := "1,000,000.00"

	expectedEnableFullRiskAcceptance := true
	expectedEnableSimpleRiskAcceptance := true
	expectedExternalAudience := true
	expectedInternetAccessible := true

	expectedProductTypeId := 42
	expectedProdNumericGrade := 43
	expectedProductManagerId := 44
	expectedTeamManagerId := 45
	expectedTechnicalContactId := 46
	expectedUserRecords := 47

	expectedTags := []string{"foo", "bar", "baz"}
	expectedTagsSet := types.Set{
		ElemType: types.StringType,
		Elems: []attr.Value{
			types.String{Value: "foo"},
			types.String{Value: "bar"},
			types.String{Value: "baz"},
		},
	}

	expectedRegulations := []int{99, 100, 101}

	expectedRegulationsSet := types.Set{
		ElemType: types.Int64Type,
		Elems: []attr.Value{
			types.Int64{Value: 99},
			types.Int64{Value: 100},
			types.Int64{Value: 101},
		},
	}

	ddProduct := productDefectdojoResource{
		Product: dd.Product{
			Id:                         expectedId,
			Description:                expectedDescription,
			Name:                       expectedName,
			Revenue:                    &expectedRevenue,
			BusinessCriticality:        (*dd.ProductBusinessCriticality)(&expectedBusinessCriticality),
			Platform:                   (*dd.ProductPlatform)(&expectedPlatform),
			Lifecycle:                  (*dd.ProductLifecycle)(&expectedLifecycle),
			Origin:                     (*dd.ProductOrigin)(&expectedOrigin),
			EnableFullRiskAcceptance:   &expectedEnableFullRiskAcceptance,
			EnableSimpleRiskAcceptance: &expectedEnableSimpleRiskAcceptance,
			ExternalAudience:           &expectedExternalAudience,
			InternetAccessible:         &expectedInternetAccessible,
			ProdType:                   expectedProductTypeId,
			ProdNumericGrade:           &expectedProdNumericGrade,
			ProductManager:             &expectedProductManagerId,
			TeamManager:                &expectedTeamManagerId,
			TechnicalContact:           &expectedTechnicalContactId,
			UserRecords:                &expectedUserRecords,
			Tags:                       &expectedTags,
			Regulations:                &expectedRegulations,
		},
	}

	productResource := productResourceData{}
	var terraformResource terraformResourceData = &productResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddProduct)
	assert.Equal(t, productResource.Id.Value, fmt.Sprint(expectedId))
	assert.Equal(t, productResource.Description.Value, expectedDescription)
	assert.Equal(t, productResource.Name.Value, expectedName)
	assert.Equal(t, productResource.Revenue.Value, expectedRevenue)
	assert.Equal(t, productResource.BusinessCriticality.Value, expectedBusinessCriticality)
	assert.Equal(t, productResource.Platform.Value, expectedPlatform)
	assert.Equal(t, productResource.Lifecycle.Value, expectedLifecycle)
	assert.Equal(t, productResource.Origin.Value, expectedOrigin)
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Value, expectedEnableFullRiskAcceptance)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Value, expectedEnableSimpleRiskAcceptance)
	assert.Equal(t, productResource.ExternalAudience.Value, expectedExternalAudience)
	assert.Equal(t, productResource.InternetAccessible.Value, expectedInternetAccessible)
	assert.Equal(t, productResource.ProductTypeId.Value, (int64)(expectedProductTypeId))
	assert.Equal(t, productResource.ProdNumericGrade.Value, (int64)(expectedProdNumericGrade))
	assert.Equal(t, productResource.ProductManagerId.Value, (int64)(expectedProductManagerId))
	assert.Equal(t, productResource.TeamManagerId.Value, (int64)(expectedTeamManagerId))
	assert.Equal(t, productResource.TechnicalContactId.Value, (int64)(expectedTechnicalContactId))
	assert.Equal(t, productResource.UserRecords.Value, (int64)(expectedUserRecords))
	assert.DeepEqual(t, productResource.Tags, expectedTagsSet)
	assert.DeepEqual(t, productResource.RegulationIds, expectedRegulationsSet)

	ddProduct = productDefectdojoResource{
		Product: dd.Product{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddProduct)

	nilStringSet := types.Set{Null: true, ElemType: types.StringType}
	nilInt64Set := types.Set{Null: true, ElemType: types.Int64Type}

	assert.Equal(t, productResource.Description.Value, "")
	assert.Equal(t, productResource.Name.Value, "")
	assert.Equal(t, productResource.Revenue.Null, true)
	assert.Equal(t, productResource.BusinessCriticality.Null, true)
	assert.Equal(t, productResource.Platform.Null, true)
	assert.Equal(t, productResource.Lifecycle.Null, true)
	assert.Equal(t, productResource.Origin.Null, true)
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Null, true)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Null, true)
	assert.Equal(t, productResource.ExternalAudience.Null, true)
	assert.Equal(t, productResource.InternetAccessible.Null, true)
	assert.Equal(t, productResource.ProductTypeId.Value, (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.Null, true)
	assert.Equal(t, productResource.ProductManagerId.Null, true)
	assert.Equal(t, productResource.TeamManagerId.Null, true)
	assert.Equal(t, productResource.TechnicalContactId.Null, true)
	assert.Equal(t, productResource.UserRecords.Null, true)
	assert.DeepEqual(t, productResource.Tags, nilStringSet)
	assert.DeepEqual(t, productResource.RegulationIds, nilInt64Set)
}
func TestProductResourcePopulateNils(t *testing.T) {

	nilStringSet := types.Set{Null: true, ElemType: types.StringType}
	nilInt64Set := types.Set{Null: true, ElemType: types.Int64Type}

	emptyStringSet := types.Set{Elems: []attr.Value{}}
	emptyInt64Set := types.Set{Elems: []attr.Value{}}

	productResource := productResourceData{}
	var terraformResource terraformResourceData = &productResource

	assert.Equal(t, productResource.Description.Value, "")
	assert.Equal(t, productResource.Name.Value, "")
	assert.Equal(t, productResource.Revenue.Value, "")
	assert.Equal(t, productResource.BusinessCriticality.Value, "")
	assert.Equal(t, productResource.Platform.Value, "")
	assert.Equal(t, productResource.Lifecycle.Value, "")
	assert.Equal(t, productResource.Origin.Value, "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Value, false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Value, false)
	assert.Equal(t, productResource.ExternalAudience.Value, false)
	assert.Equal(t, productResource.InternetAccessible.Value, false)
	assert.Equal(t, productResource.ProductTypeId.Value, (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.Value, (int64)(0))
	assert.Equal(t, productResource.ProductManagerId.Value, (int64)(0))
	assert.Equal(t, productResource.TeamManagerId.Value, (int64)(0))
	assert.Equal(t, productResource.TechnicalContactId.Value, (int64)(0))
	assert.Equal(t, productResource.UserRecords.Value, (int64)(0))
	assert.DeepEqual(t, productResource.Tags, emptyStringSet)
	assert.DeepEqual(t, productResource.RegulationIds, emptyInt64Set)

	ddProduct := productDefectdojoResource{
		Product: dd.Product{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddProduct)

	// still all empty/null values after running populate
	assert.Equal(t, productResource.Description.Value, "")
	assert.Equal(t, productResource.Name.Value, "")
	assert.Equal(t, productResource.Revenue.Value, "")
	assert.Equal(t, productResource.BusinessCriticality.Value, "")
	assert.Equal(t, productResource.Platform.Value, "")
	assert.Equal(t, productResource.Lifecycle.Value, "")
	assert.Equal(t, productResource.Origin.Value, "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Value, false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Value, false)
	assert.Equal(t, productResource.ExternalAudience.Value, false)
	assert.Equal(t, productResource.InternetAccessible.Value, false)
	assert.Equal(t, productResource.ProductTypeId.Value, (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.Value, (int64)(0))
	assert.Equal(t, productResource.ProductManagerId.Value, (int64)(0))
	assert.Equal(t, productResource.TeamManagerId.Value, (int64)(0))
	assert.Equal(t, productResource.TechnicalContactId.Value, (int64)(0))
	assert.Equal(t, productResource.UserRecords.Value, (int64)(0))
	assert.DeepEqual(t, productResource.Tags, nilStringSet)
	assert.DeepEqual(t, productResource.RegulationIds, nilInt64Set)
}

func TestProductResource__defectdojoResource(t *testing.T) {
	//expectedId := 99
	expectedDescription := "A Description"
	expectedName := "A Name"
	expectedBusinessCriticality := "medium"
	expectedPlatform := "web"
	expectedLifecycle := "construction"
	expectedOrigin := "internal"
	expectedRevenue := "1,000,000.00"

	expectedEnableFullRiskAcceptance := true
	expectedEnableSimpleRiskAcceptance := true
	expectedExternalAudience := true
	expectedInternetAccessible := true

	expectedProductTypeId := 42
	expectedProdNumericGrade := 43
	expectedProductManagerId := 44
	expectedTeamManagerId := 45
	expectedTechnicalContactId := 46
	expectedUserRecords := 47

	expectedTags := []string{"foo", "bar", "baz"}
	expectedTagsSet := types.Set{
		ElemType: types.StringType,
		Elems: []attr.Value{
			types.String{Value: "foo"},
			types.String{Value: "bar"},
			types.String{Value: "baz"},
		},
	}

	expectedRegulations := []int{99, 100, 101}

	expectedRegulationsSet := types.Set{
		ElemType: types.Int64Type,
		Elems: []attr.Value{
			types.Int64{Value: 99},
			types.Int64{Value: 100},
			types.Int64{Value: 101},
		},
	}

	productResource := productResourceData{
		//Id:                  types.String{Value: fmt.Sprint(expectedId)},
		Name:                types.String{Value: expectedName},
		Description:         types.String{Value: expectedDescription},
		BusinessCriticality: types.String{Value: expectedBusinessCriticality},
		Platform:            types.String{Value: expectedPlatform},
		Lifecycle:           types.String{Value: expectedLifecycle},
		Origin:              types.String{Value: expectedOrigin},
		Revenue:             types.String{Value: expectedRevenue},

		EnableFullRiskAcceptance:   types.Bool{Value: expectedEnableFullRiskAcceptance},
		EnableSimpleRiskAcceptance: types.Bool{Value: expectedEnableSimpleRiskAcceptance},
		ExternalAudience:           types.Bool{Value: expectedExternalAudience},
		InternetAccessible:         types.Bool{Value: expectedInternetAccessible},

		ProductTypeId:      types.Int64{Value: int64(expectedProductTypeId)},
		ProdNumericGrade:   types.Int64{Value: int64(expectedProdNumericGrade)},
		ProductManagerId:   types.Int64{Value: int64(expectedProductManagerId)},
		TeamManagerId:      types.Int64{Value: int64(expectedTeamManagerId)},
		TechnicalContactId: types.Int64{Value: int64(expectedTechnicalContactId)},
		UserRecords:        types.Int64{Value: int64(expectedUserRecords)},

		Tags:          expectedTagsSet,
		RegulationIds: expectedRegulationsSet,
	}

	ddResource := productResource.defectdojoResource()
	ddProduct := ddResource.(*productDefectdojoResource)
	var terraformResource terraformResourceData = &productResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	//assert.Equal(t, ddProduct.Id, expectedId)
	assert.Equal(t, ddProduct.Name, expectedName)
	assert.Equal(t, ddProduct.Description, expectedDescription)

	assert.Equal(t, (string)(*ddProduct.BusinessCriticality), expectedBusinessCriticality)
	assert.Equal(t, (string)(*ddProduct.Platform), expectedPlatform)
	assert.Equal(t, (string)(*ddProduct.Lifecycle), expectedLifecycle)
	assert.Equal(t, (string)(*ddProduct.Origin), expectedOrigin)
	assert.Equal(t, (string)(*ddProduct.Revenue), expectedRevenue)

	assert.Equal(t, *ddProduct.EnableFullRiskAcceptance, expectedEnableFullRiskAcceptance)
	assert.Equal(t, *ddProduct.EnableSimpleRiskAcceptance, expectedEnableSimpleRiskAcceptance)
	assert.Equal(t, *ddProduct.ExternalAudience, expectedExternalAudience)
	assert.Equal(t, *ddProduct.InternetAccessible, expectedInternetAccessible)

	assert.Equal(t, ddProduct.ProdType, expectedProductTypeId)
	assert.Equal(t, *ddProduct.ProdNumericGrade, expectedProdNumericGrade)
	assert.Equal(t, *ddProduct.ProductManager, expectedProductManagerId)
	assert.Equal(t, *ddProduct.TeamManager, expectedTeamManagerId)
	assert.Equal(t, *ddProduct.TechnicalContact, expectedTechnicalContactId)
	assert.Equal(t, *ddProduct.UserRecords, expectedUserRecords)

	assert.DeepEqual(t, *ddProduct.Tags, expectedTags)
	assert.DeepEqual(t, *ddProduct.Regulations, expectedRegulations)
}

func TestProductResource__defectdojoResource_Nulls(t *testing.T) {
	var nilBusinessCriticality *dd.ProductBusinessCriticality
	var nilPlatform *dd.ProductPlatform
	var nilLifecycle *dd.ProductLifecycle
	var nilOrigin *dd.ProductOrigin
	var nilString *string
	var nilBool *bool
	var nilInt *int
	// var nilInts *[]int
	// var nilStrings *[]string

	productResource := productResourceData{
		Id:                  types.String{Null: true},
		Name:                types.String{Null: true},
		Description:         types.String{Null: true},
		BusinessCriticality: types.String{Null: true},
		Platform:            types.String{Null: true},
		Lifecycle:           types.String{Null: true},
		Origin:              types.String{Null: true},
		Revenue:             types.String{Null: true},

		EnableFullRiskAcceptance:   types.Bool{Null: true},
		EnableSimpleRiskAcceptance: types.Bool{Null: true},
		ExternalAudience:           types.Bool{Null: true},
		InternetAccessible:         types.Bool{Null: true},

		ProductTypeId:      types.Int64{Null: true},
		ProdNumericGrade:   types.Int64{Null: true},
		ProductManagerId:   types.Int64{Null: true},
		TeamManagerId:      types.Int64{Null: true},
		TechnicalContactId: types.Int64{Null: true},
		UserRecords:        types.Int64{Null: true},

		Tags:          types.Set{Null: true},
		RegulationIds: types.Set{Null: true},
	}

	ddResource := productResource.defectdojoResource()
	ddProduct := ddResource.(*productDefectdojoResource)
	var terraformResource terraformResourceData = &productResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddProduct.Id, 0)
	assert.Equal(t, ddProduct.Name, "")
	assert.Equal(t, ddProduct.Description, "")

	assert.Equal(t, ddProduct.BusinessCriticality, nilBusinessCriticality)
	assert.Equal(t, ddProduct.Platform, nilPlatform)
	assert.Equal(t, ddProduct.Lifecycle, nilLifecycle)
	assert.Equal(t, ddProduct.Origin, nilOrigin)
	assert.Equal(t, ddProduct.Revenue, nilString)

	assert.Equal(t, ddProduct.EnableFullRiskAcceptance, nilBool)
	assert.Equal(t, ddProduct.EnableSimpleRiskAcceptance, nilBool)
	assert.Equal(t, ddProduct.ExternalAudience, nilBool)
	assert.Equal(t, ddProduct.InternetAccessible, nilBool)

	assert.Equal(t, ddProduct.ProdType, 0)
	assert.Equal(t, ddProduct.ProdNumericGrade, nilInt)
	assert.Equal(t, ddProduct.ProductManager, nilInt)
	assert.Equal(t, ddProduct.TeamManager, nilInt)
	assert.Equal(t, ddProduct.TechnicalContact, nilInt)
	assert.Equal(t, ddProduct.UserRecords, nilInt)

	assert.DeepEqual(t, *ddProduct.Tags, []string{})
	assert.DeepEqual(t, *ddProduct.Regulations, []int{})
}
