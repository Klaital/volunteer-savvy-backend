package organizations

const createOrganizationSql = `INSERT INTO organizations (name) VALUES (:name)`
const updateOrganizationSql = `UPDATE organizations SET name=:name WHERE id=:id LIMIT 1`
const deleteOrganizationNullFkeysSql = `
	UPDATE sites SET organization_id=0 WHERE organization_id=:id; 
	UPDATE users SET organization_id=0 WHERE organization_id=:id; 
	DELETE FROM organizations WHERE id=:id LIMIT 1
`
const listOrganizationsSql = `SELECT id, name FROM organizations`
const describeOrganizationSql = `SELECT id, name FROM organizations WHERE id=?`
