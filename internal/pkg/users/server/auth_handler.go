package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	log "github.com/sirupsen/logrus"
	"net/http"
)


func (server *UserServer) GetAuthAPI() *restful.WebService {

	service := new(restful.WebService)
	service.Path(server.Config.BasePath + "/auth").ApiVersion(server.ApiVersion).Doc("Volunteer-Savvy Backend")

	//
	// Auth APIs
	//
	service.Route(
		service.POST("/token").
			//Filter(filters.RateLimitingFilter).
			To(server.GrantTokenHandler).
			Doc("User login. Returns a signed JWT.").
			Produces(restful.MIME_JSON).
			Writes(AccessTokenResponse{}).
			Returns(http.StatusOK, "Successfully logged in.", AccessTokenResponse{}).
			Returns(http.StatusUnauthorized, "Email/password combination did not match.", nil))

	return service
}

type AccessTokenResponse struct {
	AccessToken string                  `json:"access_token"`
	ExpiresIn   uint                    `json:"expires_in"`
	Permissions map[uint64][]users.Role `json:"permissions"`
}

// GrantTokenHandler allows users to log in with their email/password, and get
// back a signed JWT in response
// TODO: support client_credentials grant type if we split into microservices
func (server *UserServer) GrantTokenHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "GrantTokenHandler",
	})

	email, password, ok := request.Request.BasicAuth()
	logger.WithField("email", email)
	if !ok {
		logger.Debug("no basic auth included")
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	loggedInUser, err := users.GetUserForLogin(ctx, email, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Error fetching user")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if loggedInUser == nil {
		logger.WithField("PasswordLength", len(password)).Debug("Email not found")
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if !users.CheckPassword([]byte(loggedInUser.PasswordHash), []byte(password)) {
		logger.WithField("PasswordLength", len(password)).Debug("Password mismatch")
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// The user is logged in!
	// Now fetch their roles to include in the JWT
	_, err = loggedInUser.GetRoles(ctx, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to get user's roles")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WithField("LoggingInUser", fmt.Sprintf("%+v", loggedInUser)).Debug("Added user's roles")

	// Add the user GUID to the roles
	claims := users.CreateJWT(loggedInUser, server.Config.GetTokenExpirationDuration())
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	privateKey, _ := server.Config.GetJWTKeys()
	if privateKey == nil {
		logger.Error("Failed to load JWT Keys")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		logger.WithError(err).Error("Failed to sign JWT")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseData := AccessTokenResponse{
		AccessToken: tokenString,
		ExpiresIn:   uint(server.Config.GetTokenExpirationDuration().Seconds()),
		Permissions: loggedInUser.Roles,
	}

	err = response.WriteEntity(responseData)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
