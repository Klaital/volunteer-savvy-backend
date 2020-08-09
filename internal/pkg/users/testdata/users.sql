INSERT INTO users (id, user_guid, email, password_digest) VALUES
    -- password: password
    (1, 'kit', 'kit@example.org', '$2a$04$tGTAu2Rit8j6QAjyVeshg.rrX3rDulbXxsErP3eEOMfla1/g//p6C');

INSERT INTO organizations (id, name, slug, authcode) VALUES
    (1, 'testorg1', 'testorg1', 'testorg1');

INSERT INTO roles (id, org_id, user_id, name) VALUES
    (1, 1, 1, 1); -- kit, testorg1, OrgAdmin