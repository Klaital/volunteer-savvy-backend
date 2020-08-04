package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/auth"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AuthServer struct {
	ApiVersion string
	Config     *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *AuthServer {
	return &AuthServer{
		ApiVersion: version.Version,
		Config:     cfg,
	}
}
func (server *AuthServer) GetAuthAPI() *restful.WebService {

	service := new(restful.WebService)
	service.Path(server.Config.BasePath + "/auth").ApiVersion(server.ApiVersion).Doc("Volunteer-Savvy Backend")

	//
	// Auth APIs
	//
	service.Route(
		service.POST("/token").
			//Filter(filters.RateLimitingFilter).
			To(server.GrantTokenHandler).
			Doc("List Organizations").
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
func (server *AuthServer) GrantTokenHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "GrantTokenHandler",
	})

	email, password, ok := request.Request.BasicAuth()
	logger.WithField("email", email)
	if !ok {
		logger.Error("no basic auth included")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword([]byte(password), server.Config.BcryptCost)
	if err != nil {
		logger.WithError(err).Error("Failed to hash password")
		response.WriteHeader(http.StatusInternalServerError)
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
	if !auth.CheckPassword(hashedPassword, []byte(loggedInUser.PasswordHash)) {
		logger.WithField("PasswordLength", len(password)).Debug("Password mismatch")
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// The user is logged in!
	// Now fetch their roles to include in the JWT
	roles, err := loggedInUser.GetRoles(ctx, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to get user's roles")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add the user GUID to the roles
	claims := jwt.MapClaims{
		"organizations": roles,
		"exp":           time.Now().Add(server.Config.GetTokenExpirationDuration()),
		"iat":           time.Now().Unix(),
		"sub":           loggedInUser.Guid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(server.Config.JwtPrivateKey)
	if err != nil {
		logger.WithError(err).Error("Failed to sign JWT")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseData := AccessTokenResponse{
		AccessToken: tokenString,
		ExpiresIn:   uint(server.Config.GetTokenExpirationDuration().Seconds()),
		Permissions: roles,
	}

	err = response.WriteEntity(responseData)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}