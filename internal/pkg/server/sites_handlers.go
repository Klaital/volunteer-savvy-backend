package server

import (
	"context"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/sites"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ListSitesResponse struct {
	Sites []sites.Site `json:"sites"`
}
func (server *Server) ListSitesHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Fetch sites list.
	// TODO: add optional search filters
	// TODO: add filter for only showing sites for the org that the user is logged-in to
	siteSet, err := sites.ListSites(ctx)
	if err != nil {
		// TODO: inspect err type to discern between "DB error" and "no results found"
		ctx.Logger.WithError(err).Error("Failed to fetch sites list")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response payload
	responseData := ListSitesResponse{Sites: siteSet}
	err = response.WriteEntity(responseData)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize response")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) CreateSiteHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Deserialize the request body
	requestSite := sites.Site{}
	err := request.ReadEntity(&requestSite)
	if err != nil {
		ctx.Logger.WithError(err).Error("Unable to deserialize the request body")
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	// Save it
	err = requestSite.Create(ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to save site")
		// TODO: check err type to discern between 500 and 400 responses (such as a duplicate slug, etc)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the updated site data back down
	err = response.WriteEntity(requestSite)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) DescribeSiteHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Read input - the requested site slug
	requestedSiteSlug := request.PathParameter("siteSlug")
	s, err := sites.DescribeSite(ctx, requestedSiteSlug)

	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to fetch site data")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the response payload
	err = response.WriteEntity(s)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize site response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) UpdateSiteHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Get the Slug for the site to be updated
	slug := request.PathParameter("siteSlug")
	if len(slug) == 0 {
		ctx.Logger.WithError(errors.New("no slug given")).Warn("This shouldn't happen - an empty site slug was passed to Update Site handler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Deserialize the request body
	requestSite := sites.Site{}
	err := request.ReadEntity(&requestSite)
	if err != nil {
		ctx.Logger.WithError(err).Error("Unable to deserialize the request body")
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	// TODO: Check the logged-in user's permissions for this site to determine what fields to save.

	// Save it
	updateRequest := sites.UpdateSiteRequestAdmin{Site: requestSite}
	err = requestSite.UpdateSiteAdmin(ctx, &updateRequest)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to save site")
		// TODO: check err type to discern between 500 and 400 responses (such as a duplicate slug, etc)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the updated site data back down
	err = response.WriteEntity(requestSite)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize response body")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) DeleteSiteHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Read input - the requested site slug
	requestedSiteSlug := request.PathParameter("siteSlug")
	err := sites.DeleteSite(ctx, requestedSiteSlug)

	if err != nil {
		// TODO: check the error type to discern between "slug not found" and "db error"
		ctx.Logger.WithError(err).Error("Failed to fetch site data")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success!
	response.WriteHeader(http.StatusOK)
}
