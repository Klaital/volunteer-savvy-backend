package sites

import (
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type FindAllSitesResponse struct {
	Sites []Site
}
func FindAllSitesHandler(request *restful.Request, response *restful.Response) {
	organizationId := request.PathParameter("organizationId")
	logger := log.WithFields(log.Fields{
		"operation": "FindAllSitesHandler",
		"org": organizationId,
	})

	organizationIdInt, err := strconv.Atoi(organizationId)
	if err != nil {
		logger.Errorf("Invalid organization ID. %v", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	sites, err := FindAllSites(organizationIdInt)
	if err != nil {
		logger.Errorf("Failed to fetch sites data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = response.WriteEntity(FindAllSitesResponse{Sites:sites})
	if err != nil {
		logger.Errorf("Failed to serialize sites data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func FindSiteHandler(request *restful.Request, response *restful.Response) {
	siteSlug := request.PathParameter("siteSlug")
	logger := log.WithFields(log.Fields{
		"operation": "FindSiteHandler",
		"SiteSlug": siteSlug,
	})
	site, err := FindSite(siteSlug)
	if err != nil {
		logger.Errorf("Failed to fetch site data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = response.WriteEntity(site)
	if err != nil {
		logger.Errorf("Failed to serialize site data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CreateSitesRequest struct {
	Site
}
func CreateSiteHandler(request *restful.Request, response *restful.Response) {
	logger := log.WithField("operation", "CreateSiteHandler")
	var site CreateSitesRequest
	jsonErr := request.ReadEntity(&site)
	if jsonErr != nil {
		logger.Errorf("Failed to read JSON from the request: %v", jsonErr)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err := site.Site.CreateSite()
	if err != nil {
		if strings.Index(err.Error(), "duplicate key value violates unique constraint") > 0 {
			logger.Errorf("Requested site has a duplicate slug: %v", err)
			response.WriteHeader(http.StatusBadRequest)
		} else {
			logger.Errorf("Failed to insert site data: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	response.WriteHeader(http.StatusOK)
}

func DeleteSiteHandler(request *restful.Request, response *restful.Response) {
	siteSlug := request.PathParameter("siteSlug")
	logger := log.WithFields(log.Fields{
		"operation": "DeleteSiteHandler",
		"SiteSlug": siteSlug,
	})
	err := DeleteSite(siteSlug)
	if err != nil {
		// TODO: discern between a DB error and an attempt to delete a site that doesn't exist
		logger.Errorf("Failed to delete site: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
