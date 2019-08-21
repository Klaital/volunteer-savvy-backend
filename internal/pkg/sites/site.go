package sites

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	log "github.com/sirupsen/logrus"
)

type Location struct {
	Latitude      string `json:"lat" db:"lat"`
	Longitude     string `json:"lon" db:"lon"`
	GooglePlaceId string `json:"google_place_id" db:"gplace_id"`
	Street        string `json:"street" db:"street"`
	City          string `json:"city" db:"city"`
	State         string `json:"state" db:"state"`
	ZipCode       string `json:"zip" db:"zip"`
}

type DailyScheduleRow struct {
	Id          sql.NullInt64 `db:"id"`
	SiteId      sql.NullInt64 `db:"site_id"`
	DotwDefault sql.NullString `db:"dotw_default"` // Enum for the days of the week. If this column is not null, then it specifies a site's default schedule for that day of the week.
	Day         sql.NullString `db:"override_date"`  // Expected format: YYYY-MM-DD
	OpenTime    sql.NullString `db:"open_time"`  // Expected format: HH:MM
	CloseTime   sql.NullString `db:"close_time"` // Expected format: HH:MM
	IsOpen      sql.NullBool `db:"is_open"`
}

type DailySchedule struct {
	Id        int64 `json:"-" db:"id"`
	SiteId    int64 `json:"-" db:"site_id"`
	DotwDefault string `json:"-" db:"dotw_default"` // Enum for the days of the week. If this column is not null, then it specifies a site's default schedule for that day of the week.
	Day       string `json:"date" db:"override_date"`  // Expected format: YYYY-MM-DD
	OpenTime  string `json:"open" db:"open_time"`  // Expected format: HH:MM
	CloseTime string `json:"close" db:"close_time"` // Expected format: HH:MM
	IsOpen    bool   `json:"is_open" db:"is_open"`
}

//type WeeklySchedule struct {
//	Monday    DailySchedule `json:"monday"`
//	Tuesday   DailySchedule `json:"tuesday"`
//	Wednesday DailySchedule `json:"wednesday"`
//	Thursday  DailySchedule `json:"thursday"`
//	Friday    DailySchedule `json:"friday"`
//	Saturday  DailySchedule `json:"saturday"`
//	Sunday    DailySchedule `json:"sunday"`
//}
type Site struct {
	Id     uint64  `json:"-" db:"id"`
	Slug   string `json:"slug" db:"slug"`
	Name   string `json:"name" db:"name_l10n"`
	Locale string `json:"locale" db:"locale"`

	Location

	IsActive bool `json:"active" db:"is_active"`

	// List of Site Coordinators/Managers
	Managers []users.User `json:"managers"`

	// Default Schedule
	DefaultSchedule map[string]DailySchedule `json:"default_schedule"`

	// Calendar Overrides
	CalendarOverrides []DailySchedule `json:"-"`

	// Computed Calendar
	Calendar []DailySchedule `json:"calendar"`
}

type SiteCoordinator struct {
	Id uint64 `db:"id"`
	SiteId uint64 `db:"site_id"`
	UserId uint64 `db:"user_id"`
	IsPrimary bool `db:"is_primary"`
}

func FindSite(slug string) (site *Site, err error) {
	// Setup
	logger := log.WithFields(log.Fields{
		"operation": "FindSite",
		"slug": slug,
	})
	logger.Debug("Starting queries for site data")

	svcConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.Errorf("Service Config not loaded: %v", err)
		return nil, err
	}
	rows := make([]FindAllSitesRow,0)
	err = svcConfig.DatabaseConnection.Select(&rows,
		svcConfig.DatabaseConnection.Rebind(findSingleSiteSql),
		slug)
	if err != nil {
		logger.Errorf("Failed to select all sites: %v", err)
		return nil, err
	}

	// Sort the Sites, Managers, and Calendars into the nested structs we use
	sites := CoallateSiteSet(rows)

	return &sites[0], nil
}

