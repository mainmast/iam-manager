CREATE OR REPLACE FUNCTION sp_create_iam_user (
    IN p_organisation_uuid TEXT,
    IN p_user_login TEXT,
    IN p_user_secret TEXT,
    IN p_user_status TEXT,
    IN p_custom_data JSONB DEFAULT NULL,
    IN p_is_root BOOLEAN DEFAULT FALSE,
    IN p_console_access BOOLEAN DEFAULT TRUE,
    IN p_api_access BOOLEAN DEFAULT FALSE,
    IN p_user_uuid TEXT DEFAULT NULL
) RETURNS TABLE (
    r_status TEXT,
    r_status_message TEXT,
    r_uuid TEXT,
    r_user_login TEXT,
    r_user_status TEXT,
    r_is_root BOOLEAN,
    r_console_access BOOLEAN,
    r_api_access BOOLEAN,
    r_custom_data TEXT,
    r_organisation_uuid TEXT)
AS $$
-------------------------------------------------------------------------------
-- This function creates iam users and associates them with organisations.
-------------------------------------------------------------------------------
DECLARE
    l_organisation_count INTEGER;
BEGIN

    r_uuid := NULL;
    r_user_login := p_user_login;
    r_organisation_uuid := p_organisation_uuid;
    r_custom_data := p_custom_data;
    r_user_status := p_user_status;
    r_is_root := p_is_root;
    r_console_access := p_console_access;
    r_api_access := p_api_access;

    BEGIN
        -- If a UUID is provided, use it - otherwise generate on. The calling 
        -- application should always send this where possible, as this imple
        r_uuid := COALESCE(p_organisation_uuid, gen_random_uuid()::text);

        -- Create the IAM user.
        INSERT INTO iam_users (
            uuid, 
            user_login, 
            user_secret,
            status,
            is_root,
            console_access,
            api_access, 
            custom_data)
        VALUES (
            r_uuid, 
            p_user_login, 
            p_user_secret,
            p_user_status,
            p_is_root,
            p_console_access,
            p_api_access,
            p_custom_data
        );

        -- Assign the user to the supplied organisation
        SELECT
            COUNT(*)
        INTO
            l_organisation_count
        FROM
            iam_organisations
        WHERE
            uuid = p_organisation_uuid;
        
        IF l_organisation_count > 0 THEN
            INSERT INTO iam_users_organisations (
                user_uuid, 
                organisation_uuid)
            VALUES (
                r_uuid, 
                p_organisation_uuid
            );
        
        ELSE
            RAISE EXCEPTION USING
                errcode='NOORG',
                message=FORMAT('No organisation found with the supplied uuid: %L', p_organisation_uuid);

        END IF;
	
        RETURN NEXT;
    
    -- Catch any exceptions and return a better response.
    EXCEPTION 
        WHEN SQLSTATE 'NOORG' THEN
            GET STACKED DIAGNOSTICS r_status_message = MESSAGE_TEXT;
            r_status := 'detected_failure';            
            RETURN NEXT;

        WHEN OTHERS THEN
            r_status := 'exception';
            GET STACKED DIAGNOSTICS r_status_message = MESSAGE_TEXT;
            RETURN NEXT;
    
    END;

    RETURN;

END;
$$ LANGUAGE PLpgSQL;