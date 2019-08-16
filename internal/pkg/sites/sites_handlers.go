package sites

import (
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type FindAllSitesResponse struct {
	Sites []Site
}
func FindAllSitesHandler(request *restful.Request, response *restful.Response) {
	sites, err := FindAllSites()
	if err != nil {
		log.WithField("operation", "FindAllSitesHandler").
			Errorf("Failed to fetch sites data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = response.WriteEntity(FindAllSitesResponse{Sites:sites})
	if err != nil {
		log.WithField("operation", "FindAllSitesHandler").
			Errorf("Failed to serialize sites data: %v", err)
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
	var site Site
	jsonErr := request.ReadEntity(&site)
	if jsonErr != nil {
		logger.Errorf("Failed to read JSON from the request: %v", jsonErr)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err := site.CreateSite()
	if err != nil {
		logger.Errorf("Failed to insert site data: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}