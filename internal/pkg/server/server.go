package server

import (
	"fmt"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/sites"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"net/http"
	"path"
	"time"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"

	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
)

type Server struct {
	container *restful.Container
	Config    *config.ServiceConfig
}

func New(config *config.ServiceConfig) (*Server, error) {

	server := &Server{
		container: restful.NewContainer(),
		Config:    config,
	}

	// global filters can go here.  route specific filters go in their route definitions
	// JWT handling is an example of a filter that needs to be route specific, since calls to the Swagger API would fail if it were global
	server.container.Filter(JsonLoggingFilter)
	server.container.Filter(server.SetRequestIDFilter)

	// set up the APIs we use for this server setup
	server.setupSupportAPI()

	orgServerConfig := organizations.NewOrganizationsServer(config.GetDbConn())
	orgServerConfig.BasePath = config.BasePath
	orgServer := orgServerConfig.GetOrganizationsAPI()
	server.container.Add(orgServer)

	// Expose Swagger-UI
	//http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("~/bin/swagger-ui/dist"))))

	rootService := new(restful.WebService)

	server.addHealthCheck(rootService)
	server.addSwaggerSupport()
	server.container.Add(rootService)

	return server, nil
}

func (server *Server) Serve() {

	var handler http.Handler = server.container

	innerHandler := handler
	handler = http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		innerHandler.ServeHTTP(respWriter, req)
	})

	ipPort := fmt.Sprintf(":%d", server.Config.Port)

	if server.Config.LogLevel == "debug" {
		log.Infof("Server is launching in debug mode")
	}

	log.WithFields(log.Fields{
		"operation": "Serve",
	}).Infof("Server is listening at %s", ipPort)

	httpServer := &http.Server{
		Addr:         ipPort,
		Handler:      handler,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.WithFields(log.Fields{
			"operation": "Serve",
		}).Fatalf("server failed: %v", err)
	}
}

func (server *Server) addSwaggerSupport() {
	openAPIService := restfulspec.NewOpenAPIService(restfulspec.Config{
		WebServices:     server.container.RegisteredWebServices(),
		APIPath:         server.Config.BasePath + server.Config.APIPath,
	})
	log.Infof("Enabling swagger UI at %s", server.Config.BasePath + server.Config.SwaggerPath)
	http.Handle(
		server.Config.BasePath + server.Config.SwaggerPath,
		http.StripPrefix(
			server.Config.BasePath + server.Config.SwaggerPath,
			http.FileServer(http.Dir("/home/kit/devel/volunteer-savvy-backend/web/swagger-ui/dist"))))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container: server.container,
	}
	server.container.Filter(cors.Filter)
	server.container.Add(openAPIService)
}

func (server *Server) staticFileHandler(req *restful.Request, resp *restful.Response) {
	localFilePath := path.Join(server.Config.StaticContentPath, req.PathParameter("subpath"))
	log.WithField("localPath", localFilePath).Debug("Serving static file")
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		localFilePath)
}
func (server *Server) addHealthCheck(service *restful.WebService) {
	service.Route(service.GET(server.Config.HealthCheckPath).To(server.healthCheckHandler))
}

func (server *Server) healthCheckHandler(request *restful.Request, response *restful.Response) {

	// Ensure that the database connection is useable
	if server.Config.GetDbConn() == nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Error("Database connection not initialized")
		return
	}

	response.WriteHeader(http.StatusOK)
}

