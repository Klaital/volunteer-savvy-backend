package integrationtest

import (
	"context"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"testing"
)

// TestOrganization_Create tests whether Organizations can be inserted into the database.
func (suite *IntegrationTestSuite) TestOrganization_Create() {
	if testing.Short() {
		suite.T().Skip("Skipping DB tests in short mode")
		return
	}
	initialOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Less(0, initialOrgCount, "There should be some starting data in the table from fixtures")
	o := organizations.Organization{
		Id:            0,
		Name:          "Test Organization",
		Slug:          "test-organization",
		Authcode:      "authcode1234",
		ContactUserId: 0,
		ContactUser:   nil,
		Latitude:      47.669444,
		Longitude:     -122.123889,
	}
	suite.T().Logf("Org before %+v", o)

	err := o.Create(context.Background(), suite.Config.GetDbConn())
	suite.Nil(err, "Error inserting organization")
	suite.NotEqual(0, o.Id, "Org ID not updated")

	// Validate that the Table was incremented
	finalOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Equal(initialOrgCount+1, finalOrgCount, "Organization count did not increment")

	// Fetch the org back from the DB and assert equality
	o2, err := organizations.DescribeOrganization(context.Background(), suite.Config.GetDbConn(), int64(o.Id))
	suite.Nil(err, "Error re-fetching organization")
	suite.Equal(o.Id, o2.Id, "Org ID mismatch")
}

// TestOrganization_Find tests whether Organizations can be selected from the DB.
func (suite *IntegrationTestSuite) TestOrganization_FindSlug() {
	if testing.Short() {
		suite.T().Skip("Skipping DB tests in short mode")
		return
	}
	initialOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())

	org, err := organizations.DescribeOrganizationBySlug(context.Background(), suite.Config.GetDbConn(), "test-organization-1")
	suite.Nil(err, "Error finding organization from fixture")
	suite.NotNil(org, "Org not populated")

	// Validate that all fields are populated

	finalOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Equal(initialOrgCount, finalOrgCount, "Organization count changed")
}
