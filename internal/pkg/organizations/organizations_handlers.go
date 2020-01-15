package organizations

import (
	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ListOrganizationsRequest struct {
	// Input
	Db *sqlx.DB

	// Output
	Organizations []Organization
}

func (request *ListOrganizationsRequest) ListOrganizations() error {
	organizationRows := make([]OrganizationDbRow, 0)
	sqlStmt := request.Db.Rebind(listOrganizationsSql)
	err := request.Db.Select(&organizationRows, sqlStmt)
	if err != nil {
		return err
	}

	organizations := make([]Organization, len(organizationRows), len(organizationRows))
	for i, row := range organizationRows {
		organizations[i] = *row.CopyToOrganization()
	}
	request.Organizations = organizations
	return nil
}
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
	requestConfig := ListOrganizationsRequest{
		Db:            appConfig.DatabaseConnection,
		Organizations: nil,
	}
	err = requestConfig.ListOrganizations()
	if err != nil {
		logger.WithError(err).Error("Failed to query for organizations")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = response.WriteEntity(requestConfig.Organizations)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DescribeOrganizationRequest struct {
	// Input
	Db *sqlx.DB
	OrganizationId int

	// Output
	Organization *Organization
}
func (request *DescribeOrganizationRequest) DescribeOrganization() error {
	sqlStmt := request.Db.Rebind(describeOrganizationSql)
	var orgRow OrganizationDbRow
	err := request.Db.Get(&orgRow, sqlStmt, request.OrganizationId)
	if err != nil {
		// A 404 Not Found is returned when no error is returned, and no Organization is returned either.
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return err
	}

	request.Organization = orgRow.CopyToOrganization()
	return nil
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
	orgIdInt, err := strconv.Atoi(request.PathParameter("organizationId"))
	if err != nil {
		logger.WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusNotFound)
		return
	}

	searchConfig := DescribeOrganizationRequest{
		Db:             appConfig.DatabaseConnection,
		OrganizationId: orgIdInt,
		Organization:   nil,
	}

	err = searchConfig.DescribeOrganization()
	if err != nil {
		logger.WithError(err).Error("Failed to fetch org details")
		response.WriteHeader(http.StatusInternalServerError)
		return
	} else if searchConfig.Organization == nil {
		// If no error and no data returned, then that means the request was valid, but the ID was not in the DB
		response.WriteHeader(http.StatusNotFound)
		return
	}


	err = response.WriteEntity(*searchConfig.Organization)
	if err != nil {
		logger.WithError(err).Error("Failed to write organizations response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type UpdateOrganizationRequest struct {
	// Input
	Db *sqlx.DB
	InputOrganization *Organization

	// Output
	Organization *Organization
}
func (request *UpdateOrganizationRequest) UpdateOrganization() error {
	sqlStmt := request.Db.Rebind(updateOrganizationSql)
	res, err := request.Db.NamedExec(sqlStmt, request.InputOrganization)
	if err != nil {
		// A 404 Not Found is returned when no error is returned, and no Organization is returned either.
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return err
	}
	rowCount, err := res.RowsAffected()
	if rowCount == 0 {
		return nil
	}

	// Only save the updated org if the db updated a record
	request.Organization = request.InputOrganization
	return nil
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
