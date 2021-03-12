package server

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthServerTestSuite struct {
	testhelpers.DatabaseTestingSuite
	Container *restful.Container
}

// TestAuthHandlerTestSuite is the "main" entry point for the suite.
func TestAuthHandlerTestSuite(t *testing.T) {
	// Initialize the webservice
	cfg := testhelpers.GetStandardConfig()
	cfg.FixturesPath = "../testdata/"
	testSuite := new(AuthServerTestSuite)
	testSuite.Config = &cfg
	server := New(&cfg)
	testSuite.Container = restful.NewContainer()
	testSuite.Container.Add(server.GetAuthAPI())

	if testing.Short() {
		t.Skip("Skipping Auth Handler tests in short mode")
	}
}

func (suite *AuthServerTestSuite) TestAuthServer_GrantTokenHandler() {

	// Pre-declare these to make copy/pasting the test blocks easier
	var resp *httptest.ResponseRecorder
	var req *http.Request
	var err error

	//
	// No Basic Auth included in request
	//
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	if suite.Container == nil {
		suite.T().Fatal("Container not initialized")
	}
	suite.Container.Dispatch(resp, req)
	suite.Assert().Equal(http.StatusUnauthorized, resp.Code, "GrantTokenHandler returned wrong status code")

	//
	// Incorrect basic auth included in request
	//
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.SetBasicAuth("kit@example.org", "wrongpassword")
	suite.Container.Dispatch(resp, req)
	suite.Assert().Equal(http.StatusUnauthorized, resp.Code, "GrantTokenHandler returned wrong status code")

	//
	// Correct basic auth included in request
	//
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.SetBasicAuth("kit@example.org", "password")
	suite.Container.Dispatch(resp, req)
	suite.Assert().Equal(http.StatusOK, resp.Code, "GrantTokenHandler returned wrong status code")

	// Check that kit has roles claimed
	var respData AccessTokenResponse
	err = json.Unmarshal(resp.Body.Bytes(), &respData)
	suite.Assert().Nilf(err, "Expected no err from unmarshaling token response. Instead got %+v", err)
	suite.Assert().Equal(1, len(respData.Permissions), "Expected 1 org permissions in the auth response")
	suite.Assert().Equal(1, len(respData.Permissions[1]), "expected 1 role on org 1")

	// Decode the JWT itself and validate
	suite.T().Log(respData.AccessToken)
	claims := users.DecodeJWT(respData.AccessToken, suite.Config.GetPublicKey())
	suite.Assert().NotNil(claims, "Expected valid claims")
	// TODO: extract the JWT and use it to make a follow-on request
}