func (site *Site) validate() bool {
	if site == nil {
		return false
	}
	if len(site.Slug) == 0 {
		return false
	}
	return true
}
func (site *Site) CreateSite() error {
	// Setup
	logger := log.WithFields(log.Fields{
		"operation": "CreateSite",
	})

	svcConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.Errorf("Service Config not loaded: %v", err)
		return err
	}
	db := svcConfig.DatabaseConnection
	tx, err := db.Beginx()
	if err != nil {
		logger.Errorf("Failed to create transaction: %v", err)
		return err
	} else {
		logger.Debug("Starting tx to create a site")
	}
	if tx == nil {
		logger.Error("No transaction handle returned")
		return errors.New("no transaction handle returned for CreateSite")
	}

	if !site.validate() {
		return errors.New("failed to validate site")
	}

	// Insert the Site itself
	rows, siteErr := db.NamedQuery(db.Rebind(insertSiteSql), site)
	if siteErr != nil {
		logger.Errorf("Failed to insert site: %v", siteErr)
		return siteErr
	}
	if rows.Next() {
		err = rows.Scan(&site.Id)
		if err != nil {
			logger.Errorf("Failed to scan site ID: %v", err)
			return err
		} else {
			logger.Debugf("Got inserted site ID: %d", site.Id)
		}
	} else {
		return fmt.Errorf("no site ID returned")
	}

	// TODO: Insert the Manager references

	// Insert the Default Schedule
	if _, defaultScheduleErr := db.Exec(db.Rebind(insertDefaultScheduleSql), site.Id, site.Id, site.Id, site.Id, site.Id, site.Id, site.Id); defaultScheduleErr != nil {
		logger.Errorf("Failed to insert default schedule: %v", defaultScheduleErr)
		return defaultScheduleErr
	}

	// Success!
	return nil
}

type FindAllSitesRow struct {
	Site

	CoordinatorGuid sql.NullString `db:"user_guid"`
	CoordinatorEmail sql.NullString `db:"email"`

	ScheduleDefaultDay sql.NullString `db:"dotw_default"`
	ScheduleOverride   sql.NullString `db:"override_date"`
	OpenTime sql.NullString `db:"open_time"`
	CloseTime sql.NullString `db:"close_time"`
	IsOpen   sql.NullBool `db:"is_open"`
}

func CoallateSiteSet(rows []FindAllSitesRow) []Site {
	sites := make(map[string]*Site)

	for _, row := range rows {
		thisSite := sites[row.Site.Slug]
		if thisSite == nil {
			thisSite = &Site{
				Id: row.Site.Id,
				Slug: row.Site.Slug,
				Name: row.Site.Name,
				Locale: row.Site.Locale,
				Location: row.Site.Location,
				IsActive: row.Site.IsActive,
			}
		}

		// Optional list of the Site's Managers
		if row.CoordinatorGuid.Valid {
			thisRowUser := users.User{
				Guid:         row.CoordinatorEmail.String,
				Email:        row.CoordinatorEmail.String,
			}
			thisSite.Managers = append(thisSite.Managers, thisRowUser)
		}

		// The Site's calendar
		// Any entries with dotw_default set are the defaults for that day of the week.
		// Any entries without dotw_default, but do have override_date, go on the calendar for that day.
		if row.ScheduleDefaultDay.Valid {
			thisRowCalendar := DailySchedule{
				SiteId:      0,
				DotwDefault: row.ScheduleDefaultDay.String,
				OpenTime:    row.OpenTime.String,
				CloseTime:   row.CloseTime.String,
				IsOpen:      row.IsOpen.Bool,
			}
			thisSite.DefaultSchedule[thisRowCalendar.DotwDefault] = thisRowCalendar
		} else if row.ScheduleOverride.Valid {
			thisRowCalendar := DailySchedule{
				Id:          0,
				Day:         row.ScheduleOverride.String,
				OpenTime:    row.OpenTime.String,
				CloseTime:   row.CloseTime.String,
				IsOpen:      row.IsOpen.Bool,
			}
			thisSite.CalendarOverrides = append(thisSite.CalendarOverrides, thisRowCalendar)
		}

		sites[thisSite.Slug] = thisSite
	}

	// flatten the map into an array
	siteList := make([]Site, 0, len(sites))
	for slug := range sites {
		siteList = append(siteList, *sites[slug])
	}

	return siteList
}
func FindAllSites() (sites []Site, err error) {
	logger := log.WithFields(log.Fields{
		"operation": "FindAllSites",
	})

	// Fetch the Sites, Managers, and Calendars from the database
	svcConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.Errorf("Failed to load service config: %v")
		return nil, err
	}
	rows := make([]FindAllSitesRow,0)
	err = svcConfig.DatabaseConnection.Select(&rows, findAllSitesSql)
	if err != nil {
		logger.Errorf("Failed to select all sites: %v", err)
		return nil, err
	}

	// Sort the Sites, Managers, and Calendars into the nested structs we use
	sites = CoallateSiteSet(rows)

	// success!
	return sites, nil
}

func DeleteSite(siteSlug string) error {
	logger := log.WithFields(log.Fields{
		"operation": "DeleteSite",
		"SiteSlug": siteSlug,
	})

	svcConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.Errorf("Failed to load service config: %v")
		return err
	}
	db := svcConfig.DatabaseConnection

	_, err = db.Exec(db.Rebind(deleteSiteSql), siteSlug)
	if err != nil {
		logger.Errorf("Failed to delete site: %v", err)
		return err
	}

	// Success!
	return nil
}