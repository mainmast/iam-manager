package orm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	pq "github.com/lib/pq"

	"mainmast/iam-manager/internal/util"

	org "gitlab.com/mainmast/microservices/iam/iam-models.git/pkg/organisation"
)

//CreateOrganisation ...
func CreateOrganisation(req *org.CreateOrgRQ) (*org.CreateOrgRS, error) {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return nil, fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var allow int

	if err := conn.QueryRowContext(context.Background(), NotExistsOrgQuery, req.Organisation.Name).Scan(&allow); err != nil {

		fmt.Println("error validating org name", err)
		return nil, errors.New("error validating org name")
	}

	if allow > 0 {

		fmt.Println("org name", req.Organisation.Name, "alredy exists")
		return nil, errors.New("DUPLICATED: Org name already exist please choose another one")
	}

	schema := util.NormalizeName(req.Organisation.Name, true)

	if rs := util.SetUpCustomer("upgrade", schema); !rs {

		fmt.Println("Error setting up customer database")
		util.SetUpCustomer("downgrade", schema)
		return nil, errors.New("Error setting up customer database")
	}

	orgCustom, err := json.Marshal(req.Organisation.CustomData)

	if err != nil {

		fmt.Println("error reading organisation custom data", err)
		return nil, errors.New("invalid organisation custom data object")
	}

	usrCustom, err := json.Marshal(req.User.CustomData)

	if err != nil {

		fmt.Println("error reading user custom data", err)
		return nil, errors.New("invalid user custom data object")
	}

	res := &org.CreateOrgRS{}

	res.OrganisationName = req.Organisation.Name
	res.UserAPILogin = util.GenerateAPILogin(req.User.UserLogin)
	res.UserPlatformLogin = req.User.UserLogin
	res.AccountName = util.NormalizeName(req.Organisation.Name, false)

	userSecret, err := util.SecretHash(req.User.UserSecret)

	if err != nil {

		return nil, errors.New("error generating password hash")
	}

	if err := conn.QueryRowContext(context.Background(), util.ParseQuerySchema(schema, CreateDefaultUsrQuery), &schema, &res.AccountName, &req.User.UserLogin,
		&res.UserAPILogin, util.GenerateAccesKey(req.Organisation.Name), util.GenerateSecretKey()).Scan(&res.UserPlatformUUID, &res.UserAPIUUID, &res.AccountUUID, &res.AccessUUID); err != nil {

		fmt.Println("error creating default user", err)
		return nil, errors.New("error creating default user")
	}

	if err := conn.QueryRowContext(context.Background(), CreateIamOrgAndUsrQuery, &req.Organisation.Name, &req.User.UserLogin, &userSecret, &schema, &orgCustom,
		&usrCustom, &res.UserPlatformUUID, pq.Array(&req.Organisation.WhitelistDomains)).Scan(&res.OrganisationUUID); err != nil {

		fmt.Println("error creating organization, Rollbacking everything", err)
		util.SetUpCustomer("downgrade", schema)
		return nil, errors.New("error creating organization")
	}

	fmt.Println("Organisation created and customer set up sucessfully")
	return res, nil
}

//AssociateUserWithOrg ....
func AssociateUserWithOrg(userUUID string, organisationUUID string) error {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var defaulAcc string
	var userLogin string
	var dbSchema string

	if err := conn.QueryRowContext(context.Background(), GetUserSchemaByQuery, userUUID).Scan(&userLogin); err != nil {

		fmt.Println("Not user found", err)
		return errors.New("NOT FOUND: user not found")
	}

	if err := conn.QueryRowContext(context.Background(), GetSchemaQuery, organisationUUID).Scan(&dbSchema); err != nil {

		fmt.Println("Organisation not found", err)
		return errors.New("NOT FOUND: Organisation not found")
	}

	if err := conn.QueryRowContext(context.Background(), util.ParseQuerySchema(dbSchema, DefaultOrgAccountQuery)).Scan(&defaulAcc); err != nil {

		fmt.Println("Invalid Organisation uuid, not default account created yet", err)
		return errors.New("Invalid Organisation uuid, not defaul account createy yet")
	}

	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {

		fmt.Println("User transaction could not be created")
		return errors.New("User transaction could not be created")
	}

	if _, err := tx.ExecContext(context.Background(), CreateIamUserOrgQuery, userUUID, organisationUUID); err != nil {

		tx.Rollback()
		fmt.Println("Error associating user into the organisation", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: Error already linked to this organisation")

		}
		return errors.New("Error creating user")
	}

	if _, err := tx.ExecContext(context.Background(), util.ParseQuerySchema(dbSchema, CreateUserQuery), userUUID, userLogin, false, "platform"); err != nil {

		tx.Rollback()
		fmt.Println("Error creating user in customer DB", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: user alredy present in this organisation")
		}

		return errors.New("Error associating user to the organisation")
	}

	if _, err := tx.ExecContext(context.Background(), util.ParseQuerySchema(dbSchema, CreateUserInAccQuery), userUUID, defaulAcc); err != nil {

		tx.Rollback()
		fmt.Println("Error associating user to default account", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: user alredy linked to default account")
		}

		return errors.New("Error associating user to the account")
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Error associating user", err)
		return errors.New("Error associating user")
	}

	fmt.Println("User associated with organisation sucessfuly")
	return nil
}
