CREATE TABLE IF NOT EXISTS iam_organisations (
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    status TEXT DEFAULT 'active', -- see iam-models/pkg/organisation/enum.go
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    db_schema TEXT NOT NULL,
    custom_data JSONB
);

CREATE INDEX idx_iam_organisations_status ON iam_organisations(status);

CREATE TABLE IF NOT EXISTS iam_users(
    uuid TEXT NOT NULL PRIMARY KEY,
    user_login TEXT NOT NULL UNIQUE,
    user_secret TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    is_root BOOLEAN NOT NULL DEFAULT FALSE,
    console_access BOOLEAN NOT NULL DEFAULT TRUE,
    api_access BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    custom_data JSONB
);

CREATE INDEX idx_iam_users_user_login ON iam_users(user_login);
CREATE INDEX idx_iam_users_user_secret ON iam_users(user_secret);

CREATE TABLE IF NOT EXISTS iam_users_organisations(

    user_uuid TEXT NOT NULL REFERENCES iam_users(uuid) ON DELETE CASCADE,
    organisation_uuid TEXT NOT NULL REFERENCES iam_organisations(uuid) ON DELETE CASCADE,
    status TEXT DEFAULT 'active', -- see iam-models/pkg/organisation/enum.go
    PRIMARY KEY (user_uuid, organisation_uuid)
);

CREATE INDEX idx_iam_users_organisations_user_uuid ON iam_users_organisations(user_uuid);
CREATE INDEX idx_iam_users_organisations_status ON iam_users_organisations(status);

--create the extension to generate the uuid iside the functions 
CREATE EXTENSION "uuid-ossp";
CREATE EXTENSION "pgcrypto";