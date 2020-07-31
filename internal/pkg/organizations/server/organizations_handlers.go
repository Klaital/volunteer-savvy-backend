package server

import (
	"database/sql"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/auth"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/version"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type OrganizationsServer struct {
	ApiVersion string
	Config *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *OrganizationsServer {
	return &OrganizationsServer{
		ApiVersion: version.Version,
		Config: cfg,
	}
}

//
// Call this in the server setup function
//
func (server *OrganizationsServer) GetOrganizationsAPI() *restful.WebService{

	service := new(restful.WebService)
	service.Path(server.Config.BasePath + "/organizations").ApiVersion(server.ApiVersion).Doc("Volunteer-Savvy Backend")

	//
	// Organizations APIs
	//
	service.Route(
		service.GET("/").
			//Filter(filters.RateLimitingFilter).
			To(server.ListOrganizationsHandler).
			Doc("List Organizations").
			Produces(restful.MIME_JSON).
			Writes(ListOrganizationsResponse{}).
			Returns(http.StatusOK, "Fetched all organizations", []organizations.Organization{}))
	// TODO: CreateOrganization needs a SuperAdmin permissions check filter
	service.Route(
		service.POST("/").
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireSuperAdminPermission).
			To(server.CreateOrganizationHandler).
			Doc("Create Organization").
			Produces(restful.MIME_JSON).
			Consumes(restful.MIME_JSON).
			Reads(organizations.Organization{}).
			Writes(organizations.Organization{}).
			Returns(http.StatusOK, "Organization created.", organizations.Organization{}))
	service.Route(
		service.GET("/{organizationID}").
			//Filter(filters.RateLimitingFilter).
			To(server.DescribeOrganizationHandler).
			Doc("Describe Organization").
			Param(restful.PathParameter("organizationID", "ID taken from ListOrganizations")).
			Produces(restful.MIME_JSON).
			Writes(organizations.Organization{}).
			Returns(http.StatusOK, "Organization details fetched", organizations.Organization{}).
			Returns(http.StatusNotFound, "Invalid Organization ID", nil))
	// TODO: UpdateOrganization needs a SuperAdmin permissions check filter
	service.Route(
		service.PUT("/{organizationID}").
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireSuperAdminPermission).
			To(server.UpdateOrganizationHandler).
			Doc("Update Organization").
			Param(restful.PathParameter("organizationId", "ID taken from ListOrganizations")).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(organizations.Organization{}).
			Writes(organizations.Organization{}).
			Returns(http.StatusOK, "Organization details updated", organizations.Organization{}).
			Returns(http.StatusBadRequest, "Unable to set the requested values.", nil).
			Returns(http.StatusNotFound, "Invalid Organization ID", nil))
	// TODO: DeleteOrganizations needs a SuperAdmin permissions check filter
	service.Route(
		service.DELETE("/{organizationID}").
			Filter(auth.ValidJwtFilter).
			//Filter(filters.RateLimitingFilter).
			//Filter(filters.RequireSuperAdminPermission).
			To(server.DeleteOrganizationHandler).
			Doc("Destroy Organization").
			Param(restful.PathParameter("organizationID", "ID taken from ListOrganizations")).
			Returns(http.StatusOK, "Organization deleted", nil).
			Returns(http.StatusNotFound, "Invalid Organization ID", nil))

	return service
}
type ListOrganizationsResponse struct {
	Organizations []organizations.Organization `json:"organizations"`
}
func (server *OrganizationsServer) ListOrganizationsHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation": "ListOrganizationsHandler",
	})

	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.
	orgs, err := organizations.ListOrganizations(ctx, server.Config.GetDbConn())

	if err != nil {
		logger.WithError(err).Error("Failed to fetch organizations list")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	// TODO: handle pagination
	responseBody := ListOrganizationsResponse{
		Organizations: orgs,
	}
	err = response.WriteEntity(responseBody)
	if err != nil {
		logger.WithError(err).Error("Failed to write response payload")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *OrganizationsServer) DescribeOrganizationHandler(request *restful.Request, response *restful.Response) {
	orgIDstr := request.PathParameter("organizationID")

	// Set up the context for this request thread
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation": "DescribeOrganizationsHandler",
		"OrganizationID.input": orgIDstr,
	})

	if len(orgIDstr) == 0 {
		logger.Warn("This shouldn't happen. An empty organization ID has been passed to DescribeOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	orgID, err := strconv.Atoi(orgIDstr)
	logger = logger.WithField("OrganizationID", orgID)
	if err != nil {
		logger.WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if orgID <= 0 {
		logger.WithError(errors.New("invalid org ID given")).Debug("invalid org ID")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Fetch the organization data
	org, err := organizations.DescribeOrganization(ctx, server.Config.GetDbConn(), int64(orgID))
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		logger.WithError(err).Error("Failed to fetch organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(org)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *OrganizationsServer) CreateOrganizationHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation": "CreateOrganizationHandler",
	})

	newOrg := organizations.New()
	err := request.ReadEntity(newOrg)
	if err != nil {
		// TODO: maybe do input validation in a filter function?
		logger.WithError(err).Error("Failed to deserialize organization")
		if logger.Level == logrus.DebugLevel {
			err = response.WriteError(http.StatusBadRequest, err)
			if err != nil {
				logger.WithError(err).Error("Failed to serialize error")
				response.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	}


	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Check whether the requested values form a valid Organization
	errorSet := newOrg.Validate()
	if errorSet != nil {
		logger.WithError(errorSet.Errors[0]).Errorf("Specified org is not valid - %s", errorSet.Error())
		if logger.Level == logrus.DebugLevel {
			response.WriteError(http.StatusBadRequest, errorSet)
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// Publish the new Org to the DB
	err = newOrg.Create(ctx, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to create organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(newOrg)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}


func (server *OrganizationsServer) UpdateOrganizationHandler(request *restful.Request, response *restful.Response) {
	orgID := request.PathParameter("organizationID")
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation": "UpdateOrganizationHandler",
		"OrganizationID.input": orgID,
	})

	if len(orgID) == 0 {
		logger.Warn("This shouldn't happen. An empty organization ID has been passed to UpdateOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	newOrg := organizations.New()
	err := request.ReadEntity(newOrg)
	if err != nil {
		// TODO: maybe do input validation in a filter function?
		logger.WithError(err).Error("Failed to deserialize organization input")
		if logger.Level == logrus.DebugLevel {
			err = response.WriteError(http.StatusBadRequest, err)
			if err != nil {
				logger.WithError(err).Error("Failed to serialize error")
				response.WriteHeader(http.StatusBadRequest)
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}


	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Publish the updates to the DB
	err = newOrg.Update(ctx, server.Config.GetDbConn())
	if err != nil {
		logger.WithError(err).Error("Failed to create organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(newOrg)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *OrganizationsServer)  DeleteOrganizationHandler(request *restful.Request, response *restful.Response) {
	orgIDstr := request.PathParameter("organizationID")
	ctx := filters.GetRequestContext(request)
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation": "DeleteOrganizationHandler",
		"OrganizationID.input": orgIDstr,
	})

	// Get the ID for the requested Organization
	if len(orgIDstr) == 0 {
		logger.Warn("This shouldn't happen. An empty organization ID has been passed to DescribeOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	orgID, err := strconv.Atoi(orgIDstr)
	logger = logger.WithField("OrganizationID", orgID)
	if err != nil {
		logger.WithError(err).Debug("Invalid Org ID given")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	if orgID <= 0 {
		logger.WithError(errors.New("invalid org ID")).Debug("Negative Org ID given")
		response.WriteHeader(http.StatusBadRequest)
	}

	// TODO: get logged-in user and check their permissions

	// Fetch the organization data
	err = organizations.DeleteOrganization(ctx, uint64(orgID), server.Config.GetDbConn())
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteHeader(http.StatusNotFound)
		} else {
			logger.WithError(err).Error("Failed to fetch organization")
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Format and send the response
	response.WriteHeader(http.StatusOK)
}
