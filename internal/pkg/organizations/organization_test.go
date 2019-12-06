package organizations

import "fmt"

func (suite *OrganizationsTestSuite) TestNew() {
	var o *Organization
	o = New()
	suite.NotNil(o, "Expected New() to return a valid pointer to an Organization")
	suite.Equal(uint64(0), o.Id, "Incorrectly initialized Org ID")
}

func (suite *OrganizationsTestSuite) TestOrganization_Validate() {
	o := New()

	isValid, _ := o.Validate()
	suite.Equal(false, isValid, "Expected New() to return a stub that does not pass validation")

	// Try to make a valid organization
	validOrg := &Organization{
		Id:            0,
		Name:          "testorg",
		Slug:          "testorg",
		Authcode:      "supersecret",
		ContactUserId: 5,
		ContactUser:   nil,
		Latitude:      0,
		Longitude:     0,
	}

	isValid, errs := validOrg.Validate()
	suite.Equal(true, isValid, "Expected Org to pass validation")
	suite.Equal(0, len(errs), "Expected no errors when validation passes")

	// Set an invalid slug
	o = &Organization{
		Id:            0,
		Name:          "testorg",
		Slug:          "TestOrg&!",
		Authcode:      "supersecret",
		ContactUserId: 5,
		ContactUser:   nil,
		Latitude:      0,
		Longitude:     0,
	}

	isValid, _ = o.Validate()
	suite.Equal(false, isValid, fmt.Sprintf("Expected slug '%s' to be invalid, but successfully validated", o.Slug))
}
