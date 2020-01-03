CREATE TABLE IF NOT EXISTS iam_organisations(
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    status TEXT DEFAULT 'active', -- Can be pending, active, disabled, blocked
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    whitelist_domains TEXT[],
    db_schema TEXT NOT NULL,
    custom_data JSONB
);

CREATE INDEX idx_org_status ON iam_organisations(status);

-- Used for platform users
CREATE TABLE IF NOT EXISTS iam_users(
    uuid TEXT NOT NULL PRIMARY KEY,
    user_login TEXT NOT NULL UNIQUE,
    user_secret TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    custom_data JSONB
);

CREATE INDEX idx_usr_lgin ON iam_users(user_login);
CREATE INDEX idx_usr_scret ON iam_users(user_secret);

CREATE TABLE IF NOT EXISTS iam_users_organisations(

    user_uuid TEXT NOT NULL REFERENCES iam_users(uuid) ON DELETE CASCADE,
    organisation_uuid TEXT NOT NULL REFERENCES iam_organisations(uuid) ON DELETE CASCADE,
    status TEXT DEFAULT 'active', -- Can be active, inactive, bolcked and deleted
    PRIMARY KEY (user_uuid, organisation_uuid)
);

CREATE INDEX idx_usrorg_lgin ON iam_users_organisations(user_uuid);
CREATE INDEX idx_usrorg_status ON iam_users_organisations(status);

--create the extension to generate the uuid iside the functions 
CREATE EXTENSION "uuid-ossp";

CREATE OR REPLACE FUNCTION create_default_org (org_name TEXT, usr_login TEXT, user_secret TEXT, db_schema TEXT, 
 org_custom_data JSONB, usr_custom_data JSONB, user_platform_uuid TEXT, whitelist_domains TEXT[])RETURNS TEXT AS $$
DECLARE
org_uuid TEXT;
BEGIN

    org_uuid := uuid_generate_v4();
    --Creating the deafult account
    INSERT INTO iam_organisations (uuid, name, whitelist_domains, db_schema, custom_data)
    VALUES (org_uuid, org_name, whitelist_domains, db_schema, org_custom_data);

    -- Creating the Platform user
    INSERT INTO iam_users (uuid, user_login, user_secret, custom_data)
    VALUES (user_platform_uuid, usr_login, user_secret, usr_custom_data);

    --  creating user in the organisation
    INSERT INTO iam_users_organisations (user_uuid, organisation_uuid)
    VALUES (user_platform_uuid, org_uuid);

   RETURN org_uuid;
END;
$$ LANGUAGE PLpgSQL;