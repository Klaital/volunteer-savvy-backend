package organizations

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"regexp"
)

type Organization struct {
	Id uint64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Authcode string `json:"authcode" db:"authcode"`

	// Contact info
	ContactUserId uint64 `json:"contact_user_id" db:"contact_user_id"`
	ContactUser *users.User `json:"contact"`

	// Geographical Center - used for map view defaults
	Latitude float64 `json:"lat" db:"lat"`
	Longitude float64 `json:"lon" db:"lon"`
}

type OrganizationDbRow struct {
	Id uint64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Authcode string `json:"authcode" db:"authcode"`

	// Contact info
	ContactUserId sql.NullInt64 `json:"contact_user_id" db:"contact_user_id"`

	// Geographical Center - used for map view defaults
	Latitude float64 `json:"lat" db:"lat"`
	Longitude float64 `json:"lon" db:"lon"`
}
func (row OrganizationDbRow) CopyToOrganization() *Organization {
	o := Organization{
		Id:            row.Id,
		Name:          row.Name,
		Slug:          row.Slug,
		Authcode:      row.Authcode,
		ContactUserId: 0,
		Latitude:      row.Latitude,
		Longitude:     row.Longitude,
	}

	if row.ContactUserId.Valid {
		o.ContactUserId = uint64(row.ContactUserId.Int64)
	}
	return &o
}

func New() *Organization {
	return &Organization{
		Id:            0,
		Name:          "",
		Slug:          "",
		Authcode:      "",
		ContactUserId: 0,
		Latitude:      0,
		Longitude:     0,
	}
}

func (o Organization) Validate() (valid bool, errs []error) {
	errs = make([]error, 0)
	if len(o.Name) == 0 {
		errs = append(errs, errors.New("name must be present"))
	}
	if len(o.Slug) == 0 {
		errs = append(errs, errors.New("slug must be present"))
	}
	pattern, err := regexp.Compile("^[a-z0-9]+(?:-[a-z0-9]+)*$")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed to compile regex: %v", err))
		return false, errs
	}
	if !pattern.Match([]byte(o.Slug)) {
		errs = append(errs, errors.New("slug must match pattern /^[a-z0-9]+(?:-[a-z0-9]+)*$/"))
	}
	if len(o.Authcode) == 0 {
		errs = append(errs, errors.New("authcode must be present"))
	}

	return len(errs) == 0, errs
}