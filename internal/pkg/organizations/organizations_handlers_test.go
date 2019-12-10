package organizations


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
