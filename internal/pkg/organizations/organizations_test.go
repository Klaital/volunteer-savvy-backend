package organizations

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationsTestSuite struct {
	suite.Suite
	Config *config.ServiceConfig
}

func TestOrganizationsTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OrganizationsTestSuite in short mode")
	} else {
		suite.Run(t, new(OrganizationsTestSuite))
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
