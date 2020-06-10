package organizations

import (
	"fmt"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
)


func ListOrganizations(ctx *config.Context) ([]Organization, error) {
	// Fetch the logger and DB connection from the context
	db := ctx.Config.GetDbConn()

	// Init complete. Start DB operations
	ctx.Logger.Debug("Starting ListOrganizations")

	organizationSet := make([]Organization, 0)
	sqlStmt := db.Rebind(listOrganizationsSql)
	err := db.Select(&organizationSet, sqlStmt)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to find organizaitons")
		return organizationSet, err
	}
	return organizationSet, nil
}

func DescribeOrganization(ctx *config.Context, organizationID int64) (*Organization, error) {
	db := ctx.Config.GetDbConn()

	sqlStmt := db.Rebind(describeOrganizationSql)
	var orgRow OrganizationDbRow
	err := db.Get(&orgRow, sqlStmt, organizationID)
	if err != nil {
		// A 404 Not Found is returned when no error is returned, and no Organization is returned either.
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	org := orgRow.CopyToOrganization()
	return org, nil
}

func (o *Organization) Create(ctx *config.Context) error {
	logger := ctx.Logger.WithFields(logrus.Fields{
		"operation": "Organization.Create",
		"OrganizationId": o.Id,
		"OrganizationSlug": o.Slug,
	})
	errorSet := o.Validate()
	if errorSet != nil {
		logger.WithError(errorSet.Errors[0]).Errorf("Cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
		return fmt.Errorf("cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
	}

	db := ctx.Config.GetDbConn()

	sqlStmt := db.Rebind(createOrganizationSql)
	_, err := db.NamedExec(sqlStmt, o)
	if err != nil {
		logger.WithError(err).Error("Failed to insert organization")
		return err
	}
	return nil
}

func (o *Organization) Update(ctx *config.Context) error {
	logger := ctx.Logger.WithFields(logrus.Fields{
		"operation": "Organization.Update",
		"OrganizationID": o.Id,
		"OrganizationSlug": o.Slug,
	})
	errorSet := o.Validate()
	if errorSet != nil {
		logger.WithError(errorSet.Errors[0]).Errorf("Cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
		return fmt.Errorf("cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
	}

	db := ctx.Config.GetDbConn()
	sqlStmt := db.Rebind(updateOrganizationSql)
	_, err := db.NamedExec(sqlStmt, o)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update organization")
		return err
	}

	// Success!
	return nil
}

func DeleteOrganization(ctx *config.Context, organizationID uint64) error {
	logger := ctx.Logger.WithFields(logrus.Fields{
		"operation": "DeleteOrganization",
		"OrganizationID": organizationID,
	})

	org := Organization{
		Id: organizationID,
	}
	db := ctx.Config.GetDbConn()
	sqlStmt := db.Rebind(deleteOrganizationNullFkeysSql)
	_, err := db.NamedExec(sqlStmt, &org)
	if err != nil {
		logger.WithError(err).Error("Error deleting organization")
		return err
	}

	// Success!
	return nil
}