package organizations

import (
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/helpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationsTestSuite struct {
	suite.Suite
	DatabaseConnection *sqlx.DB
}

func TestOrganizationsTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationsTestSuite))
}

// Perform global setup
func (suite *OrganizationsTestSuite) SetupAllSuite() {
	db, err := helpers.InitializeTestDb()
	if err != nil {
		suite.T().Fatalf("Failed to init test db: %v", err)
		return
	}
	if db == nil {
		suite.T().Fatal("Nil database handle returned")
		return
	}
	suite.DatabaseConnection = db
}

// Perform initialization required by each test function
func (suite *OrganizationsTestSuite) BeforeTest(suiteName, testName string) {
	if suite.DatabaseConnection == nil {
		suite.SetupAllSuite()
	}
	helpers.LoadFixtures(suite.DatabaseConnection)
}
func (suite *OrganizationsTestSuite) AfterTest(suiteName, testName string) {
	helpers.CleanupTestDb(suite.DatabaseConnection)
}

