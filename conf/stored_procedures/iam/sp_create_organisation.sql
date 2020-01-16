CREATE OR REPLACE FUNCTION sp_create_organisation (
    IN p_name TEXT, 
    IN p_db_schema TEXT, 
    IN p_status TEXT,
    IN p_user_login TEXT,
    IN p_user_secret TEXT,
    IN p_custom_data JSONB DEFAULT NULL,
    IN p_user_custom_data JSONB DEFAULT NULL,
    IN p_organisation_uuid TEXT DEFAULT NULL,
    IN p_user_uuid TEXT DEFAULT NULL
)
RETURNS TABLE (
    r_status TEXT,
    r_status_message TEXT,
	r_organisation_uuid TEXT,
	r_user_uuid TEXT)
AS $$
-------------------------------------------------------------------------------
-- This function creates iam organisations, and the associated root user for
-- the organisation.
-------------------------------------------------------------------------------
DECLARE
    l_organisation_uuid TEXT;
    l_status TEXT;

BEGIN
    l_organisation_uuid := COALESCE(p_organisation_uuid, gen_random_uuid()::text);

    --Creating the deafult account
    BEGIN
        INSERT INTO iam_organisations (
            uuid, 
            name, 
            db_schema, 
            custom_data)
        VALUES (
            l_organisation_uuid, 
            p_name, 
            p_db_schema, 
            p_custom_data);

        -- Creating the Platform user
        SELECT 
            sp_create_iam_user.r_status, 
            sp_create_iam_user.r_status_message, 
            sp_create_iam_user.r_uuid
        INTO
            r_status,
            r_status_message,
            r_user_uuid
        FROM 
            sp_create_iam_user(
                p_organisation_uuid := l_organisation_uuid, 
                p_user_login := p_user_login, 
                p_user_secret := p_user_secret, 
                p_user_status := 'active', 
                p_custom_data := p_user_custom_data, 
                p_is_root := true, 
                p_console_access := true, 
                p_api_access := false);
        
        IF r_status != 'success' THEN
            RAISE EXCEPTION USING
                errcode='USRER',
                message=r_status_message;
        END IF;

        -- Completed successfully, return the uuid.
        r_organisation_uuid := l_organisation_uuid;

        r_status := 'success';

    EXCEPTION 
        WHEN SQLSTATE 'USRER' THEN
            r_status := 'failure';
            r_status_message := FORMAT('Failed to create root user: %L ', r_status_message);

        WHEN unique_violation THEN
            r_status := 'failure';
            r_status_message := FORMAT('Organisation name %L already exists', p_name);

        WHEN others THEN
            r_status := 'exception';
            GET STACKED DIAGNOSTICS r_status_message = MESSAGE_TEXT;
    END;

    RETURN NEXT;
    RETURN;
END;
$$ LANGUAGE PLpgSQL;

-- select * from sp_create_organisation ( p_name := 'mainmast', p_db_schema := 'mainmast_9xhs8', p_status := 'active', p_user_login := 'trentm', p_user_secret := 'XXXX', p_custom_data := NULL, p_user_custom_data := NULL );
