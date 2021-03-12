package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserServerTestSuite struct {
	suite.Suite
	Config          *config.ServiceConfig
	Container       *restful.Container
	suiteConfigured bool
}

// TestUserHandlerTestSuite is the "main" entry point for the suite.
func TestUserHandlerTestSuite(t *testing.T) {
	// Initialize the webservice
	cfg := testhelpers.GetStandardConfig()
	cfg.FixturesPath = "../testdata/"
	testSuite := new(UserServerTestSuite)
	testSuite.Config = &cfg
	server := New(&cfg)
	testSuite.Container = restful.NewContainer()
	testSuite.Container.Add(server.GetUsersAPI())
	if testing.Short() {
		t.Skip("Skipping User Handlers tests in short mode")
	} else {
		suite.Run(t, testSuite)
	}
}

func (suite *UserServerTestSuite) SetupAllSuite() {
	cfg, err := config.GetServiceConfig()
	if err != nil {
		suite.T().Fatalf("Failed to load environment: %v", err)
		return
	}
	if cfg == nil {
		suite.T().Fatalf("Failed to read config")
		return
	}
	err = testhelpers.InitializeDatabase(suite.Config.GetDbConn(), suite.Config.MigrationsPath, suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to init test db: %v", err)
		return
	}
}

// Perform initialization required by each test function
func (suite *UserServerTestSuite) BeforeTest(suiteName, testName string) {
	if !suite.suiteConfigured {
		suite.SetupAllSuite()
		suite.suiteConfigured = true
	}
	err := testhelpers.CleanupTestDb(suite.Config.GetDbConn())
	if err != nil {
		suite.T().Fatalf("Failed to cleanup test db: %v", err)
	}
	err = testhelpers.LoadFixtures(suite.Config.GetDbConn(), suite.Config.FixturesPath)
	if err != nil {
		suite.T().Fatalf("Failed to load fixtures: %v", err)
	}
}
func (suite *UserServerTestSuite) AfterTest(suiteName, testName string) {
	testhelpers.CleanupTestDb(suite.Config.GetDbConn())
}

func (suite *UserServerTestSuite) TestListUsersHandler() {
	// Pre-declare these to make copy/pasting the test blocks easier
	var resp *httptest.ResponseRecorder
	var req *http.Request
	var err error
	var tokenStr string
	var listUsersResponse ListUsersResponse

	//
	// No JWT included
	//resp = httptest.NewRecorder()
	//req, err = http.NewRequest(http.MethodGet, "/vs/users/", nil)
	//if err != nil {
	//	suite.T().Fatal(err)
	//}
	//suite.Container.Dispatch(resp, req)
	//suite.Assert().Equal(http.StatusForbidden, resp.Code, "ListUsers API returned incorrect response code")

	//
	// Get a valid JWT with an Org permission
	tokenStr, err = GetUserAuthHeader("kit@example.org", suite.Config)
	suite.Assert().NotEqual("", tokenStr, "Expected to get a JWT from helper")
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/vs/users/", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Authorization", tokenStr)
	suite.Container.Dispatch(resp, req)
	suite.Assert().Equal(http.StatusOK, resp.Code, "ListUsers API returned incorrect response code")

	// TODO: check that the users listed are only from the Orgs where the logged-in user is an Admin. If the user is a superadmin, all of them.
	// Kit should be able to see user2 and user3, but not user4
	err = json.Unmarshal(resp.Body.Bytes(), &listUsersResponse)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.Assert().NotEqualf(0, len(listUsersResponse.Users), "Expected some users to be returned. Got %s instead.", resp.Body.String())

	expectUsersListed := map[string]bool{
		"user2": true,
		"user3": true,
		"user4": false,
	}
	usersListed := make(map[string]bool, 3)
	for _, user := range listUsersResponse.Users {
		usersListed[user.Guid] = true
	}
	for userGuid, expectFound := range expectUsersListed {
		if expectFound != usersListed[userGuid] {
			suite.T().Errorf("Expected user %s to be %t. Got %t. Full list: %+v", userGuid, expectFound, usersListed[userGuid], listUsersResponse)
		}
	}
}

func GetUserAuthHeader(email string, config *config.ServiceConfig) (string, error) {
	user, err := users.FindUser(context.Background(), email, config.GetDbConn())
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("no user returned")
	}

	_, err = user.GetRoles(context.Background(), config.GetDbConn())
	if err != nil {
		return "", err
	}

	claims := users.CreateJWT(user, config.GetTokenExpirationDuration())
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	privateKey, _ := config.GetJWTKeys()
	if privateKey == nil {
		return "", errors.New("failed to load private key")
	}
	tokenString, err := token.SignedString(privateKey)
	return fmt.Sprintf("Bearer %s", tokenString), err
}
