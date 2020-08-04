package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/testhelpers"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthServerTestSuite struct {
	suite.Suite
	Config *config.ServiceConfig
	Container *restful.Container
}

// TestAuthHandlerTestSuite is the "main" entry point for the suite.
func TestAuthHandlerTestSuite(t *testing.T) {
	testSuite := new(AuthServerTestSuite)
	suite.Run(t, testSuite)
}

func (suite *AuthServerTestSuite) SetupAllSuite() {
	// TODO: launch the test database server with docker
	err := testhelpers.InitializeDatabase(suite.Config.GetDbConn(), "file://internal/pkg/auth/migrations/", "file://internal/pkg/auth/testdata/")
	if err != nil {
		suite.T().Fatalf("Error initializing the db %v", err)
		return
	}

	// Initialize the webservice
	cfg := config.ServiceConfig{
		BasePath:            "/vs",
		ServiceVersion:      "0.0.0",
		DatabaseHost:        "localhost",
		DatabaseDriver:      "postgres",
		DatabaseUser:        "vstest",
		DatabasePassword:    "rootpw",
		DatabasePort:        5432,
		DatabaseName:        "vstest",
		LogLevel:            "debug",
		LogStyle:            "prettyjson",
		Port:                8080,
		StaticContentPath:   "/web",
		SwaggerFilePath:     "/apidocs.json",
		APIPath:             "",
		SwaggerPath:         "/apidocs",
		HealthCheckPath:     "/GetServiceStatus",
		Logger:              log.NewEntry(log.New()),
		BcryptCost:          1,
		TokenExpirationTime: "4h",
		JwtPrivateKey:       "", // TODO:
		JwtPublicKey:        "", // TODO:
	}
	testSuite := new(AuthServerTestSuite)
	testSuite.Config = &cfg
	server := New(&cfg)
	suite.Container = restful.NewContainer()
	suite.Container.Add(server.GetAuthAPI())
}

// Perform initialization required by each test function
func (suite *AuthServerTestSuite) BeforeTest(suiteName, testName string) {
	// For some reason, testify runs the BeforeTest hook before the suite-wide SetupAllSuite,
	// so we must hook on it manually here
	if suite.Config.GetDbConn() == nil {
		suite.SetupAllSuite()
	}

	// Run some handcrafted SQL to inject common test data from the top-level testdata directory
	err := testhelpers.LoadFixtures(suite.Config.GetDbConn(), "file://internal/pkg/auth/testdata/")
	if err != nil {
		suite.T().Fatalf("Failed to load fixtures: %v", err)
	}
}

func (suite *AuthServerTestSuite) TestAuthServer_GrantTokenHandler(t *testing.T) {

	// Pre-declare these to make copy/pasting the test blocks easier
	var resp *httptest.ResponseRecorder
	var req *http.Request
	var err error

	//
	// No Basic Auth included in request
	//
	fmt.Println("Starting test 1")
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	fmt.Println("Request/response created")
	if err != nil {
		suite.T().Fatal(err)
	}
	fmt.Println("Dispatching request")
	suite.Container.Dispatch(resp, req)

	fmt.Println("Testing response")
	if resp.Code != http.StatusUnauthorized {
		suite.T().Errorf("GrantTokenHandler returned wrong status code: Expected %v want %v",
			resp.Code, http.StatusUnauthorized)
	}

	//
	// Incorrect basic auth included in request
	//
	//resp = httptest.NewRecorder()
	//req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//suite.Container.Dispatch(resp, req)
	//
	//if resp.Code != http.StatusUnauthorized {
	//	t.Errorf("GrantTokenHandler returned wrong status code: Expected %v want %v",
	//		resp.Code, http.StatusUnauthorized)
	//}

}
