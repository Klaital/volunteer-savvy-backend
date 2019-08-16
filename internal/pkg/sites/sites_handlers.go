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