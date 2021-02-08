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

	userSet, err := ListUsersInSameOrgs(context.Background(), &claims, suite.Config.GetDbConn())
	suite.Assert().Nilf(err, "Expected no error from ListUsersInSameOrgs. Got %+v", err)
	suite.Assert().Equalf(3, len(userSet), "Expected %d userSet, got %d: %+v", 3, len(userSet), userSet)
	expectedUsers := map[string]bool{
		"kit":   true,
		"user2": true,
		"user3": true,
		"user4": false,
	}
	for userGuid, expectedInList := range expectedUsers {
		suite.Assert().Equalf(expectedInList, userInSet(userGuid, userSet), "UserSet incorrect. Expected user in list? %t, got %t", expectedInList, userInSet(userGuid, userSet))
	}
}
func userInSet(guid string, set []User) bool {
	for _, u := range set {
		if u.Guid == guid {
			return true
		}
	}
	return false
}
