package orm

/*
	Queries used for the System IAM side of the application. This is the base company for which the service
	is being offered.
*/

// IAM Organisation
const (

	//CreateOrganisationQuery Creates an IAM organisation.
	CreateOrganisationQuery string = "SELECT * FROM public.sp_create_organisation(name := $1, db_schema := $2, status := $3, custom_data := $4);"

	//CreateDefaultOrganisationQuery Creates an IAM default user for the organisation.
	CreateDefaultOrganisationQuery string = "SELECT * FROM public.sp_create_organisation(uuid := $1, name := $2, whitelist_domains := $3, db_schema := $4, custom_data := $5, user_login := $6, user_secret := $7, user_custom_data := $8);"

	//CheckOrganisationByNameQuery Checks to see if an organisation name is already in use.
	CheckOrganisationByNameQuery string = "SELECT * FROM public.ss_organisation_exists(name := $1);"

	//GetSchemaByUUIDQuery Gets the schema name for an organisation from the UUID
	GetSchemaByUUIDQuery string = "SELECT db_schema FROM public.ss_get_organisation(uuid := $1);"

	//GetOrganisationBySchemaQuery Gets the name of an organisation from the schema name.
	GetOrganisationBySchemaQuery string = "SELECT name FROM public.ss_get_organisation(db_schema := $1);"
)

// IAM Users
const (
	//GetIamUserLoginByUUIDQuery Gets the IAM user login from the UUID
	GetIamUserLoginByUUIDQuery string = "SELECT user_login FROM public.ss_get_iam_user(uuid := $1);"

	//CreateIamUserQuery Creates an IAM user
	CreateIamUserQuery string = "INSERT INTO public.iam_users(uuid, user_login, user_secret, custom_data) VALUES ($1, $2, $3, $4);"

	//DeleteIamUserQuery Deletes an IAM user.
	DeleteIamUserQuery string = "DELETE FROM public.iam_users WHERE uuid = $1;"
)

const (
	//AssociateUserWithOrganisationQuery Associates an IAM user with an IAM organisation
	AssociateUserWithOrganisationQuery string = "INSERT INTO public.iam_users_organisations (user_uuid, organisation_uuid) VALUES ($1,$2)"
)

/*
	Queries used to create an organisation's data. This is their private part of the DB (separated by schemas)
*/

// Platforms
const (

	//CreatePlatformQuery ...
	CreatePlatformQuery string = "SELECT * FROM #*.sp_create_platform(name := $1, status := $2, custom_data := $3);"

	//ChangePlatformStatusQuery ...
	ChangePlatformStatusQuery string = "SELECT * FROM #*.ss_get_platform_status(uuid := $1, effective_at := $2);"

	//UpdatePlatformQuery ...
	UpdatePlatformQuery string = "SELECT * FROM #*.sp_update_platform(uuid := $1, name := $2, status := $3, custom_data := $4);"

	//DeletePlatformQuery ...
	DeletePlatformQuery string = "SELECT * FROM #*.sp_delete_platform(uuid := $1, permanently_delete := $2);"

	//GetPlatformByUUIDQuery ...
	GetPlatformByUUIDQuery string = "SELECT * FROM #*.ss_get_platform(uuid := $1);"

	//GetPlatformByNameQuery ...
	GetPlatformByNameQuery string = "SELECT * FROM #*.ss_get_platform(name := $1);"
)

// Accounts
const (
	//CreateAccountQuery Creates an account
	CreateAccountQuery string = "INSERT INTO #*.accounts (uuid, name, custom_data) VALUES ($1, $2, $3);"

	//UpdateAccountQuery Updates an account
	UpdateAccountQuery string = "UPDATE #*.accounts (uuid, name, custom_data) VALUES ($1, $2, $3);"

	//AssociateUserToAccountQuery Associates an organisation user with an account
	AssociateUserToAccountQuery string = "INSERT INTO #*.users_in_accounts (user_uuid, account_uuid) VALUES ($1, $2)"

	//UserInDefaultAccountQuery ...
	UserInDefaultAccountQuery string = "SELECT (EXISTS (SELECT FROM #*.accounts acc INNER JOIN #*.users_in_accounts uia ON uia.account_uuid = acc.uuid WHERE uia.user_uuid = $1 AND acc.is_root = 'true'))::bool as is_user_in_org;"

	//DefaultOrgAccountQuery
	DefaultOrgAccountQuery string = "SELECT acc.uuid FROM #*.accounts acc WHERE acc.is_root = 'true'"
)

// Organisation Users
const (

	//CreateUserQuery ...
	CreateUserQuery = "INSERT INTO #*.users (uuid, user_login, is_root, user_type) VALUES ($1, $2, $3, $4);"

	//CreateUserInOrg ...
	CreateUserInAccQuery = "INSERT INTO #*.users_in_accounts (user_uuid, account_uuid) VALUES($1, $2);"

	//GetUserAccesByUUIDQuery ...
	GetUserAccesByUUIDQuery = "SELECT usr.user_type, acc.access_key, acc.secret_key, acc.status FROM #*.users usr INNER JOIN #*.access_control acc ON acc.user_uuid=usr.uuid WHERE usr.uuid = $1;"

	// CreateDefaultUserQuery creates a default user, with an associated account if provided.
	CreateDefaultUserQuery = "SELECT user_uuid, account_uuid FROM #*.create_root_user($1,$2,$3);"
)

//AccessControl

const (
	//GenerateAccessQuery ...
	GenerateAccessQuery = "INSERT INTO #*.access_control (uuid, access_key, secret_key, user_uuid, expires_at) VALUES ($1, $2, $3, $4, $5)"
)
