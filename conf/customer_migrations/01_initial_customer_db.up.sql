CREATE TABLE IF NOT EXISTS platforms (
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- see iam-models/pkg/platform/enum.go
    whitelist_domains TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    custom_data JSONB
);


CREATE TABLE IF NOT EXISTS users(
    uuid TEXT NOT NULL PRIMARY KEY,
    user_login TEXT UNIQUE,
    user_secret TEXT NOT NULL,
    is_root boolean,
    console_access boolean NOT NULL DEFAULT false, -- Access to any GUI consoles
    api_access boolean NOT NULL DEFAULT false, -- Direct access to an API
    status TEXT DEFAULT 'active', -- see iam-models/pkg/user/enum.go
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    updated_at TIMESTAMPTZ
    custom_data JSONB
);

CREATE INDEX idx_users_user_login ON iam_users(user_login);
CREATE INDEX idx_users_user_secret ON iam_users(user_secret);

CREATE TABLE IF NOT EXISTS user_platforms (
    user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    platform_uuid TEXT NOT NULL REFERENCES platforms(uuid) ON DELETE CASCADE,
    PRIMARY KEY(user_uuid, platform_uuid)
);

CREATE TABLE IF NOT EXISTS accounts(
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- see iam-models/pkg/account/enum.go
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    custom_data JSONB,
    platform_uuid TEXT NOT NULL REFERENCES platforms(uuid) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS account_users (
    account_uuid TEXT NOT NULL REFERENCES accounts(uuid) ON DELETE CASCADE,
    user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE, 
    PRIMARY KEY (user_uuid, account_uuid)
);

CREATE INDEX idx_account_users_user_uuid ON account_users(user_uuid);

CREATE TABLE IF NOT EXISTS roles (
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    custom_data JSONB
);

CREATE TABLE IF NOT EXISTS resources (

    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    custom_data JSONB
);

CREATE TABLE IF NOT EXISTS permissions (
    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    resource TEXT NOT NULL REFERENCES resources(uuid) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    custom_data JSONB
);

CREATE INDEX idx_permissions_name  ON permissions(name);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_uuid TEXT NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
    permission_uuid TEXT NOT NULL REFERENCES permissions(uuid) ON DELETE CASCADE
);

CREATE TABLE account_user_roles (
    user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE, 
    role_uuid TEXT NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
    account_uuid TEXT NOT NULL REFERENCES accounts(uuid) ON DELETE CASCADE,
    PRIMARY KEY (
        user_uuid,
        role_uuid,
        account_uuid)
);

CREATE TABLE IF NOT EXISTS access_keys (
    uuid TEXT NOT NULL PRIMARY KEY,
    access_key TEXT NOT NULL UNIQUE,
    secret_key TEXT NOT NULL UNIQUE,
    status TEXT DEFAULT 'active', -- Can be active, expired, revoked, renewed
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
    user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE INDEX idx_access_keys_user_uuid  ON access_keys(user_uuid);
CREATE INDEX idx_access_keys_access_key ON access_keys(access_key);
CREATE INDEX idx_access_keys_secret_key ON access_keys(secret_key);

CREATE TABLE IF NOT EXISTS access_key_permissions (
    access_key_uuid TEXT NOT NULL REFERENCES access_keys(uuid) ON DELETE CASCADE,
    permission_uuid TEXT NOT NULL REFERENCES permissions(uuid) ON DELETE CASCADE
);