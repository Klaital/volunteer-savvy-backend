package users

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsersTestSuite struct {
	testhelpers.DatabaseTestingSuite
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
