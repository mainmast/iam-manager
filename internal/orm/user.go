package orm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	//postgres driver
	_ "github.com/lib/pq"

	"mainmast/iam-manager/internal/util"

	usr "github.com/mainmast/iam-models/pkg/user"
)

//CreateUser ...
func CreateUser(user *usr.User, orgUUID string) error {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {

		fmt.Println("User transaction could not be created")
		return errors.New("User transaction could not be created")
	}

	user.UUID = uuid.NewV4().String()
	user.IsRoot = false
	user.UserType = usr.Plaftorm

	if err := util.CheckUserType(user); err != nil {

		fmt.Println("Error checking user", err)
		return errors.New("Error checking user")
	}

	usrCustom, err := json.Marshal(user.CustomData)

	if err != nil {

		fmt.Println("error reading user custom data", err)
		return errors.New("invalid user custom data object")
	}

	if _, err := tx.ExecContext(context.Background(), CreateIamUserQuery, user.UUID, user.UserLogin, user.UserSecret, usrCustom); err != nil {

		tx.Rollback()
		fmt.Println("Error creating user", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: Error creating user")

		}
		return errors.New("Error creating user")
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Error creating user", err)
		return errors.New("Error creating user")
	}

	if strings.TrimSpace(orgUUID) != "" {

		if err := AssociateUserWithOrg(user.UUID, orgUUID); err != nil {

			tx.Rollback()
			if _, err := conn.ExecContext(context.Background(), DeleteUsrQuery, user.UUID); err != nil {

				fmt.Println("User created but not linked with any organisation, error rollbacking user", err)
			}
			fmt.Println("Error associating user with the org: ", err)
			return errors.New("Error associating user with the org")
		}

	}

	user.Status = "active"
	user.UserSecret = "REDACTED"
	fmt.Println("user created")
	return nil
}

//CreateUserAPI ...
func CreateUserAPI(user *usr.User, schema string) error {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var defaulAcc string

	if err := conn.QueryRowContext(context.Background(), util.ParseQuerySchema(schema, DefaultOrgAccountQuery)).Scan(&defaulAcc); err != nil {

		fmt.Println("User can not be created in this organisation", schema, err)
		return errors.New("Invalid Organisation, not default account created yet")
	}

	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {

		fmt.Println("User transaction could not be created")
		return errors.New("User transaction could not be created")
	}

	user.UUID = uuid.NewV4().String()
	user.IsRoot = false
	user.UserType = usr.API

	if _, err := tx.ExecContext(context.Background(), util.ParseQuerySchema(schema, CreateUserQuery), user.UUID, user.UserLogin, false, "api"); err != nil {

		tx.Rollback()
		fmt.Println("Error creating user", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: Error creating API user")

		}
		return errors.New("Error creating API user")
	}

	if _, err := tx.ExecContext(context.Background(), util.ParseQuerySchema(schema, CreateUserInAccQuery), user.UUID, defaulAcc); err != nil {

		tx.Rollback()
		fmt.Println("Error associating user with default account", err)

		if strings.HasPrefix(err.Error(), "pq: duplicate key") {

			return errors.New("DUPLICATED: Error creating API user")

		}
		return errors.New("Error creating API user")
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Error creating API user", err)
		return errors.New("Error creating API user")
	}

	user.Status = "active"
	user.UserSecret = "NO-REQUIRED"
	fmt.Println("user created")
	return nil
}

//GenerateAccessKeys ...
func GenerateAccessKeys(userUUID string, schema string) (string, string, error) {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return "", "", fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var orgName string

	if err := conn.QueryRowContext(context.Background(), GetIamOrgQuery, schema).Scan(&orgName); err != nil {

		fmt.Println("Error validating schema db", err)
		return "", "", errors.New("Error validating schema db")
	}

	var usrType, accKey, scrtKey, accStatus string

	if err := conn.QueryRowContext(context.Background(), util.ParseQuerySchema(schema, GetUserAccesByUUIDQuery), userUUID).Scan(&usrType, &accKey, &scrtKey, &accStatus); err != nil {

		fmt.Println("Error validating type and access keys", err)

		if !strings.HasPrefix(err.Error(), "sql: no rows") {

			return "", "", errors.New("Error validating type and access keys")
		}
	}

	if usrType == "platform" {

		return "", "", errors.New("BAD TYPE: user not allowed to have access keys")
	}

	if (accKey != "" || scrtKey != "") && accStatus == "active" {

		return "", "", errors.New("CONFLICT: User with active access keys")
	}

	stmt, err := conn.PrepareContext(context.Background(), util.ParseQuerySchema(schema, GenerateAccessQuery))

	if err != nil {

		fmt.Println("error preparing query", err)
		return "", "", errors.New("Error preparing Query")
	}

	accessCtrlUUID := uuid.NewV4().String()

	// TODO: remove random and crate keys based on ids, so can be really uniques
	access := util.GenerateAccesKey(orgName)
	secret := util.GenerateSecretKey()
	expiresAt := time.Now().AddDate(0, 6, 0).Format(time.RFC3339)

	if _, err := stmt.ExecContext(context.Background(), accessCtrlUUID, access, secret, userUUID, expiresAt); err != nil {

		fmt.Println("Error generating access and secret", err)
		return "", "", errors.New("Error generating access and secret")
	}

	fmt.Println("Access tokens generated")
	return access, secret, nil

}
