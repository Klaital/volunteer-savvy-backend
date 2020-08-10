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
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

type UserServerTestSuite struct {
	suite.Suite
	Config              *config.ServiceConfig
	Container           *restful.Container
	databaseInitialized bool
}

// TestUserHandlerTestSuite is the "main" entry point for the suite.
func TestUserHandlerTestSuite(t *testing.T) {
	// Initialize the webservice
	cfg := config.ServiceConfig{
		BasePath:            "/vs",
		ServiceVersion:      "0.0.0",
		DatabaseHost:        "localhost",
		DatabaseDriver:      "postgres",
		DatabaseUser:        "vstester",
		DatabasePassword:    "rootpw",
		DatabasePort:        5432,
		DatabaseName:        "vstester",
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
		JwtPrivateKey:       "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS1FJQkFBS0NBZ0VBMEIwWjN6b0Q5UGhvVlZNZ3FPaVhkeit4MjRoNVY3QzJraVlkTzBXellFUUlCZXJsCjR5NjgrS3VveU1UcFJDNENUZDdWNUJ6WWY1S2tLQUxHTVBhL2lMMzRnSUhhV3JmTDUrK0tobUJoRmVQL0lSMXkKUWFmcDZpTXlSY0Q3cFpHb1E2aW90STNwR2k5eHF5VEJ1TDM1NEtTR0VLY1cyMVNqQzkwVU1oRTBnS20zdCtFTwpIWVVkeWdaWVl0eWFTN2pRSURuYUxpZnFCekowcEcrSFFUWVlIajY2VVN2Q1haMXo3bDFrQURvYkphdHJ1eUFOCklxeTRVdlY4NHVocjNFM0hGU0JJOGZGaUlPRFZGTnUxNXJvcVhOaTJKSmF4RjJmNVlZSkkzZkpKRGtWS0dXaEoKYWhVSU1ocU1FTjdYbkMrUVpsdTNER2VKb3FGSlFyeXN6K3JnMzQwRUNMWTZGODJPMEJIbWlCUTJPdVhRcnBTcApuc0EvZjZ2aUJmRkY5SWdmbjU1b2NIZTVwc3hqOThaeWtMa3p6ZE4xTFIxNmh5QWpCbW14TmliNm1sNnBrV1RVCnVhbEM3SVlNWm5IcTBTNnJObUdmeGVXMzhtMDZXSFZSN0RiOTZFb1JKeDRPSnBMZys1WHZtMVc2cHo2TXljcDIKSWZodDJ5SW9yZnRaVXNpVVZ3SUc0c01xUjFleFRNeVV5MGJOQWFwZ1luMXlsMGh0b2g0MEVHV1dyeFpJUCtsRgpOMUd4VFMxMkpDUjVWRzhZcmgxaUNQcmJWSWpRdHdyVS9zTzVvYXI2R2R2WStiMXdTbEJUbUFJbFBpdFJmSkxMCmxXYWpkd2Y5cXRQa200bmZRbnI0MC82ZDBSeGxoSTE2VW1YUEtLOXpaYk1IblZabnFZOENwR1gySTNNQ0F3RUEKQVFLQ0FnQno2TkZoRDdubWRYZitsY2JwN0dsMzVFVFdCYU8zb0ZkKy9MVnBMci9pRE9IL0ViNHFFdnp5N3dDWgptWHBtRzgzeXV2cWNDeWpWbk1ISyt3aVJlc3hnaDVYaFZQRmRkMktjOGtCUDZWd0pTaXZ0c0szVFBZYzlmWTdoClpNT0RpcVdSMFZ6cyt1RHFVYVJZY3FkbWtvQ2FpbWVVM01zUks0bUg4UUR2aGIrZExXbUNnMWxPUGJiQ3IxZ0kKNzk4TDc0b1RlTVU0MGNnNWEvT2xKZXpKK0N2a3BJRTI5azFSd0lFeU1GZWkvTG5qK0I2UFlTR1B2MjBGMzI1SQpIclQ4UldGdE5jY0s0YXNLcnM1ZXlLUCtObytqbUp1UnMxQTh2ZzhhTE9zU21uL3Y3ODErYXlRQWJtaGJKNGY5CldXL3lzRmNYZTF0dkVEZUxHWk5wRlJhVzByV2RLNERKRVo0Nk1NcHJrSWdhVjdmUTJaNEtHcSt6c1ltdVpXRDMKNG44ZWJ4WEdRVERCc0x3aUhyMzdxb3ovTEpzS3graXVub3FITDVKYUkxUWFJSjU2cnA4bVI1blZmTHpYcjZ6dwp1NC82TmxJeEErTVFGdVVudzhqQjl0RVYyN0xodWtBODJsUEQxRlNudm04SDRUT203VFNFaHNPbmlLOERkNTZNCmllZFhVU2pGT3kyNmxXWmJlYlNZS1lFekpLaTJHMXRSTjZ0MU9MRjhIVnhjbmg0azkzQVBQaXNKMWxURzhIVk8KY2NjK3UraktnODNaMUVNQThsQTJuSDJrNFVQZkRRYTVBZ1ZwOStlUjNEMW94Q3lJQ3JpOXE5emY3TnFTN1RtaQptT29nRjhKdE1zV1c4OFI1YmR0R0pPUHRkZ1R3YkFKT3BQOVhmc0UyN1hHbzRHdm5DUUtDQVFFQTdXNktKQXRlCk9reUhQQzZkNkRDdk11OTdsNWtWWVR3RGQrTSs1eWZocWZZNmhCM1QzRUdzV0o5dUV5R1VLUVhxajJEWU9GTGUKeFRyajFadUZjRFFNVlNSNHdtdjhYWFBNV05WOUR2RWFtbWt6N0tOTEp5RzJ4SnpOcWM4T2dmSGZEenhMbEhKSwo3YnpwMU90NDZqaHFQSjZ6eE5nSDdzT0FwQmo1ZkRtTHppclFBbG5lc1o0Z0ZYN3hhckRZcUhVa21wWFRiM0FsCk5VZm43bWo5Z3kxVmxDUWd5eENMb0VqZmtTZHlLQllPZDBGcm9SNlVSWWdlajRMNmNZTXJDV1FqK2F5QUo3SEkKY0k5Q0lpRGhDNEtDUXo3bG1ldjFtdXFxcEtrNXhBMGk2dmVvYTlIRlU1b0gzUzA5eGMwY2luWEZhT0pzbng3aQp4a0ZGclZoV3h3ODFqUUtDQVFFQTRHT2FkREtjbWpEbFFlV21rN2psNnQ2Sk51blJ3MHdQdDcrWFNjNmVRTGlXCldnN2tmOUtaMGkwQ0dzWmpPUUhCaldjQ21OcWI4cEFsU2JSdVNMT05tZ1UyYjhHOWZ0bmYyaVJpMlpxYXlGTEQKTTV6ZzUweHMzM05KblVxMkQ3WmlWZUlBQ2ZPZnVKUTFNVTNRNzNiUnhhT0l5ZnNBSW03dFR3a2gzbms5Y0piVQo5RWQ0RUdRQ1NHNWt2TnBpUHhuUFJjeWQ0MjBDZENyU0NLQmgrTzB5S2wrM3NmcDFZb1FFZEphVHptR2c2VU5iCkhOZVpCZ0NwRURIUzBCQzJDSWQrdlFLZ016VHlOaXBmdEE1VjZkZEk0aHpYT2NYUEdNYVljb2Y4TXZCd2I4LzMKY1ZDcDBYem1wdnJQWTA1SmdJVitIZUEyeDh4eVpMZWtZVTZVREwzOC93S0NBUUVBeTBkdXdqbHhiVnlFRkVTZApNV1F0TytESjRodFFzTFVmQ3cvbWxTWVNFT0FkYld2VUxhbVVrek84bkdpTlh5b1BqcjROb1B0aWUxNVdIbFpPCndxZnRQeUJBdThTVWhyWlQ2R0t2OVpEN2crUTZib25JR0RMSE5rSkIydmJKcHZ1Y1RJRUUvSTEyRldFK21lc1kKMVArRUJXNmkzdzlPaTErYXplUU1CZzJHNHZiSXJKcWhEVlVpaHdUdVVMZ2tadVlVZHIxOER0Ym5KRnp4OTY2dwpEaFZNUmM5QXZGcm9FRTBVREVUSGVnYVlVQVlVemhkT1B0R3h5SkVOTnc1a1ZHQUdaUWNKbWZLZWQ3QlBvTVNoCnFLY09PK0NuMTBhc092eGJLU3N4cCtiUFZIakJHNzYzd1VJSkpaWk1Zd09mUWZSZkZkTjF5QzQ3WGg0WU43ZUEKWkdGaktRS0NBUUJodjJINFBsZnozMXJ2VXVBMnQ3UUlsWHFHbm1MUFJhSVBOSG51SUFEV1J0TFFWbTU1dEQ5bgp6RTEvWm02dzFiQUFMaUIyZjd5eGROT1pnTzBONUpISngzMklQNGlSNnMxV0ZNV3U3MmQvM25YRVZSR3dFSjNZCjFUcjdOeUdLUkxRZm4zek8yUDc2QkM0TDFVOHdFYjJkNy9oVnJHN0prVEwxWEJBUi94U2hxRU1LU3R2bG4vdFMKbkN4c0RHSUNCUGRDKzdqSDZxUElBU05QdUxZVkh4YmNXS2dIOHdnUnExclpnd0xPUTc4NS9pcUJyUFd2SkxpQgpJb01vT2k5aUZjeERBQkxUVzd3cmxsZnNjdFpBTUpWQ1VyZjdKYzFxaUpDK1M3aTBaQU5HNGZ4enMyVFdxaWM4CkZFUWxpV1FCaGFXRDFEbG8zZ256RUFDZWx3RnRiNUgzQW9JQkFRQ3ZUYkVNRkpqc2xFYXpTZmpiamhTZU9TY0cKSWtCZ3dxQVB5UUdxU3oxRGx6WFROZ1VSclJQWUFqRFdjSm9qWlFFWStqNmhxYUhUZVdpU0tFR1ZPS0doSEovSwoweFFIT3hLUkUyVi9QRzZkQXVjRU1QS3pac3lnODJaWkExRHhQeWhHRkhwOW8rdFlFQXZrcHpJM1RsVkNnSi80CmVHN0N0N0YzeU1nZWZ5MFNkdjJuczJlNkxvVkxsVDZzME9JYTZsRHJkMkIwQ0VCVS9mRjJ1eEE1OFUrbG1FZC8KckN6Y0p3eVlFeWlXZzNwM05NOEx5UDh6ZDhzeEc4bm5manF3dk85NHoyclcwTHNBQW5tODEyczR3dm9KU2ErKwowZGlrd1NjckZXZ2ZrTW9WQXF0aFB3djcxbFVmdUw5blZLOXRiYnIxUXBzMmZWYjhQWHlqQmM2RmZUakEKLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K",
		JwtPublicKey:        "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEwQjBaM3pvRDlQaG9WVk1ncU9pWApkeit4MjRoNVY3QzJraVlkTzBXellFUUlCZXJsNHk2OCtLdW95TVRwUkM0Q1RkN1Y1QnpZZjVLa0tBTEdNUGEvCmlMMzRnSUhhV3JmTDUrK0tobUJoRmVQL0lSMXlRYWZwNmlNeVJjRDdwWkdvUTZpb3RJM3BHaTl4cXlUQnVMMzUKNEtTR0VLY1cyMVNqQzkwVU1oRTBnS20zdCtFT0hZVWR5Z1pZWXR5YVM3alFJRG5hTGlmcUJ6SjBwRytIUVRZWQpIajY2VVN2Q1haMXo3bDFrQURvYkphdHJ1eUFOSXF5NFV2Vjg0dWhyM0UzSEZTQkk4ZkZpSU9EVkZOdTE1cm9xClhOaTJKSmF4RjJmNVlZSkkzZkpKRGtWS0dXaEphaFVJTWhxTUVON1huQytRWmx1M0RHZUpvcUZKUXJ5c3orcmcKMzQwRUNMWTZGODJPMEJIbWlCUTJPdVhRcnBTcG5zQS9mNnZpQmZGRjlJZ2ZuNTVvY0hlNXBzeGo5OFp5a0xregp6ZE4xTFIxNmh5QWpCbW14TmliNm1sNnBrV1RVdWFsQzdJWU1abkhxMFM2ck5tR2Z4ZVczOG0wNldIVlI3RGI5CjZFb1JKeDRPSnBMZys1WHZtMVc2cHo2TXljcDJJZmh0MnlJb3JmdFpVc2lVVndJRzRzTXFSMWV4VE15VXkwYk4KQWFwZ1luMXlsMGh0b2g0MEVHV1dyeFpJUCtsRk4xR3hUUzEySkNSNVZHOFlyaDFpQ1ByYlZJalF0d3JVL3NPNQpvYXI2R2R2WStiMXdTbEJUbUFJbFBpdFJmSkxMbFdhamR3ZjlxdFBrbTRuZlFucjQwLzZkMFJ4bGhJMTZVbVhQCktLOXpaYk1IblZabnFZOENwR1gySTNNQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=",
	}
	testSuite := new(UserServerTestSuite)
	testSuite.Config = &cfg
	server := New(&cfg)
	testSuite.Container = restful.NewContainer()
	testSuite.Container.Add(server.GetUsersAPI())
	suite.Run(t, testSuite)
}

