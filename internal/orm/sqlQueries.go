package orm

//## Public/default schema

//iam_org
const (

	//CreateIamOrgQuery ...
	CreateIamOrgQuery = "INSERT INTO public.iam_organisations (uuid, name, whitelist_domains, db_schema, custom_data) VALUES ($1,$2,$3,$4,$5);"
	//CreateIamOrgAndUsrQuery ...
	CreateIamOrgAndUsrQuery = "SELECT * FROM public.create_default_org($1,$2,$3,$4,$5,$6,$7,$8)"
	//ExistsOrgQuery ORG
	NotExistsOrgQuery = "SELECT COUNT(NAME) as count_org FROM public.iam_organisations WHERE name = $1"
	// GetSchemaQuery
	GetSchemaQuery = "SELECT db_schema FROM public.iam_organisations WHERE uuid = $1"
	//GetOrgNameQuery ...
	GetIamOrgQuery = "SELECT name FROM public.iam_organisations WHERE db_schema = $1"
)

//iam_users_and_orgs
const (
	//CreateIamUserQuery ...
	CreateIamUserOrgQuery = "INSERT INTO public.iam_users_organisations (user_uuid, organisation_uuid) VALUES ($1,$2)"
)

//iam_users
const (
	//CreateIamUserQuery ...
	GetUserSchemaByQuery = "SELECT user_login FROM public.iam_users WHERE uuid = $1"
	//CreateUserQuery ...
	CreateIamUserQuery = "INSERT INTO public.iam_users(uuid, user_login, user_secret, custom_data) VALUES ($1, $2, $3, $4);"
	//DeleteUsrQuery
	DeleteUsrQuery = "DELETE FROM public.iam_users WHERE uuid = $1;"
)

//##Customer Schemas

//Account
const (
	//CreateOrganisationQuery ...
	CreateAccountQuery = "INSERT INTO #*.accounts (uuid, name, is_root, custom_data) VALUES ($1, $2, $3, $4);"
	//AssociateUserToAccountQuery ...
	AssociateUserToAccountQuery = "INSERT INTO #*.users_in_accounts (user_uuid, account_uuid) VALUES ($1, $2)"
	//UserInDefaultAccountQuery ...
	UserInDefaultAccountQuery = "SELECT (EXISTS (SELECT FROM #*.accounts acc INNER JOIN #*.users_in_accounts uia ON uia.account_uuid = acc.uuid WHERE uia.user_uuid = $1 AND acc.is_root = 'true'))::bool as is_user_in_org;"
	//DefaultOrgAccountWuery
	DefaultOrgAccountQuery = "SELECT acc.uuid FROM #*.accounts acc WHERE acc.is_root = 'true'"
)

//User
const (

	//CreateUserQuery ...
	CreateUserQuery = "INSERT INTO #*.users (uuid, user_login, is_root, user_type) VALUES ($1, $2, $3, $4);"
	//CreateUserInOrg ...
	CreateUserInAccQuery = "INSERT INTO #*.users_in_accounts (user_uuid, account_uuid) VALUES($1, $2);"
	//GetUserAccesByUUIDQuery ...
	GetUserAccesByUUIDQuery = "SELECT usr.user_type, acc.access_key, acc.secret_key, acc.status FROM #*.users usr INNER JOIN #*.access_control acc ON acc.user_uuid=usr.uuid WHERE usr.uuid = $1;"
	//CreateOrganisationQuery creates the organisation and the default accoutn and platform , api root users
	CreateDefaultUsrQuery = "SELECT usr_platform_uuid, usr_api_uuid, acc_uuid, access_uuid FROM #*.create_defaul_user($1,$2,$3,$4,$5,$6);"
)

//AccessControl

const (
	//GenerateAccessQuery ...
	GenerateAccessQuery = "INSERT INTO #*.access_control (uuid, access_key, secret_key, user_uuid, expires_at) VALUES ($1, $2, $3, $4, $5)"
)
