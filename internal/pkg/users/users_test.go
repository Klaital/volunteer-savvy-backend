package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsersTestSuite struct {
	suite.Suite
	Config             *config.ServiceConfig
	databaseConnection *sqlx.DB
	suiteConfigured    bool
}

func TestUsersTestSuite(t *testing.T) {
	cfg := testhelpers.GetStandardConfig()
	cfg.FixturesPath = "./testdata/"
	cfg.MigrationsPath = "file://../../../db/migrations/"
	testSuite := new(UsersTestSuite)
	testSuite.Config = &cfg
	if testing.Short() {
		t.Skip("Skipping UsersTestSuite in short mode")
	} else {
		suite.Run(t, testSuite)
	}
}

// Perform global setup
func (suite *UsersTestSuite) SetupAllSuite() {
	suite.Assert().NotNil(suite.Config.GetDbConn(), "could not connect")
	err := testhelpers.InitializeDatabase(suite.Config.GetDbConn(),
		suite.Config.MigrationsPath,
		suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Error initializing the db %v", err)
		return
	}
	suite.databaseConnection = suite.Config.GetDbConn()
}

// Perform initialization required by each test function
func (suite *UsersTestSuite) BeforeTest(suiteName, testName string) {
	if !suite.suiteConfigured {
		suite.SetupAllSuite()
		suite.suiteConfigured = true
	}
	err := testhelpers.CleanupTestDb(suite.Config.GetDbConn())
	if err != nil {
		suite.T().Fatalf("Failed to cleanup test db: %v", err)
	}
	err = testhelpers.LoadFixtures(suite.Config.GetDbConn(), suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to load fixtures: %v", err)
	}
}

func (suite *UsersTestSuite) AfterTest(suiteName, testName string) {
	testhelpers.CleanupTestDb(suite.Config.GetDbConn())
}
