package organizations

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/sirupsen/logrus"
)

func ListOrganizations(ctx context.Context, db *sqlx.DB) ([]Organization, error) {
	// Fetch the logger from the context
	logger := filters.GetContextLogger(ctx)

	// Init complete. Start DB operations
	organizationSet := make([]Organization, 0)
	sqlStmt := db.Rebind(listOrganizationsSql)
	err := db.Select(&organizationSet, sqlStmt)
	if err != nil {
		logger.WithError(err).Error("Failed to find organizaitons")
		return organizationSet, err
	}
	return organizationSet, nil
}

func DescribeOrganization(ctx context.Context, db *sqlx.DB, organizationID int64) (*Organization, error) {
	sqlStmt := db.Rebind(describeOrganizationSql)
	var orgRow OrganizationDbRow
	err := db.Get(&orgRow, sqlStmt, organizationID)
	if err != nil {
		return nil, err
	}

	org := orgRow.CopyToOrganization()
	return org, nil
}

func (o *Organization) Create(ctx context.Context, db *sqlx.DB) error {
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation":        "Organization.Create",
		"OrganizationId":   o.Id,
		"OrganizationSlug": o.Slug,
	})
	errorSet := o.Validate()
	if errorSet != nil {
		logger.WithError(errorSet.Errors[0]).Errorf("Cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
		return fmt.Errorf("cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
	}

	sqlStmt := db.Rebind(createOrganizationSql)
	_, err := db.NamedExec(sqlStmt, o)
	if err != nil {
		logger.WithError(err).Error("Failed to insert organization")
		return err
	}
	return nil
}

func (o *Organization) Update(ctx context.Context, db *sqlx.DB) error {
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation":        "Organization.Update",
		"OrganizationID":   o.Id,
		"OrganizationSlug": o.Slug,
	})
	errorSet := o.Validate()
	if errorSet != nil {
		logger.WithError(errorSet.Errors[0]).Errorf("Cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
		return fmt.Errorf("cannot save - organiztaion failed validation with %d errors: %v", len(errorSet.Errors), errorSet)
	}

	sqlStmt := db.Rebind(updateOrganizationSql)
	_, err := db.NamedExec(sqlStmt, o)
	if err != nil {
		logger.WithError(err).Error("Failed to update organization")
		return err
	}

	// Success!
	return nil
}

func DeleteOrganization(ctx context.Context, organizationID uint64, db *sqlx.DB) error {
	logger := filters.GetContextLogger(ctx).WithFields(logrus.Fields{
		"operation":      "DeleteOrganization",
		"OrganizationID": organizationID,
	})

	org := Organization{
		Id: organizationID,
	}
	sqlStmt := db.Rebind(deleteOrganizationNullFkeysSql)
	_, err := db.NamedExec(sqlStmt, &org)
	if err != nil {
		logger.WithError(err).Error("Error deleting organization")
		return err
	}

	// Success!
	return nil
}
