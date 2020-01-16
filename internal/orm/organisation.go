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

	org "github.com/mainmast/iam-models/pkg/organisation"
)

// Create an organisation entity.
func CreateOrganisation(req *org.CreateOrgRQ) (*org.CreateOrgRS, error) {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting to database", err)
		return nil, fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var allow int

	// Check to see if the organisation name already exists.
	if err := conn.QueryRowContext(context.Background(), CheckOrganisationByNameQuery, req.Organisation.Name).Scan(&allow); err != nil {

		fmt.Println("Error validating organisation name", err)
		return nil, errors.New("Error validating organisation name")
	}

	if allow > 0 {

		fmt.Println("Organisation name", req.Organisation.Name, "already exists")
		return nil, errors.New("DUPLICATED: Organisation name already exists. Please choose another one.")
	}

	schema := util.NormaliseOrganisationName(req.Organisation.Name, true)

	if rs := util.SetUpCustomer("upgrade", schema); !rs {

		fmt.Println("Error setting up customer database")
		util.SetUpCustomer("downgrade", schema)
		return nil, errors.New("Error setting up customer database")
	}

	organisationCustomData, err := json.Marshal(req.Organisation.CustomData)

	if err != nil {
		fmt.Println("Error reading organisation custom data", err)
		return nil, errors.New("Invalid organisation custom data object")
	}

	userCustomData, err := json.Marshal(req.User.CustomData)

	if err != nil {

		fmt.Println("Error reading user custom data", err)
		return nil, errors.New("Invalid user custom data object")
	}

	res := &org.CreateOrgRS{}

	res.OrganisationName = req.Organisation.Name
	res.UserLogin = req.User.UserLogin
	userSecret, err := util.SecretHash(req.User.UserSecret)

	if err != nil {
		return nil, errors.New("Error generating password hash")
	}

	if err := conn.QueryRowContext(context.Background(), util.ParseQuerySchema(schema, CreateDefaultUserQuery), &schema, &res.AccountName, &req.User.UserLogin, &req.User.UserSecret).Scan(&res.UserUUID, &res.AccountUUID); err != nil {

		fmt.Println("Error creating the root user for the organisation", err)
		return nil, errors.New("Error creating root user")
	}

	if err := conn.QueryRowContext(context.Background(), CreateIamOrgAndUsrQuery, &req.Organisation.Name, &req.User.UserLogin, &userSecret, &schema, &organisationCustomData,
		&usrCustom, &res.UserPlatformUUID, pq.Array(&req.Organisation.WhitelistDomains)).Scan(&res.OrganisationUUID); err != nil {

		fmt.Println("error creating organization, Rollbacking everything", err)
		util.SetUpCustomer("downgrade", schema)
		return nil, errors.New("error creating organization")
	}

	fmt.Println("Organisation created and customer set up sucessfully")
	return res, nil
}

// Associate an IAM user with an organisation
func AssociateUserWithOrganisation(userUUID string, organisationUUID string) error {

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
