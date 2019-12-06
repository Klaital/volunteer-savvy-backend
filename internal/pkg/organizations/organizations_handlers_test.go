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
