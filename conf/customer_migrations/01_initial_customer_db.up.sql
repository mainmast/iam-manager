
CREATE TABLE IF NOT EXISTS users(

    uuid TEXT NOT NULL PRIMARY KEY,
    user_login TEXT UNIQUE,
    is_root boolean,
    user_type TEXT DEFAULT 'platform', -- Can be platform or api, platform for dashboard/admin users , api for users which call apis
    status TEXT DEFAULT 'active', -- Can be active, inactive, bolcked and deleted
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts(

    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    is_root boolean,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    custom_data JSONB

);

CREATE TABLE IF NOT EXISTS users_in_accounts (

 user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE, 
 account_uuid TEXT NOT NULL,
 PRIMARY KEY (user_uuid, account_uuid)
);

CREATE INDEX idx_uia_uuid ON users_in_accounts(user_uuid);

CREATE TABLE IF NOT EXISTS roles (

    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    custom_data JSONB
);

/*EXPERIMENTAL*/
CREATE TABLE IF NOT EXISTS resources (

    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS permissions (

    uuid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    resource TEXT NOT NULL REFERENCES resources(uuid) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    custom_data JSONB
);

CREATE INDEX idx_per_nm  ON permissions(name);

CREATE TABLE IF NOT EXISTS permission_in_role (
 role_uuid TEXT NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
 permission_uuid TEXT NOT NULL REFERENCES permissions(uuid) ON DELETE CASCADE
);

CREATE TABLE users_with_role (
 user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE, 
 role_uuid TEXT NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
 account_uuid TEXT NOT NULL REFERENCES accounts(uuid) ON DELETE CASCADE,
 PRIMARY KEY (user_uuid,role_uuid,account_uuid)
);

CREATE TABLE IF NOT EXISTS access_control (

    uuid TEXT NOT NULL PRIMARY KEY,
    access_key TEXT NOT NULL UNIQUE,
    secret_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used TIMESTAMPTZ,
    status TEXT DEFAULT 'active', -- Can be active, expired, revoked, renewed
    user_uuid TEXT NOT NULL REFERENCES users(uuid) ON DELETE CASCADE, 
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_usr_uuid  ON access_control(user_uuid);
CREATE INDEX idx_access_key ON access_control(access_key);
CREATE INDEX idx_secret_key ON access_control(secret_key);


CREATE OR REPLACE FUNCTION create_defaul_user (db_schema TEXT, account_name TEXT, usr_login TEXT, user_api_login TEXT,
access_ctrl TEXT, access_key TEXT)
RETURNS TABLE (
	usr_platform_uuid TEXT,
	usr_api_uuid TEXT,
	acc_uuid TEXT,
	access_uuid TEXT
    )
AS $$
DECLARE
    role_uuid TEXT;
    roletable TEXT;
BEGIN
	
	usr_platform_uuid:= uuid_generate_v4();
	usr_api_uuid := uuid_generate_v4();
	acc_uuid := uuid_generate_v4();
	access_uuid := uuid_generate_v4();


    --Creating the deafult account
    EXECUTE FORMAT('INSERT INTO %s.accounts(uuid, name, is_root) VALUES (%L,%L,%L);', db_schema, acc_uuid, account_name, 'true');

    -- Creating the default Platform user
     EXECUTE FORMAT('INSERT INTO %s.users(uuid, user_login, is_root) VALUES (%L,%L,%L);', db_schema, usr_platform_uuid, usr_login, 'true');

    -- Creating the default API user
    EXECUTE FORMAT('INSERT INTO %s.users(uuid, user_login, is_root, user_type) VALUES (%L,%L,%L,%L);', db_schema, usr_api_uuid, user_api_login, 'true', 'api');

    -- Crate the access keys for API user
    EXECUTE FORMAT('INSERT INTO %s.access_control(uuid, access_key, secret_key, user_uuid, expires_at) VALUES (%L,%L,%L,%L,%L);', db_schema, access_uuid, access_ctrl, access_key, usr_api_uuid, NOW() + INTERVAL '6 month');
	
    -- associate platform user in the account
    EXECUTE FORMAT('INSERT INTO %s.users_in_accounts(user_uuid, account_uuid) VALUES (%L,%L);', db_schema, usr_platform_uuid, acc_uuid);

     -- associate api user in the account
    EXECUTE FORMAT('INSERT INTO %s.users_in_accounts(user_uuid, account_uuid) VALUES (%L,%L);', db_schema, usr_api_uuid, acc_uuid);

	
    EXECUTE FORMAT('SELECT uuid FROM %s.roles WHERE name = %L', db_schema,'root') INTO role_uuid;
	
	IF role_uuid IS NULL THEN
        EXECUTE FORMAT('INSERT INTO %s.roles(uuid, name, description) VALUES (%L,%L,%L) returning uuid;', db_schema, uuid_generate_v4(), 'root', 'root / super admin user') INTO role_uuid;
    END IF;
		
	-- associate api user with the role
    EXECUTE FORMAT('INSERT INTO %s.users_with_role(user_uuid, role_uuid, account_uuid) VALUES (%L,%L,%L);', db_schema, usr_api_uuid, role_uuid, acc_uuid);

    -- associate api user with the role
    EXECUTE FORMAT('INSERT INTO %s.users_with_role(user_uuid, role_uuid, account_uuid) VALUES (%L,%L,%L);', db_schema, usr_platform_uuid, role_uuid, acc_uuid);
	
   RETURN NEXT;
END;
$$ LANGUAGE PLpgSQL;