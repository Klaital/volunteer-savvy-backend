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
