package organizations

import (
	"context"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"testing"
)

// TestOrganization_Create tests whether Organizations can be inserted into the database.
func (suite *OrganizationsTestSuite) TestOrganization_Create() {
	if testing.Short() {
		suite.T().Skip("Skipping DB tests in short mode")
		return
	}
	initialOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())

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

	err := o.Create(context.Background(), suite.Config.GetDbConn())
	suite.Nil(err, "Error inserting organization")

	finalOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Equal(initialOrgCount+1, finalOrgCount, "Organization count did not increment")
}

// TestOrganization_Find tests whether Organizations can be selected from the DB.
func (suite *OrganizationsTestSuite) TestOrganization_Find() {
	if testing.Short() {
		suite.T().Skip("Skipping DB tests in short mode")
		return
	}
	initialOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())

	org, err := DescribeOrganization(context.Background(), suite.Config.GetDbConn(), 1)
	suite.Nil(err, "Error finding organization")
	suite.NotNil(org, "Org not populated")

	// Validate that all fields are populated

	finalOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Equal(initialOrgCount, finalOrgCount, "Organization count changed")
}
