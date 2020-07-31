package server

import (
	"database/sql"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/auth"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/sites"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/version"
	"net/http"
)

type SitesServer struct {
	ApiVersion string
	Config *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *SitesServer {
	return &SitesServer{
		ApiVersion: version.Version,
		Config: cfg,
	}
}

func (server *SitesServer) GetSitesAPI() *restful.WebService {
	service := new(restful.WebService)
	service.Path(server.Config.BasePath + "/sites").ApiVersion(server.ApiVersion)

	//
	// Sites APIs
	//
	service.Route(
		service.GET("/sites/").
			//Filter(filters.RateLimitingFilter).
			To(server.ListSitesHandler).
			Doc("Fetch all sites").
			Produces(restful.MIME_JSON).
			Writes(ListSitesResponse{}).
			Returns(http.StatusOK, "Fetched all sites", ListSitesResponse{}))
	service.Route(
		service.POST("/sites/").
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			// TODO: how do I write a filter that can inspect the Site object for the Organization ID, then validate that that user has permissions on that org+site?
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
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			To(server.DescribeSiteHandler).
			Doc("Fetch all sites").
			Produces(restful.MIME_JSON).
			Writes(sites.Site{}).
			Returns(http.StatusOK, "Fetched site data", sites.Site{}))
	service.Route(
		service.PUT("/sites/{siteSlug}").
			Filter(auth.ValidJwtFilter).
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
			Filter(auth.ValidJwtFilter).
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
	//		Filter(filters.ValidJwtFilter).
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
	//		Filter(filters.ValidJwtFilter).
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
	//		Filter(filters.ValidJwtFilter).
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
	//		Filter(filters.ValidJwtFilter).
	//		//Filter(filters.RateLimitingFilter).
	//		//Filter(filters.RequireSiteUpdatePermission).
	//		To(DeleteSiteFeatureHandler).
	//		Doc("Remove a Coordinator from a Site").
	//		Produces(restful.MIME_JSON).
	//		Returns(http.StatusOK, "Site coordinator removed", sites.Site{}).
	//		Returns(http.StatusUnauthorized, "Not logged in", nil).
	//		Returns(http.StatusForbidden, "Logged-in user is not authorized to update this site", nil))
	//

	return service
}

type ListSitesResponse struct {
	Sites []sites.Site `json:"sites"`
}
func (server *SitesServer) ListSitesHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Fetch sites list.
	// TODO: add optional search filters
	// TODO: add filter for only showing sites for the org that the user is logged-in to
	siteSet, err := sites.ListSites(ctx, server.Config.GetDbConn())
	if err != nil {
		// TODO: inspect err type to discern between "DB error" and "no results found"
		logger.WithError(err).Error("Failed to fetch sites list")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response payload
	responseData := ListSitesResponse{Sites: siteSet}
	err = response.WriteEntity(responseData)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize response")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *SitesServer) CreateSiteHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Deserialize the request body
	requestSite := sites.Site{}
	err := request.ReadEntity(&requestSite)
	if err != nil {
		logger.WithError(err).Error("Unable to deserialize the request body")
		err = response.WriteError(http.StatusBadRequest, err)
		if err != nil {
			logger.WithError(err).Error("Failed to write error response")
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Save it
	err = requestSite.Create(ctx, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to save site")
		// TODO: check err type to discern between 500 and 400 responses (such as a duplicate slug, etc)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the updated site data back down
	err = response.WriteEntity(requestSite)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *SitesServer) DescribeSiteHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Read input - the requested site slug
	requestedSiteSlug := request.PathParameter("siteSlug")
	s, err := sites.DescribeSite(ctx, requestedSiteSlug, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to fetch site data")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the response payload
	err = response.WriteEntity(s)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize site response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *SitesServer) UpdateSiteHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Get the Slug for the site to be updated
	slug := request.PathParameter("siteSlug")
	if len(slug) == 0 {
		logger.WithError(errors.New("no slug given")).Warn("This shouldn't happen - an empty site slug was passed to Update Site handler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Deserialize the request body
	requestSite := sites.Site{}
	err := request.ReadEntity(&requestSite)
	if err != nil {
		logger.WithError(err).Error("Unable to deserialize the request body")
		err = response.WriteError(http.StatusBadRequest, err)
		if err != nil {
			logger.WithError(err).Error("Error writing error")
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// TODO: Check the logged-in user's permissions for this site to determine what fields to save.

	// Save it
	updateRequest := sites.UpdateSiteRequestAdmin{Site: requestSite}
	err = requestSite.UpdateSiteAdmin(ctx, server.Config.GetDbConn(), &updateRequest)
	if err != nil {
		logger.WithError(err).Error("Failed to save site")
		// TODO: check err type to discern between 500 and 400 responses (such as a duplicate slug, etc)
		if err == sql.ErrNoRows {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the updated site data back down
	err = response.WriteEntity(requestSite)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *SitesServer) DeleteSiteHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Read input - the requested site slug
	requestedSiteSlug := request.PathParameter("siteSlug")
	err := sites.DeleteSite(ctx, server.Config.GetDbConn(), requestedSiteSlug)

	if err != nil {
		// TODO: check the error type to discern between "slug not found" and "db error"
		logger.WithError(err).Error("Failed to fetch site data")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success!
	response.WriteHeader(http.StatusOK)
}