func (suite *UserServerTestSuite) SetupAllSuite() {
	// TODO: launch the test database server with docker
	suite.Assert().NotNil(suite.Config.GetDbConn(), "could not connect")
	err := testhelpers.InitializeDatabase(suite.Config.GetDbConn(),
		"file://../../../../db/migrations/",
		"../testdata/")
	if err != nil {
		suite.T().Fatalf("Error initializing the db %v", err)
		return
	}
	suite.databaseInitialized = true
}

// Perform initialization required by each test function
func (suite *UserServerTestSuite) BeforeTest(suiteName, testName string) {
	// For some reason, testify runs the BeforeTest hook before the suite-wide SetupAllSuite,
	// so we must hook on it manually here
	if !suite.databaseInitialized {
		suite.SetupAllSuite()
	}

	// Run some handcrafted SQL to inject common test data from the top-level testdata directory
	dir, err := filepath.Abs(filepath.Join("..", "testdata"))
	if err != nil {
		suite.T().Fatalf("Failed to generate fixture dir path: %v\n", err)
	}
	err = testhelpers.ResetFixtures(suite.Config.GetDbConn(), dir)
	if err != nil {
		suite.T().Fatalf("Failed to load fixtures: %v\n", err)
	}
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
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/vs/users/", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.Container.Dispatch(resp, req)
	suite.Assert().Equal(http.StatusForbidden, resp.Code, "ListUsers API returned incorrect response code")

	//
	// Get a valid JWT with an Org permission
	tokenStr, err = GetUserAuthHeader( "kit@example.org", suite.Config)
	//suite.Assert().NotNilf(err, "Expected no error from loading fixture user. Got %v", err)
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
	expectUsersListed := map[string]bool {
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
