-- Organizations

CREATE TABLE organizations (
   id SERIAL PRIMARY KEY,
   name VARCHAR(128) NOT NULL,
   slug VARCHAR(64) NOT NULL,
   authcode VARCHAR(64) UNIQUE NOT NULL,

   contact_user_id INTEGER, -- REFERENCES users(id),
   lat FLOAT,
   lon FLOAT
);

-- Sites

CREATE TABLE sites (
  id SERIAL PRIMARY KEY,
  organization_id INTEGER REFERENCES organizations(id),
  slug VARCHAR(64) UNIQUE NOT NULL,

  name_l10n VARCHAR(128) NOT NULL,
  locale VARCHAR(128) NOT NULL,

  is_active BOOLEAN NOT NULL DEFAULT false,

  -- Location
  lat VARCHAR(64) NOT NULL DEFAULT '',
  lon VARCHAR(64) NOT NULL DEFAULT '',
  gplace_id VARCHAR(64) NOT NULL DEFAULT '',
  street VARCHAR(64) NOT NULL DEFAULT '',
  city VARCHAR(32) NOT NULL DEFAULT '',
  state VARCHAR(16) NOT NULL DEFAULT '',
  zip VARCHAR(8) NOT NULL DEFAULT ''
);

CREATE INDEX site_slug_index ON sites(slug);

-- Users & Roles

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  organization_id INTEGER REFERENCES organizations(id),
  user_guid VARCHAR(64) UNIQUE NOT NULL,
  email VARCHAR(128) UNIQUE NOT NULL,
  password_digest VARCHAR(128) NOT NULL
);
CREATE INDEX user_guid_index ON users(user_guid);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    name INTEGER NOT NULL
);
CREATE INDEX roles_users_index ON roles(user_id);
CREATE INDEX roles_org_users_index ON roles(user_id, org_id);
CREATE UNIQUE INDEX roles_unique_index ON roles(user_id, org_id, name);

-- Site Coordinators

CREATE TABLE site_coordinators (
  id SERIAL PRIMARY KEY,
  site_id INTEGER REFERENCES sites(id),
  user_id INTEGER REFERENCES users(id)
);

CREATE INDEX sites_coordinators_index ON site_coordinators(site_id);
CREATE INDEX users_coordinated_sites_index ON site_coordinators(user_id);

-- Schedules

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dotw_type') THEN
        CREATE TYPE dotw_type AS ENUM('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday');
    END IF;
END$$;
CREATE TABLE daily_schedules (
  id SERIAL PRIMARY KEY,
  site_id INTEGER REFERENCES sites(id),
  dotw_default dotw_type, -- if this column is null, it must be an override. That means that the override_date column must be set
  override_date DATE, -- either this or dotw_default must be set. The other must be null.
  open_time VARCHAR(6) NOT NULL,
  close_time VARCHAR(6) NOT NULL,
  is_open BOOLEAN NOT NULL
);
CREATE INDEX sites_schedules_index ON daily_schedules(site_id);

