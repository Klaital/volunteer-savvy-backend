package integrationtest

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IntegrationTestSuite struct {
	testhelpers.DatabaseTestingSuite
}

func TestIntegrationTestSuite(t *testing.T) {

	testSuite := new(IntegrationTestSuite)
	cfg := testhelpers.GetStandardConfig()
	cfg.FixturesPath = "../../../testdata/"
	cfg.MigrationsPath = "file://../../../db/migrations/"
	testSuite.Config = &cfg

	if testing.Short() {
		t.Skip("Skipping Organizations integration tests in short mode")
	} else {
		suite.Run(t, testSuite)
	}
}
