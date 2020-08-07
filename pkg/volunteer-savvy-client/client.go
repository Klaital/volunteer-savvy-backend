package volunteersavvyclient

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Client struct {
	Host string
	OauthHost string
	GetTokenPath string

	// Authentication
	Username string
	Password string

	// To enable threaded logging
	LogContext *log.Entry

	// Runtime caching
	httpClient *http.Client
	jwt string
}

func New() *Client {
	return &Client{
		Host:         "",
		Username: "",
		Password: "",
		OauthHost:    "",
		GetTokenPath: "",
		LogContext:   nil,
		httpClient:   nil,
		jwt:          "",
	}
}

func (c *Client) getJwt() string {
	if len(c.jwt) > 0 {
		return c.jwt
	}
	if len(c.Username) == 0 || len(c.Password) == 0{
		return ""
	}

	resp, err := c.GetNewJwt()

}
