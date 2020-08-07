-- Organizations
DROP TABLE IF EXISTS organizations;

-- Sites
DROP INDEX IF EXISTS site_slug_index;
DROP TABLE IF EXISTS sites;

-- Users & Roles
DROP INDEX IF EXISTS user_guid_index;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE ;
DROP INDEX IF EXISTS roles_users_index;
DROP INDEX IF EXISTS roles_org_users_index;
DROP INDEX IF EXISTS roles_unique_index;


-- Site Coordinators
DROP INDEX IF EXISTS sites_coordinators_index;
DROP INDEX IF EXISTS users_coordinated_sites_index;
DROP TABLE IF EXISTS site_coordinators;

-- Schedules

DROP INDEX IF EXISTS sites_schedules_index;
DROP TABLE IF EXISTS daily_schedules;
DROP TYPE IF EXISTS dotw_type;
