package sites

import (
	"database/sql"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"net/http"
)

type SitesServer struct {
	Config *config.ServiceConfig
}

type ListSitesResponse struct {
	Sites []Site `json:"sites"`
}
func (server *SitesServer) ListSitesHandler(request *restful.Request, response *restful.Response) {
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx)

	// Fetch sites list.
	// TODO: add optional search filters
	// TODO: add filter for only showing sites for the org that the user is logged-in to
	siteSet, err := ListSites(ctx, server.Config.GetDbConn())
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
	requestSite := Site{}
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
	s, err := DescribeSite(ctx, requestedSiteSlug, server.Config.GetDbConn())
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
	requestSite := Site{}
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
	updateRequest := UpdateSiteRequestAdmin{Site: requestSite}
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
	err := DeleteSite(ctx, server.Config.GetDbConn(), requestedSiteSlug)

	if err != nil {
		// TODO: check the error type to discern between "slug not found" and "db error"
		logger.WithError(err).Error("Failed to fetch site data")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success!
	response.WriteHeader(http.StatusOK)
}
