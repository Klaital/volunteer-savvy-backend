package organizations

const createOrganizationSql = `
INSERT INTO organizations 
		(name, slug, authcode, contact_user_id, lat, lon) 
	VALUES 
		(:name, :slug, :authcode, :contact_user_id, :lat, :lon)`
const updateOrganizationSql = `
UPDATE organizations 
SET 
	name=:name,
	slug=:slug,
	authcode=:authcode,
	contact_user_id=:contact_user_id,
	lat=:lat,
	lon=:lon
WHERE id=:id`
const deleteOrganizationNullFkeysSql = `
	UPDATE sites SET organization_id=0 WHERE organization_id=:id; 
	UPDATE users SET organization_id=0 WHERE organization_id=:id; 
	DELETE FROM organizations WHERE id=:id LIMIT 1
`
const listOrganizationsSql = `SELECT id, name, slug, authcode, contact_user_id, lat, lon FROM organizations`
const describeOrganizationSql = `SELECT id, name, slug, authcode, contact_user_id, lat, lon FROM organizations WHERE id=?`
