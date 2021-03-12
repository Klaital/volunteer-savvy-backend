package users

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/intmath"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Id           uint64 `json:"-" db:"id"`
	Guid         string `json:"user_guid" db:"user_guid"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_digest"`

	Roles map[uint64][]Role `json:"roles"` // the map key is the organization ID
}

// FindUser queries the database for the user and all other data needed to
// display their profile.
func FindUser(ctx context.Context, email string, db *sqlx.DB) (*User, error) {
	user, err := GetUserForLogin(ctx, email, db)
	if err != nil {
		return nil, err
	}
	_, err = user.GetRoles(ctx, db)
	return user, err
}

func GetUserForLogin(ctx context.Context, email string, db *sqlx.DB) (*User, error) {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "GetUserForLogin",
		"email":     email,
	})

	var u User
	sqlStmt := db.Rebind(`SELECT id, user_guid, password_digest FROM users WHERE email = ?`)
	err := db.Get(&u, sqlStmt, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.WithError(err).Error("Failed to select user with email")
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
		"UserID":    u.Id,
		"UserGuid":  u.Guid,
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

func ListUsersInSameOrgs(ctx context.Context, jwtClaims *Claims, db *sqlx.DB) ([]User, error) {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "ListUsersInSameOrgs",
	})

	// Compile a list of all orgs where the user claims any Role
	orgIdSet := intmath.NewSet()
	for orgId := range jwtClaims.Roles {
		orgIdSet.Add(int64(orgId))
	}

	if orgIdSet.Length() == 0 {
		logger.WithField("Claims", *jwtClaims).Debug("No claims to find users against")
		return []User{}, nil
	}
	logger = logger.WithField("OrgIdSet", orgIdSet)
	logger.Debug("Generated Org ID list for user")

	sqlStmt, args, err := sqlx.In(`
SELECT DISTINCT
	u.id, u.user_guid, u.email
FROM users AS u JOIN roles AS r 
	ON u.id = r.user_id
WHERE r.org_id IN (?)`, orgIdSet.GetItems())
	if err != nil {
		logger.WithError(err).Error("Failed to compile IN query")
		return []User{}, err
	}

	// Load the users from the DB
	var users []User
	err = db.Select(&users, db.Rebind(sqlStmt), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("SQL", sqlStmt).WithError(err).Debug("No users returned")
			return []User{}, nil
		}
		logger.WithFields(log.Fields{
			"sqlstmt": sqlStmt,
			"orgIds":  args,
		}).WithError(err).Error("Failed to query users")
		return []User{}, err
	}

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
