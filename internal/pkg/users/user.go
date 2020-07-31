package users

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	log "github.com/sirupsen/logrus"
	"github.com/klaital/intmath"
)

type User struct {
	Id    uint64  `json:"-" db:"id"`
	Guid  string `json:"user_guid" db:"user_guid"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_digest"`

	Roles map[uint64][]Role `json:"roles"` // the map key is the organization ID
}

type RoleType int
const (
	SiteAdmin  RoleType = 0 // Administrative permissions for Volunteer-Savvy as a whole

	OrgAdmin    RoleType = 1 // Administrative permissions for a single Organization
	Volunteer   RoleType = 2 // User is able to sign up, log work.
	SiteManager RoleType = 3 // User is able to sign up as a Site Coordinator for sites, then manage those sites' settings.
	BackOffice  RoleType = 4 // User is able to log work, read and update suggestions, generate reports.
	                         // Not able to modify Users or Site settings.
 	Mobile      RoleType = 5 // User is interested in working at the Mobile sites.
 	                         // Enables the user to opt-in to notifications about mobile sites specifically.
)
type Role struct {
	Id     uint64 `json:"-" db:"id"`
	OrgId  uint64 `json:"organization_id" db:"org_id"`
	UserId uint64 `json:"-" db:"user_id"`
	UserGuid string `json:"user_guid"`
	Role   RoleType `json:"name" db:"name"`
}

func GetUserForLogin(ctx context.Context, email string, db *sqlx.DB) (*User, error) {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "GetUserForLogin",
		"email": email,
	})

	var u User
	sqlStmt := db.Rebind(`SELECT id, user_guid, password_digest FROM users WHERE email = ?`)
	err := db.Get(&u, sqlStmt, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.WithError(err).Error("Failed to select user with email/password")
		return nil, err
	}

	// Success!
	return &u, nil
}

// GetUserRoles fetches all permissions granted to the user, sorted by the
// Organization ID they are granted on. If an Organization ID is not found
// among the keys, the user does not have any access to that org.
func (u *User) GetRoles(ctx context.Context, db *sqlx.DB) (map[uint64][]Role, error) {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "GetUserRoles",
		"UserID": u.Id,
		"UserGuid": u.Guid,
	})

	// If the user's roles have already been loaded, just use that
	if len(u.Roles) > 0 {
		return u.Roles, nil
	}

	var roles []Role
	sqlStmt := db.Rebind(`SELECT id, org_id, user_id, name FROM roles WHERE user_id = ?`)
	err := db.Select(&roles, sqlStmt, u.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return map[uint64][]Role{}, nil
		}
		logger.WithError(err).Error("Error fetching user roles")
		return nil, err
	}

	// Sort the roles by Organization
	mappedRoles := make(map[uint64][]Role, 1)
	for i, role := range roles {
		existingRoles, ok := mappedRoles[role.OrgId]
		roles[i].UserGuid = u.Guid
		if ok {
			existingRoles = append(existingRoles, roles[i])
		} else {
			existingRoles = []Role{roles[i]}
		}
		mappedRoles[role.OrgId] = existingRoles
	}

	// Success!
	u.Roles = mappedRoles // cache them for future use
	return mappedRoles, nil
}

func ListUsersInSameOrgs(ctx context.Context, userGuid string, db *sqlx.DB) ([]User, error) {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "ListUsersInSameOrgs",
	})

	// Pull the roles belonging to the user
	var roles []Role
	sqlStmt := db.Rebind(`
SELECT 
	r.id AS id, 
	r.org_id AS org_id, 
	r.user_id AS user_id, 
	r.name AS name 
FROM roles AS r JOIN users AS u 
	ON r.user_id = u.id
WHERE u.user_guid = ?`)

	err := db.Select(&roles, sqlStmt, userGuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []User{}, nil
		}
		return []User{}, err
	}

	// Sort the roles by Organization
	orgIdSet := intmath.NewSet()
	for _, role := range roles {
		orgIdSet.Add(int64(role.OrgId))
	}

	sqlStmt = `
SELECT 
	u.id, u.user_guid, u.email
FROM users AS u JOIN roles AS r 
	ON u.id = r.user_id
WHERE r.org_id IN (?)`
	sqlStmt, args, err := sqlx.In(sqlStmt, orgIdSet.GetItems())
	if err != nil {
		logger.WithError(err).Error("Failed to compile IN query")
		return []User{}, err
	}

	// Load the users from the DB
	var users []User
	err = db.Select(&users, sqlStmt, args...)

	// Load those user's Roles
	// OPTIMIZATION: make this a JOIN query with the users SELECT
	for _, user := range users {
		_, err := user.GetRoles(ctx, db)
		if err != nil {
			logger.WithField("UserGUID", user.Guid).WithError(err).Error("Failed to load roles for user")
			// don't block return on a failure here, just return what we do have
		}
	}

	return users, nil
}
