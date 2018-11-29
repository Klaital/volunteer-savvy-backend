package sites

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
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

type DailySchedule struct {
	Day       string `json:"date"`  // Expected format: YYYY-MM-DD
	OpenTime  string `json:"open"`  // Expected format: HH:MM
	CloseTime string `json:"close"` // Expected format: HH:MM
	IsOpen    bool   `json:"is_open"`
}
type WeeklySchedule struct {
	Monday    DailySchedule `json:"monday"`
	Tuesday   DailySchedule `json:"tuesday"`
	Wednesday DailySchedule `json:"wednesday"`
	Thursday  DailySchedule `json:"thursday"`
	Friday    DailySchedule `json:"friday"`
	Saturday  DailySchedule `json:"saturday"`
	Sunday    DailySchedule `json:"sunday"`
}
type Site struct {
	Id     int32  `json:"-" db:"id"`
	Slug   string `json:"slug" db:"slug"`
	Name   string `json:"name" db:"name_l10n"`
	Locale string `json:"locale" db:"locale"`

	Location

	IsActive bool `json:"active" db:"is_active"`

	// List of Site Coordinators/Managers
	Managers []users.User `json:"managers"`

	// Default Schedule
	DefaultSchedule WeeklySchedule `json:"default_schedule"`

	// Calendar Overrides
	CalendarOverrides []DailySchedule `json:"-"`

	// Computed Calendar
	Calendar []DailySchedule `json:"calendar"`
}
