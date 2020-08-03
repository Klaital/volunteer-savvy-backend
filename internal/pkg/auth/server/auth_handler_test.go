package server

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthServer_GrantTokenHandler(t *testing.T) {

	server := New(&config.ServiceConfig{
		BasePath:            "/vs",
		ServiceVersion:      "0.0.0",
		DatabaseHost:        "localhost",
		DatabaseDriver:      "postgres",
		DatabaseUser:        "vstester",
		DatabasePassword:    "rootpw",
		DatabasePort:        5432,
		DatabaseName:        "vstest",
		LogLevel:            "debug",
		LogStyle:            "prettyjson",
		Port:                8080,
		HealthCheckPath:     "/healthz",
		Logger:              log.NewEntry(log.New()),
		BcryptCost:          0,
		TokenExpirationTime: "4h",
		JwtPrivateKey:       "", // TODO
		JwtPublicKey:        "", // TODO
	})
	handler := http.HandlerFunc(server.GrantTokenHandler)

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
		t.Fatal(err)
	}
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("GrantTokenHandler returned wrong status code: Expected %v want %v",
			resp.Code, http.StatusUnauthorized)
	}

	//
	// Incorrect basic auth included in request
	//
	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/vs/auth/token", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("GrantTokenHandler returned wrong status code: Expected %v want %v",
			resp.Code, http.StatusUnauthorized)
	}

}
