package server

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ListOrganizationsResponse struct {
	Organizations []organizations.Organization `json:"organizations"`
}
func (server *Server) ListOrganizationsHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.
	orgs, err := organizations.ListOrganizations(ctx)

	if err != nil {
		server.Config.Logger.WithError(err).Error("Failed to fetch organizations list")
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
		server.Config.Logger.WithError(err).Error("Failed to write response payload")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) DescribeOrganizationHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	orgIDstr := request.PathParameter("organizationID")
	if len(orgIDstr) == 0 {
		ctx.Logger.Warn("This shouldn't happen. An empty organization ID has been passed to DescribeOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	orgID, err := strconv.Atoi(orgIDstr)
	if err != nil {
		ctx.Logger.WithField("orgID", orgIDstr).WithError(err).Error("Invalid org ID")
		response.WriteHeader(http.StatusBadRequest)
		return
	}



	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Fetch the organization data
	org, err := organizations.DescribeOrganization(ctx, int64(orgID))

	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to fetch organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(org)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) CreateOrganizationHandler(request *restful.Request, response *restful.Response) {
	// Set up the context for this request thread
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	newOrg := organizations.New()
	err := request.ReadEntity(newOrg)
	if err != nil {
		// TODO: maybe do input validation in a filter function?
		server.Config.Logger.WithError(err).Error("Failed to deserialize organization")
		if server.Config.LogLevel == "debug" {
			err = response.WriteError(http.StatusBadRequest, err)
			if err != nil {
				server.Config.Logger.WithError(err).Error("Failed to serialize error")
				response.WriteHeader(http.StatusBadRequest)
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}


	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Check whether the requested values form a valid Organization
	errorSet := newOrg.Validate()
	if errorSet != nil {
		ctx.Logger.WithError(errorSet.Errors[0]).Errorf("Specified org is not valid - %s", errorSet.Error())
		if ctx.Config.Logger.Level == logrus.DebugLevel {
			response.WriteError(http.StatusBadRequest, errorSet)
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// Publish the new Org to the DB
	err = newOrg.Create(ctx)
	if err != nil {
		server.Config.Logger.WithError(err).Error("Failed to create organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(newOrg)
	if err != nil {
		server.Config.Logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}


func (server *Server) UpdateOrganizationHandler(request *restful.Request, response *restful.Response) {
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	orgID := request.PathParameter("organizationID")
	if len(orgID) == 0 {
		ctx.Logger.Warn("This shouldn't happen. An empty organization ID has been passed to UpdateOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	newOrg := organizations.New()
	err := request.ReadEntity(newOrg)
	if err != nil {
		// TODO: maybe do input validation in a filter function?
		ctx.Logger.WithError(err).Error("Failed to deserialize organization input")
		if server.Config.LogLevel == "debug" {
			err = response.WriteError(http.StatusBadRequest, err)
			if err != nil {
				ctx.Logger.WithError(err).Error("Failed to serialize error")
				response.WriteHeader(http.StatusBadRequest)
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}


	// TODO: get logged-in user and add it to the context so that permissions and scope can be determined.

	// Publish the updates to the DB
	err = newOrg.Update(ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	err = response.WriteEntity(newOrg)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to serialize organization")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) DeleteOrganizationHandler(request *restful.Request, response *restful.Response) {
	ctx := config.NewContext(context.Background(), server.Config)
	ctx.Logger = request.Attribute("Logger").(*logrus.Entry)

	// Get the ID for the requested Organization
	orgIDstr := request.PathParameter("organizationID")
	if len(orgIDstr) == 0 {
		ctx.Logger.Warn("This shouldn't happen. An empty organization ID has been passed to DescribeOrganizationHandler")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	orgID, err := strconv.Atoi(orgIDstr)
	if err != nil {
		ctx.Logger.WithError(err).Errorf("Invalid Org ID given: %s", orgIDstr)
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	// TODO: get logged-in user and check their permissions

	// Fetch the organization data
	err = organizations.DeleteOrganization(ctx, uint64(orgID))
	if err != nil {
		// TODO: check the err type to discern between "not found" and "there was a DB error"
		server.Config.Logger.WithError(err).Error("Failed to fetch organization")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format and send the response
	response.WriteHeader(http.StatusOK)
}
