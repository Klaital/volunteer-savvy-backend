package sites


const findSiteSql = `
	SELECT
		id, slug, name_l10n, locale, 
		lat, lon, gplace_id, street, city, state, zip, 
		is_active 
	FROM sites WHERE slug=? LIMIT 1
`

const listAllSitesSql = `
	SELECT 
		sites.id, sites.slug, sites.name_l10n, sites.locale, 
		sites.lat, sites.lon, sites.gplace_id, sites.street, 
		sites.city, sites.state, sites.zip, sites.is_active,

		users.user_guid, users.email,

		daily_schedules.dotw_default, daily_schedules.override_date, 
		daily_schedules.open_time, daily_schedules.close_time,
		daily_schedules.is_open

	FROM sites 
		LEFT OUTER JOIN site_coordinators ON site_coordinators.site_id = sites.id
		LEFT OUTER JOIN users ON site_coordinators.user_id = users.id 
		LEFT OUTER JOIN daily_schedules on daily_schedules.site_id = sites.id
`
const describeSiteSql = listAllSitesSql + `
	WHERE sites.slug = ?	
`

const listOrganizationSitesSql = `
	SELECT 
		sites.id, sites.slug, sites.name_l10n, sites.locale, 
		sites.lat, sites.lon, sites.gplace_id, sites.street, 
		sites.city, sites.state, sites.zip, sites.is_active,

		users.user_guid, users.email,

		daily_schedules.dotw_default, daily_schedules.override_date, 
		daily_schedules.open_time, daily_schedules.close_time,
		daily_schedules.is_open

	FROM sites 
		LEFT OUTER JOIN site_coordinators ON site_coordinators.site_id = sites.id
		LEFT OUTER JOIN users ON site_coordinators.user_id = users.id 
		LEFT OUTER JOIN daily_schedules on daily_schedules.site_id = sites.id
	WHERE
		sites.organization_id = ?
`

const selectSiteCoordinatorsForSiteSql = `
	SELECT users.id, users.user_guid, users.email 
	FROM users JOIN site_coordinators ON users.id = site_coordinators.user_id
	WHERE site_coordinators.site_id = ?`

const selectSiteSchedulesSql = `
	SELECT id, site_id, dotw_default, override_date, open_time, close_time, is_open 
	FROM daily_schedules 
	WHERE site_id = ?
`

const insertSiteSql = `
	INSERT INTO sites (
		slug, name_l10n, locale, lat, lon, gplace_id, street, city, state, zip, is_active
	) VALUES (
		:slug, :name_l10n, :locale, :lat, :lon, :gplace_id, :street, :city, :state, :zip, :is_active
	) RETURNING id
`

const insertDefaultScheduleSql = `
	INSERT INTO daily_schedules 
		(site_id, dotw_default, override_date, open_time, close_time, is_open)
	VALUES
		(?, 'sunday', null, '09:00', '17:00', true),
		(?, 'monday', null, '09:00', '17:00', true),
		(?, 'tuesday', null, '09:00', '17:00', true),
		(?, 'wednesday', null, '09:00', '17:00', true),
		(?, 'thursday', null, '09:00', '17:00', true),
		(?, 'friday', null, '09:00', '17:00', true),
		(?, 'saturday', null, '09:00', '17:00', true)
`

const deleteSiteSql = `
	DELETE FROM sites WHERE slug = ?
`
const updateSiteSql = `
	UPDATE sites SET 
		name_l10n = :name_l10n,
		locale = :locale,
		lat = :lat,
		lon = :lon,
		gplace_id = :gplace_id,
		street = :street,
		city = :city,
		state = :state,
		zip = :zip,
		is_active = :is_active
	WHERE slug = :slug
`
