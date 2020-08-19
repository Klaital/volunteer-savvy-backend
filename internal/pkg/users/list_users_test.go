package users

import "context"

func (suite *UsersTestSuite) TestListUsersInSameOrgs() {
	claims := Claims{
		Roles: map[uint64][]RoleType{
			1: []RoleType{
				OrgAdmin,
			},
			2: []RoleType{
				OrgAdmin,
			},
		},
	}

	userSet, err := ListUsersInSameOrgs(context.Background(), &claims, suite.DatabaseConnection)
	suite.Assert().Nilf(err, "Expected no error from ListUsersInSameOrgs. Got %+v", err)
	suite.Assert().Equalf(4, len(userSet), "Expected %d userSet, got %d: %+v", 2, len(userSet), userSet)
}
