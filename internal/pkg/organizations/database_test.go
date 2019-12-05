package organizations

//import (
//	"github.com/klaital/volunteer-savvy-backend/internal/pkg/helpers"
//	"github.com/sirupsen/logrus"
//	"runtime/debug"
//)
//
//func (suite *OrganizationsTestSuite) TestOrganization_Create() {
//	defer func() {
//		if r := recover(); r == nil {
//			return
//		}
//
//		logrus.Errorf("panichandler trace: %s\n", string(debug.Stack()))
//	}()
//	initialOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)
//
//	o := Organization{
//		Id:            0,
//		Name:          "Test Organization",
//		Slug:          "test-organization",
//		Authcode:      "authcode1234",
//		ContactUserId: 0,
//		ContactUser:   nil,
//		Latitude:      47.669444,
//		Longitude:     -122.123889,
//	}
//
//	err := o.Create(suite.DatabaseConnection)
//	suite.Assertions.NotNil(err, "Error inserting organization")
//
//	finalOrgCount := helpers.CountTable("organizations", suite.DatabaseConnection)
//	suite.Assertions.Equal(initialOrgCount+1, finalOrgCount, "Organization count did not increment")
//}
