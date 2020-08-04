package server

import (
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/auth"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/version"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserServer struct {
	ApiVersion string
	Config     *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *UserServer {
	return &UserServer{
		ApiVersion: version.Version,
		Config:     cfg,
	}
}

func (server *UserServer) GetUsersAPI() *restful.WebService {
	service := new(restful.WebService)
	service.Path(server.Config.BasePath + "/users").ApiVersion(server.ApiVersion)

	service.Route(
		service.GET("/").
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			To(server.ListUsersHandler).
			Doc("Fetch all users' details").
			Produces(restful.MIME_JSON).
			Writes(ListUsersResponse{}).
			Returns(http.StatusOK, "Got list of users", ListUsersResponse{}))
	//service.Route(
	//	service.GET("/{userGuid}").
	//		Filter(filters.ValidJwtFilter).
	//		To(FindUserHandler).
	//		Doc("Fetch details on a specific user").
	//		Produces(restful.MIME_JSON).
	//		Writes(User{}).
	//		Returns(http.StatusOK, "User data fetched", users.User{}))
	//service.Route(
	//	service.POST("/").
	//		Filter(filters.ValidJwtFilter).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireAdminPermission).
	//		To(CreateUserHandler).
	//		Doc("Create a new user account").
	//		Produces(restful.MIME_JSON).
	//		Consumes(restful.MIME_JSON).
	//		Reads(CreateUserRequest{}).
	//		Writes(users.User{}).
	//		Returns(http.StatusOK, "User created", User{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to create new users", nil))
	//service.Route(
	//	service.PUT("/{userGuid}").
	//		Filter(filters.ValidJwtFilter).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSOAPermission).
	//		To(UpdateUserHandler).
	//		Doc("Update a user's details").
	//		Produces(restful.MIME_JSON).
	//		Reads(UpdateUserRequest{}).
	//		Writes(users.User{}).
	//		Returns(http.StatusOK, "User updated", User{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this user", nil))

	return service
}

type ListUsersResponse struct {
	Users []users.User `json:"users"`
}

func (server *UserServer) ListUsersHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := log.WithFields(log.Fields{
		"operation": "ListUsersHandler",
	})

	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fetch users for organizations where the logged-in user is an Admin

	users, err := users.ListUsersInSameOrgs(ctx, request.Attribute("jwt.sub").(string), appConfig.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to list users")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = response.WriteEntity(ListUsersResponse{Users: users})
	if err != nil {
		logger.WithError(err).Error("Failed to serialize users")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}