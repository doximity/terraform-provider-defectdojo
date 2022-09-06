package provider

import (
	"fmt"
	"testing"

	dd "github.com/doximity/defect-dojo-client-go"
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

	ddProduct := productDefectdojoResource{
		Product: dd.Product{
			Id:                         expectedId,
			Description:                expectedDescription,
			Name:                       expectedName,
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
		},
	}
	productResource := productResourceData{}
	productResource.populate(&ddProduct)
	assert.Equal(t, productResource.Id.Value, fmt.Sprint(expectedId))
	assert.Equal(t, productResource.Description.Value, expectedDescription)
	assert.Equal(t, productResource.Name.Value, expectedName)
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
}
func TestProductResourcePopulateNils(t *testing.T) {

	productResource := productResourceData{}
	assert.Equal(t, productResource.Description.Value, "")
	assert.Equal(t, productResource.Name.Value, "")
	assert.Equal(t, productResource.BusinessCriticality.Value, "")
	assert.Equal(t, productResource.Platform.Value, "")
	assert.Equal(t, productResource.Lifecycle.Value, "")
	assert.Equal(t, productResource.Origin.Value, "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Value, false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Value, false)
	assert.Equal(t, productResource.ExternalAudience.Value, false)
	assert.Equal(t, productResource.InternetAccessible.Value, false)

	ddProduct := productDefectdojoResource{
		Product: dd.Product{},
	}
	productResource.populate(&ddProduct)

	// still all empty/null values after running populate
	assert.Equal(t, productResource.Description.Value, "")
	assert.Equal(t, productResource.Name.Value, "")
	assert.Equal(t, productResource.BusinessCriticality.Value, "")
	assert.Equal(t, productResource.Platform.Value, "")
	assert.Equal(t, productResource.Lifecycle.Value, "")
	assert.Equal(t, productResource.Origin.Value, "")
	assert.Equal(t, productResource.EnableFullRiskAcceptance.Value, false)
	assert.Equal(t, productResource.EnableSimpleRiskAcceptance.Value, false)
	assert.Equal(t, productResource.ExternalAudience.Value, false)
	assert.Equal(t, productResource.InternetAccessible.Value, false)
}
