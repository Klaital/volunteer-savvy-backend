package organizations

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/helpers"
)

func (suite *OrganizationsTestSuite) TestOrganization_Create() {
	initialOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)

	o := Organization{
		Id:            0,
		Name:          "Test Organization",
		Slug:          "test-organization",
		Authcode:      "authcode1234",
		ContactUserId: 0,
		ContactUser:   nil,
		Latitude:      47.669444,
		Longitude:     -122.123889,
	}

	err := o.Create(suite.DatabaseConnection)
	suite.Nil(err, "Error inserting organization")

	finalOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)
	suite.Equal(initialOrgCount+1, finalOrgCount, "Organization count did not increment")
}

func (suite *OrganizationsTestSuite) TestOrganization_Find() {
	initialOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)

	org, err := FindOrganization(1, suite.DatabaseConnection)
	suite.Nil(err, "Error finding organization")
	suite.NotNil(org, "Org not populated")

	// Validate that all fields are populated

	finalOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)
	suite.Equal(initialOrgCount, finalOrgCount, "Organization count changed")
}