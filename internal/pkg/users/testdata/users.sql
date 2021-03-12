INSERT INTO users (id, user_guid, email, password_digest) VALUES
    -- password: password
    (1, 'kit', 'kit@example.org', '$2a$04$tGTAu2Rit8j6QAjyVeshg.rrX3rDulbXxsErP3eEOMfla1/g//p6C')
    , (2, 'user2', 'user2@example.org', '$2a$04$tGTAu2Rit8j6QAjyVeshg.rrX3rDulbXxsErP3eEOMfla1/g//p6C')
    , (3, 'user3', 'user3@example.org', '$2a$04$tGTAu2Rit8j6QAjyVeshg.rrX3rDulbXxsErP3eEOMfla1/g//p6C')
    , (4, 'user4', 'user4@example.org', '$2a$04$tGTAu2Rit8j6QAjyVeshg.rrX3rDulbXxsErP3eEOMfla1/g//p6C')
;

INSERT INTO organizations (id, name, slug, authcode) VALUES
    (1, 'testorg1', 'testorg1', 'testorg1')
    , (2, 'testorg2', 'testorg2', 'testorg2')
    , (3, 'testorg3', 'testorg3', 'testorg3')
;

INSERT INTO roles (id, org_id, user_id, name) VALUES
    (1, 1, 1, 1) -- kit, testorg1, OrgAdmin
    , (2, 2, 1, 1) -- kit, testorg2, OrgAdmin
    , (3, 1, 2, 2) -- user2, testorg1, Volunteer
    , (4, 2, 3, 2) -- user3, testorg2, Volunteer
    , (5, 3, 4, 2) -- user4, testorg3, Volunteer
;