package organizations

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func (o *Organization) Create(db *sqlx.DB) error {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "Organization.Save",
		"OrganizationId": o.Id,
		"OrganizationSlug": o.Slug,
	})
	isValid, errors := o.Validate()
	if !isValid {
		logger.Errorf("Cannot save - organiztaion failed validation with %d errors: %v", len(errors), errors)
		return fmt.Errorf("cannot save - organiztaion failed validation with %d errors: %v", len(errors), errors)
	}

	sqlStmt := db.Rebind(createOrganizationSql)
	_, err := db.NamedExec(sqlStmt, o)
	if err != nil {
		logger.WithError(err).Error("Failed to insert organization")
		return err
	}
	return nil
}

func FindOrganization(orgId uint64, db *sqlx.DB) (*Organization, error) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "FindOrganization",
		"OrganizationId": orgId,
	})

	var org OrganizationDbRow
	sqlStmt := db.Rebind(describeOrganizationSql)
	err := db.Get(&org, sqlStmt, orgId)
	if err != nil {
		// TODO: discern between db error and legit "not found"
		logger.WithError(err).Error("Failed to select organization")
		return nil, err
	}

	return org.CopyToOrganization(), nil
}