// do this for each subpackage/service
func (server *Server) setupSupportAPI() {

	service := new(restful.WebService)
	service.Path(server.Config.BasePath).ApiVersion("0.0.0").Doc("Volunteer-Savvy Backend")

	//
	// Static File Server
	//
	service.Route(
		service.GET("/{subpath:*}").To(server.staticFileHandler))


	//
	// Sites APIs
	//
	service.Route(
		service.GET("/sites/").
			//Filter(filters.RateLimitingFilter).
			To(server.ListSitesHandler).
			Doc("Fetch all sites").
			Produces(restful.MIME_JSON).
			Writes(sites.ListSitesResponse{}).
			Returns(http.StatusOK, "Fetched all sites", sites.ListSitesResponse{}))
	service.Route(
		service.POST("/sites/").
			//Filter(filters.RequireValidJWT).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireAdminPermission).
			To(server.CreateSiteHandler).
			Doc("Fetch all sites").
			Produces(restful.MIME_JSON).
			Consumes(restful.MIME_JSON).
			Reads(sites.Site{}).
			Returns(http.StatusOK, "Created site", nil).
			Returns(http.StatusUnauthorized, "Not logged in", nil).
			Returns(http.StatusForbidden, "Logged-in user is not authorized to create sites", nil))
	service.Route(
		service.GET("/sites/{siteSlug}").
			//Filter(filters.RequireValidJWT).
			//Filter(filters.RateLimitingFilter).
			To(server.DescribeSiteHandler).
			Doc("Fetch all sites").
			Produces(restful.MIME_JSON).
			Writes(sites.Site{}).
			Returns(http.StatusOK, "Fetched site data", sites.Site{}))
	service.Route(
		service.PUT("/sites/{siteSlug}").
			//Filter(filters.RequireValidJWT).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireSiteUpdatePermission).
			To(server.UpdateSiteHandler).
			Doc("Update site config").
			Produces(restful.MIME_JSON).
			Consumes(restful.MIME_JSON).
			Reads(sites.UpdateSiteRequestAdmin{}).
			Writes(sites.Site{}).
			Returns(http.StatusOK, "Site updated", sites.Site{}).
			Returns(http.StatusUnauthorized, "Not logged in", nil).
			Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	service.Route(
		service.DELETE("/sites/{siteSlug}").
			//Filter(filters.RequireValidJWT).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireAdminPermission).
			To(server.DeleteSiteHandler).
			Doc("Delete site and related calendars").
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, "Site deleted", nil).
			Returns(http.StatusUnauthorized, "Not logged in", nil).
			Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//service.Route(
	//	service.PUT("/sites/{siteSlug}/feature/{featureId}").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSiteUpdatePermission).
	//		To(AddSiteFeatureHandler).
	//		Doc("Add a Feature to a Site").
	//		Produces(restful.MIME_JSON).
	//		Reads(AddSiteFeatureRequest{}).
	//		Writes(sites.Site{}).
	//		Returns(http.StatusOK, "Site updated", sites.Site{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//service.Route(
	//	service.DELETE("/sites/{siteSlug}/feature/{featureId}").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSiteUpdatePermission).
	//		To(DeleteSiteFeatureHandler).
	//		Doc("Remove a Feature from a Site").
	//		Produces(restful.MIME_JSON).
	//		Returns(http.StatusOK, "Site feature removed", nil).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//service.Route(
	//	service.PUT("/sites/{siteSlug}/coordinators/{userId}").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSiteUpdatePermission).
	//		To(AddSiteCoordinatorHandler).
	//		Doc("Add a Coordinator to a Site").
	//		Produces(restful.MIME_JSON).
	//		Writes(sites.Site{}).
	//		Returns(http.StatusOK, "Site updated", sites.Site{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//service.Route(
	//	service.DELETE("/sites/{siteSlug}/feature/{featureId}").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSiteUpdatePermission).
	//		To(DeleteSiteFeatureHandler).
	//		Doc("Remove a Coordinator from a Site").
	//		Produces(restful.MIME_JSON).
	//		Returns(http.StatusOK, "Site coordinator removed", sites.Site{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//
	service.Route(
		service.GET("/users/").
			//Filter(filters.RequireValidJWT).
			//Filter(filters.RateLimitingFilter).
			To(users.ListUsersHandler).
			Doc("Fetch all users' details").
			Produces(restful.MIME_JSON).
			Writes([]users.User{}).
			Returns(http.StatusOK, "Got list of users", []users.User{}))
	//service.Route(
	//	service.GET("/users/{userGuid}").
	//		//Filter(filters.RequireValidJWT).
	//		To(FindUserHandler).
	//		Doc("Fetch details on a specific user").
	//		Produces(restful.MIME_JSON).
	//		Writes(users.User{}).
	//		Returns(http.StatusOK, "User data fetched", users.User{}))
	//service.Route(
	//	service.POST("/users/").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireAdminPermission).
	//		To(CreateUserHandler).
	//		Doc("Create a new user account").
	//		Produces(restful.MIME_JSON).
	//		Consumes(restful.MIME_JSON).
	//		Reads(CreateUserRequest{}).
	//		Writes(users.User{}).
	//		Returns(http.StatusOK, "User created", users.User{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to create new users", nil))
	//service.Route(
	//	service.PUT("/users/{userGuid}").
	//		//Filter(filters.RequireValidJWT).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSOAPermission).
	//		To(UpdateUserHandler).
	//		Doc("Update a user's details").
	//		Produces(restful.MIME_JSON).
	//		Reads(UpdateUserRequest{}).
	//		Writes(users.User{}).
	//		Returns(http.StatusOK, "User updated", users.User{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this user", nil))

	server.container.Add(service)
}
