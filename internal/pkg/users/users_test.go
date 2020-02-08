package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/helpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsersTestSuite struct {
	suite.Suite
	DatabaseConnection *sqlx.DB
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

// Perform global setup
func (suite *UsersTestSuite) SetupAllSuite() {
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
func (suite *UsersTestSuite) BeforeTest(suiteName, testName string) {
	if suite.DatabaseConnection == nil {
		suite.SetupAllSuite()
	}
	helpers.CleanupTestDb(suite.DatabaseConnection)
	helpers.LoadFixtures(suite.DatabaseConnection)
}
func (suite *UsersTestSuite) AfterTest(suiteName, testName string) {
	helpers.CleanupTestDb(suite.DatabaseConnection)
}
