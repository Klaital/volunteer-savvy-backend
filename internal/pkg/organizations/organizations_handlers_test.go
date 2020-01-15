package organizations

import "github.com/sirupsen/logrus"

func (suite *OrganizationsTestSuite) TestListOrganizationsRequest_ListOrganizations() {
	request := ListOrganizationsRequest{
		Db:            suite.DatabaseConnection,
		Organizations: nil,
	}

	err := request.ListOrganizations()
	suite.Nil(err, "Error when fetching organizations list")
	suite.Equal(1, len(request.Organizations), "Incorrect organizations list")
}

func (suite *OrganizationsTestSuite) TestDescribeOrganizationsRequest_DescribeOrganization() {
	request := DescribeOrganizationRequest{
		Db: suite.DatabaseConnection,
		OrganizationId: 1,

		Organization: nil,
	}

	err := request.DescribeOrganization()
	suite.Nilf(err, "Error on DescribeOrganization(): %v", err)
	suite.Equal(uint64(1), request.Organization.Id, "Organization ID not returned")
	suite.Equal(true, len(request.Organization.Name) > 0, "Organization Name not returned")
	suite.Equal(true, len(request.Organization.Slug) > 0, "Organization Slug not returned")

	// Test 2: Validate "Not Found" behavior
	request = DescribeOrganizationRequest{
		Db: suite.DatabaseConnection,
		OrganizationId: 2,

		Organization: nil,
	}

	err = request.DescribeOrganization()
	suite.Nil(err, "Should return a nil error and a nil Organization")
	suite.Nil(request.Organization, "Should return a nil error and a nil Organization")
}

func (suite *OrganizationsTestSuite) TestUpdateOrganizationsRequest_UpdateOrganization() {
	request := UpdateOrganizationRequest{
		Db: suite.DatabaseConnection,
		InputOrganization: &Organization{
			Id:            1,
			Name:          "new organization name",
			Slug:          "new-organization-slug",
			Authcode:      "newauthcode",
			ContactUserId: 0,
			Latitude:      123.45,
			Longitude:     -123.45,
		},

		// Output
		Organization: nil,
	}

	err := request.UpdateOrganization()
	suite.Nilf(err, "Error on UpdateOrganization(): %v", err)
	suite.Equal(uint64(1), request.Organization.Id, "Organization ID not returned")
	suite.Equal("new organization name", request.Organization.Name, "Organization Name not returned")
	suite.Equal("new-organization-slug", request.Organization.Slug, "Organization Slug not returned")
	suite.Equal("newauthcode", request.Organization.Authcode, "Organization authcode not returned")
	suite.Equal(123.45, request.Organization.Latitude, "Organization Latitude not returned")
	suite.Equal(-123.45, request.Organization.Longitude, "Organization Latitude not returned")

	// Reload the organization from the DB to validate that everything was updated
	updatedOrg, err := FindOrganization(1, suite.DatabaseConnection)
	suite.Equal("new organization name", updatedOrg.Name, "Organization Name not returned")
	suite.Equal("new-organization-slug", updatedOrg.Slug, "Organization Slug not returned")
	suite.Equal("newauthcode", updatedOrg.Authcode, "Organization authcode not returned")
	suite.Equal(123.45, updatedOrg.Latitude, "Organization Latitude not returned")
	suite.Equal(-123.45, updatedOrg.Longitude, "Organization Latitude not returned")

	// Test 2: Validate "Not Found" behavior
	logrus.SetLevel(logrus.DebugLevel)
	request = UpdateOrganizationRequest{
		Db: suite.DatabaseConnection,
		InputOrganization: &Organization{
			Id:            200,
			Name:          "sitenotfound",
			Slug:          "sitenotfound",
			Authcode:      "asdf",
			ContactUserId: 0,
			Latitude:      0,
			Longitude:     0,
		},

		Organization: nil,
	}

	err = request.UpdateOrganization()
	suite.Nil(err, "Should return a nil error and a nil Organization")
	suite.Nil(request.Organization, "Should return a nil error and a nil Organization")
}

