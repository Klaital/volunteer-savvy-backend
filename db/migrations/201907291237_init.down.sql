-- Organizations
DROP TABLE IF EXISTS organizations;

-- Sites
DROP INDEX IF EXISTS site_slug_index;
DROP TABLE IF EXISTS sites;

-- Users & Roles
DROP INDEX IF EXISTS user_guid_index;
DROP TABLE IF EXISTS users CASCADE;

-- Site Coordinators
DROP INDEX IF EXISTS sites_coordinators_index;
DROP INDEX IF EXISTS users_coordinated_sites_index;
DROP TABLE IF EXISTS site_coordinators;

-- Schedules

DROP INDEX IF EXISTS sites_schedules_index;
DROP TABLE IF EXISTS daily_schedules;
DROP TYPE IF EXISTS dotw_type;
