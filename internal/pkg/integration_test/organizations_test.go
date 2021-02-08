package integrationtest

import (
	"context"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationsTestSuite struct {
	suite.Suite
	Config *config.ServiceConfig
}

func TestOrganizationsTestSuite(t *testing.T) {

	testSuite := new(OrganizationsTestSuite)
	testSuite.Config = getStaticConfig()

	if testing.Short() {
		t.Skip("Skipping Organizations integration tests in short mode")
	} else {
		suite.Run(t, testSuite)
	}
}

// Perform global setup
func (suite *OrganizationsTestSuite) SetupAllSuite() {
	cfg, err := config.GetServiceConfig()
	if err != nil {
		suite.T().Fatalf("Failed to load environment: %v", err)
		return
	}
	if cfg == nil {
		suite.T().Fatalf("Failed to read config")
		return
	}
	err = testhelpers.InitializeDatabase(suite.Config.GetDbConn(), suite.Config.MigrationsPath, suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to init test db: %v", err)
		return
	}
}

// Perform initialization required by each test function
func (suite *OrganizationsTestSuite) BeforeTest(suiteName, testName string) {
	if suite.Config.GetDbConn() == nil {
		suite.SetupAllSuite()
	}
	testhelpers.CleanupTestDb(suite.Config.GetDbConn())
	testhelpers.LoadFixtures(suite.Config.GetDbConn(), suite.Config.FixturesPath)
}
func (suite *OrganizationsTestSuite) AfterTest(suiteName, testName string) {
	testhelpers.CleanupTestDb(suite.Config.GetDbConn())
}

// TestOrganization_Create tests whether Organizations can be inserted into the database.
func (suite *OrganizationsTestSuite) TestOrganization_Create() {
	if testing.Short() {
		suite.T().Skip("Skipping DB tests in short mode")
		return
	}
	initialOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())

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

	org, err := organizations.DescribeOrganization(context.Background(), suite.Config.GetDbConn(), 1)
	suite.Nil(err, "Error finding organization")
	suite.NotNil(org, "Org not populated")

	// Validate that all fields are populated

	finalOrgCount := testhelpers.CountTable("organizations", suite.Config.GetDbConn())
	suite.Equal(initialOrgCount, finalOrgCount, "Organization count changed")
}
