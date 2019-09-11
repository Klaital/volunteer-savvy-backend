package users

type User struct {
	Id    int32  `json:"-" db:"id"`
	OrganizationId uint64 `json:"organization_id" db:"organization_id"`
	Guid  string `json:"user_guid" db:"user_guid"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_digest"`

	Roles []Role `json:"roles"`
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
	Id     int32 `json:"-" db:"id"`
	OrgId  int32 `json:"organization_id" db:"org_id"`
	UserId int32 `json:"-" db:"user_id"`
	UserGuid string `json:"user_guid"`
	Role   RoleType `json:"name" db:"name"`
}

