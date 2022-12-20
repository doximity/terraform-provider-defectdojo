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
	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	)

	expectedRegulations := []int{99, 100, 101}

	expectedRegulationsSet := types.SetValueMust(
		types.Int64Type,
		[]attr.Value{
			types.Int64Value(99),
			types.Int64Value(100),
			types.Int64Value(101),
		},
	)

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
	assert.Equal(t, productResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, productResource.Description.ValueString(), expectedDescription)
	assert.Equal(t, productResource.Name.ValueString(), expectedName)
	assert.Equal(t, productResource.Revenue.ValueString(), expectedRevenue)
	assert.Equal(t, productResource.BusinessCriticality.ValueString(), expectedBusinessCriticality)
	assert.Equal(t, productResource.Platform.ValueString(), expectedPlatform)
	assert.Equal(t, productResource.Lifecycle.ValueString(), expectedLifecycle)
	assert.Equal(t, productResource.Origin.ValueString(), expectedOrigin)
	assert.Equal(t, productResource.EnableFullRiskAcceptance.ValueBool(), expectedEnableFullRiskAcceptance)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.ValueBool(), expectedEnableSimpleRiskAcceptance)
	assert.Equal(t, productResource.ExternalAudience.ValueBool(), expectedExternalAudience)
	assert.Equal(t, productResource.InternetAccessible.ValueBool(), expectedInternetAccessible)
	assert.Equal(t, productResource.ProductTypeId.ValueInt64(), (int64)(expectedProductTypeId))
	assert.Equal(t, productResource.ProdNumericGrade.ValueInt64(), (int64)(expectedProdNumericGrade))
	assert.Equal(t, productResource.ProductManagerId.ValueInt64(), (int64)(expectedProductManagerId))
	assert.Equal(t, productResource.TeamManagerId.ValueInt64(), (int64)(expectedTeamManagerId))
	assert.Equal(t, productResource.TechnicalContactId.ValueInt64(), (int64)(expectedTechnicalContactId))
	assert.Equal(t, productResource.UserRecords.ValueInt64(), (int64)(expectedUserRecords))
	assert.DeepEqual(t, productResource.Tags, expectedTagsSet)
	assert.DeepEqual(t, productResource.RegulationIds, expectedRegulationsSet)

	ddProduct = productDefectdojoResource{
		Product: dd.Product{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddProduct)

	nilStringSet := types.SetNull(types.StringType)
	nilInt64Set := types.SetNull(types.Int64Type)

	assert.Equal(t, productResource.Description.ValueString(), "")
	assert.Equal(t, productResource.Name.ValueString(), "")
	assert.Equal(t, productResource.Revenue.IsNull(), true)
	assert.Equal(t, productResource.BusinessCriticality.IsNull(), true)
	assert.Equal(t, productResource.Platform.IsNull(), true)
	assert.Equal(t, productResource.Lifecycle.IsNull(), true)
	assert.Equal(t, productResource.Origin.IsNull(), true)
	assert.Equal(t, productResource.EnableFullRiskAcceptance.IsNull(), true)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.IsNull(), true)
	assert.Equal(t, productResource.ExternalAudience.IsNull(), true)
	assert.Equal(t, productResource.InternetAccessible.IsNull(), true)
	assert.Equal(t, productResource.ProductTypeId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.IsNull(), true)
	assert.Equal(t, productResource.ProductManagerId.IsNull(), true)
	assert.Equal(t, productResource.TeamManagerId.IsNull(), true)
	assert.Equal(t, productResource.TechnicalContactId.IsNull(), true)
	assert.Equal(t, productResource.UserRecords.IsNull(), true)
	assert.DeepEqual(t, productResource.Tags, nilStringSet)
	assert.DeepEqual(t, productResource.RegulationIds, nilInt64Set)
}
func TestProductResourcePopulateNils(t *testing.T) {

	nilStringSet := types.SetNull(types.StringType)
	nilInt64Set := types.SetNull(types.Int64Type)

	productResource := productResourceData{}
	var terraformResource terraformResourceData = &productResource

	assert.Equal(t, productResource.Description.ValueString(), "")
	assert.Equal(t, productResource.Name.ValueString(), "")
	assert.Equal(t, productResource.Revenue.ValueString(), "")
	assert.Equal(t, productResource.BusinessCriticality.ValueString(), "")
	assert.Equal(t, productResource.Platform.ValueString(), "")
	assert.Equal(t, productResource.Lifecycle.ValueString(), "")
	assert.Equal(t, productResource.Origin.ValueString(), "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.ValueBool(), false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.ValueBool(), false)
	assert.Equal(t, productResource.ExternalAudience.ValueBool(), false)
	assert.Equal(t, productResource.InternetAccessible.ValueBool(), false)
	assert.Equal(t, productResource.ProductTypeId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.ProductManagerId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.TeamManagerId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.TechnicalContactId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.UserRecords.ValueInt64(), (int64)(0))

	assert.DeepEqual(t, productResource.Tags.Elements(), []attr.Value{})
	assert.DeepEqual(t, productResource.RegulationIds.Elements(), []attr.Value{})

	ddProduct := productDefectdojoResource{
		Product: dd.Product{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddProduct)

	// still all empty/null values after running populate
	assert.Equal(t, productResource.Description.ValueString(), "")
	assert.Equal(t, productResource.Name.ValueString(), "")
	assert.Equal(t, productResource.Revenue.ValueString(), "")
	assert.Equal(t, productResource.BusinessCriticality.ValueString(), "")
	assert.Equal(t, productResource.Platform.ValueString(), "")
	assert.Equal(t, productResource.Lifecycle.ValueString(), "")
	assert.Equal(t, productResource.Origin.ValueString(), "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.ValueBool(), false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.ValueBool(), false)
	assert.Equal(t, productResource.ExternalAudience.ValueBool(), false)
	assert.Equal(t, productResource.InternetAccessible.ValueBool(), false)
	assert.Equal(t, productResource.ProductTypeId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.ProdNumericGrade.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.ProductManagerId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.TeamManagerId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.TechnicalContactId.ValueInt64(), (int64)(0))
	assert.Equal(t, productResource.UserRecords.ValueInt64(), (int64)(0))
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
	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	)

	expectedRegulations := []int{99, 100, 101}

	expectedRegulationsSet := types.SetValueMust(
		types.Int64Type,
		[]attr.Value{
			types.Int64Value(99),
			types.Int64Value(100),
			types.Int64Value(101),
		},
	)

	productResource := productResourceData{
		//Id:                  types.String{Value: fmt.Sprint(expectedId)},
		Name:                types.StringValue(expectedName),
		Description:         types.StringValue(expectedDescription),
		BusinessCriticality: types.StringValue(expectedBusinessCriticality),
		Platform:            types.StringValue(expectedPlatform),
		Lifecycle:           types.StringValue(expectedLifecycle),
		Origin:              types.StringValue(expectedOrigin),
		Revenue:             types.StringValue(expectedRevenue),

		EnableFullRiskAcceptance:   types.BoolValue(expectedEnableFullRiskAcceptance),
		EnableSimpleRiskAcceptance: types.BoolValue(expectedEnableSimpleRiskAcceptance),
		ExternalAudience:           types.BoolValue(expectedExternalAudience),
		InternetAccessible:         types.BoolValue(expectedInternetAccessible),

		ProductTypeId:      types.Int64Value(int64(expectedProductTypeId)),
		ProdNumericGrade:   types.Int64Value(int64(expectedProdNumericGrade)),
		ProductManagerId:   types.Int64Value(int64(expectedProductManagerId)),
		TeamManagerId:      types.Int64Value(int64(expectedTeamManagerId)),
		TechnicalContactId: types.Int64Value(int64(expectedTechnicalContactId)),
		UserRecords:        types.Int64Value(int64(expectedUserRecords)),

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
		Id:                  types.StringNull(),
		Name:                types.StringNull(),
		Description:         types.StringNull(),
		BusinessCriticality: types.StringNull(),
		Platform:            types.StringNull(),
		Lifecycle:           types.StringNull(),
		Origin:              types.StringNull(),
		Revenue:             types.StringNull(),

		EnableFullRiskAcceptance:   types.BoolNull(),
		EnableSimpleRiskAcceptance: types.BoolNull(),
		ExternalAudience:           types.BoolNull(),
		InternetAccessible:         types.BoolNull(),

		ProductTypeId:      types.Int64Null(),
		ProdNumericGrade:   types.Int64Null(),
		ProductManagerId:   types.Int64Null(),
		TeamManagerId:      types.Int64Null(),
		TechnicalContactId: types.Int64Null(),
		UserRecords:        types.Int64Null(),

		Tags:          types.SetNull(types.StringType),
		RegulationIds: types.SetNull(types.Int64Type),
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
