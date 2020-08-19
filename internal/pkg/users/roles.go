package users

type RoleType int

const (
	SiteAdmin   RoleType = 0 // Administrative permissions for Volunteer-Savvy as a whole
	OrgAdmin    RoleType = 1 // Administrative permissions for a single Organization
	Volunteer   RoleType = 2 // User is able to sign up, log work.
	SiteManager RoleType = 3 // User is able to sign up as a Site Coordinator for sites, then manage those sites' settings.
	BackOffice  RoleType = 4 // User is able to log work, read and update suggestions, generate reports. Not able to modify Users or Site settings.
	Mobile      RoleType = 5 // User is interested in working at the Mobile sites. Enables the user to opt-in to notifications about mobile sites specifically.
)

type Role struct {
	Id       uint64   `json:"-" db:"id"`
	OrgId    uint64   `json:"org_id" db:"org_id"`
	UserId   uint64   `json:"-" db:"user_id"`
	UserGuid string   `json:"user_guid"`
	Role     RoleType `json:"name" db:"name"`
}
