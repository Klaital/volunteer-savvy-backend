package organizations

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationsTestSuite struct {
	suite.Suite
}

func TestOrganizationsTestSuite(t *testing.T) {
	testSuite := new(OrganizationsTestSuite)
	suite.Run(t, testSuite)
}
