package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
	"time"
)

type Server struct {
	container *restful.Container
	Config    *config.ServiceConfig
}

func New(config *config.ServiceConfig, services []*restful.WebService) (*Server, error) {

	server := &Server{
		container: restful.NewContainer(),
		Config:    config,
	}

	// global filters can go here.  route specific filters go in their route definitions
	// JWT handling is an example of a filter that needs to be route specific, since calls to the Swagger API would fail if it were global
	server.container.Filter(filters.JsonLoggingFilter)

	// set up the APIs we use for this server setup
	server.setupSupportAPI()

	for i := range services {
		server.container.Add(services[i])
	}

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
		WebServices: server.container.RegisteredWebServices(),
		APIPath:     server.Config.BasePath + server.Config.APIPath,
	})
	log.Infof("Enabling swagger UI at %s", server.Config.BasePath+server.Config.SwaggerPath)
	http.Handle(
		server.Config.BasePath+server.Config.SwaggerPath,
		http.StripPrefix(
			server.Config.BasePath+server.Config.SwaggerPath,
			http.FileServer(http.Dir("/home/kit/devel/volunteer-savvy-backend/web/swagger-ui/dist"))))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      server.container,
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

	server.container.Add(service)
}
