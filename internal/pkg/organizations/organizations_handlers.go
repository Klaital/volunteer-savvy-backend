package organizations

import (
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)
func ListOrganizationsHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "ListOrganizationsHandler",
	})
	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	organizations := make([]Organization, 0)
	sqlStmt := appConfig.DatabaseConnection.Rebind(listOrganizationsSql)
	err = appConfig.DatabaseConnection.Select(&organizations, sqlStmt)
	if err != nil {
		logger.WithError(err).Error("Failed to query organizations")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = response.WriteEntity(organizations)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DescribeOrganizationHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "ListOrganizationsHandler",
		"OrganizationID": request.PathParameter("organizationId"),
	})
	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var organization Organization
	orgIdInt, err := strconv.Atoi(request.PathParameter("organizationId"))
	if err != nil {
		logger.WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusNotFound)
		return
	}

	sqlStmt := appConfig.DatabaseConnection.Rebind(describeOrganizationSql)
	err = appConfig.DatabaseConnection.Get(&organization, sqlStmt, orgIdInt)
	if err != nil {
		logger.WithError(err).Error("Failed to select organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = response.WriteEntity(organization)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateOrganizationHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "UpdateOrganizationHandler",
		"OrganizationID": request.PathParameter("organizationId"),
	})
	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	var organization Organization
	err = request.ReadEntity(organization)
	if err != nil {
		logger.WithError(err).Error("Failed to unmarshal request body")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	orgIdInt, err := strconv.Atoi(request.PathParameter("organizationId"))
	if err != nil {
		logger.WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusNotFound)
		return
	}
	organization.Id = uint64(orgIdInt)
	sqlStmt := appConfig.DatabaseConnection.Rebind(updateOrganizationSql)
	_, err = appConfig.DatabaseConnection.NamedExec(sqlStmt, &organization)
	if err != nil {
		logger.WithError(err).Error("Failed to update organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = response.WriteEntity(organization)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DestroyOrganizationHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "DestroyOrganizationsHandler",
		"OrganizationID": request.PathParameter("organizationId"),
	})
	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	orgIdInt, err := strconv.Atoi(request.PathParameter("organizationId"))
	if err != nil {
		logger.WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusNotFound)
		return
	}
	sqlStmt := appConfig.DatabaseConnection.Rebind(deleteOrganizationNullFkeysSql)
	_, err = appConfig.DatabaseConnection.NamedExec(sqlStmt, &Organization{Id:uint64(orgIdInt)})
	if err != nil {
		// TODO: differentiate between "not found" and "db error". Currently this code assumes the database never fails - it's always due to a bad ID
		logger.WithError(err).Error("Failed to select organization")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	response.WriteHeader(http.StatusOK)
}


func CreateOrganizationHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "ListOrganizationsHandler",
		"OrganizationID": request.PathParameter("organizationId"),
	})
	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	var organization Organization
	err = request.ReadEntity(&organization)
	if err != nil {
		logger.WithError(err).Error("Failed to read request body")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	sqlStmt := appConfig.DatabaseConnection.Rebind(createOrganizationSql)
	_, err = appConfig.DatabaseConnection.NamedExec(sqlStmt, &organization)
	if err != nil {
		logger.WithError(err).Error("Failed to update organization")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = response.WriteEntity(organization)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
