CREATE OR REPLACE FUNCTION sp_create_platform (
    IN p_db_schema TEXT,
    IN p_platform_name TEXT,
    IN p_status TEXT,
    IN p_whitelist_domains TEXT[],
    IN p_custom_data JSONB
    IN p_root_user_login TEXT,
    IN p_root_user_secret TEXT,
    IN p_root_user_custom_data JSONB)
RETURNS TABLE (
    r_status TEXT,
    r_status_message TEXT,
	r_platform_uuid TEXT,
    r_platform_name TEXT,
    r_status TEXT,
    r_whitelist_domains TEXT[],
    r_custom_data JSONB)
AS $$
DECLARE
    r_platform_uuid TEXT;
    role_uuid TEXT;
BEGIN	
	r_platform_uuid := uuid_generate_v4();

    -- Create the first organisation user.
    EXECUTE FORMAT('INSERT INTO %s.users(uuid, user_login, user_secret, is_root, platform_access, api_access) VALUES (%L,%L,%L,%L,%L,%L);', 
        db_schema, 
        user_uuid, 
        user_login, 
        user_secret,
        'true',
        'true',
        'true'        
    );

    -- Create the first account if supplied.
    IF account_name != "" THEN
        account_uuid := uuid_generate_v4();
        
        EXECUTE FORMAT('INSERT INTO %s.accounts(uuid, name, is_root) VALUES (%L,%L,%L);', 
            db_schema, 
            account_uuid, 
            account_name, 
            'true'
        );    

        -- associate user with the account
        EXECUTE FORMAT('INSERT INTO %s.account_users(user_uuid, account_uuid) VALUES (%L,%L);', db_schema, user_uuid, account_uuid);

    END IF;
		
    EXECUTE FORMAT('SELECT uuid FROM %s.roles WHERE name = %L', 
        db_schema,'root') 
    INTO role_uuid;
	
	IF role_uuid IS NULL THEN
        EXECUTE FORMAT('INSERT INTO %s.roles(uuid, name, description) VALUES (%L,%L,%L) returning uuid;', 
            db_schema, 
            uuid_generate_v4(), 
            'root', 
            'root / super admin user') 
        INTO role_uuid;
    END IF;
		
	-- Associate user with the root role.
    EXECUTE FORMAT('INSERT INTO %s.account_user_roles(user_uuid, role_uuid, account_uuid) VALUES (%L,%L,%L);', 
        db_schema, 
        user_uuid, 
        role_uuid, 
        account_uuid);
	
   RETURN NEXT;
END;
$$ LANGUAGE PLpgSQL;