package testhelpers

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestingSuite struct {
	suite.Suite
	Config          *config.ServiceConfig
	suiteConfigured bool
}

// Perform global setup
func (suite *DatabaseTestingSuite) SetupAllSuite() {
	cfg, err := config.GetServiceConfig()
	if err != nil {
		suite.T().Fatalf("Failed to load environment: %v", err)
		return
	}
	if cfg == nil {
		suite.T().Fatalf("Failed to read config")
		return
	}
	err = InitializeDatabase(suite.Config.GetDbConn(), suite.Config.MigrationsPath, suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to init test db: %v", err)
		return
	}
}

// Perform initialization required by each test function
func (suite *DatabaseTestingSuite) BeforeTest(suiteName, testName string) {
	if !suite.suiteConfigured {
		suite.SetupAllSuite()
		suite.suiteConfigured = true
	}
	err := CleanupTestDb(suite.Config.GetDbConn())
	if err != nil {
		suite.T().Fatalf("Failed to cleanup test db: %v", err)
	}
	err = LoadFixtures(suite.Config.GetDbConn(), suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to load fixtures: %v", err)
	}
}

func (suite *DatabaseTestingSuite) AfterTest(suiteName, testName string) {
	CleanupTestDb(suite.Config.GetDbConn())
}
