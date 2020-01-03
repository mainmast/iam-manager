package orm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"

	//postgres driver
	_ "github.com/lib/pq"

	"mainmast/iam-manager/internal/util"

	acc "gitlab.com/mainmast/microservices/iam/iam-models.git/pkg/account"
)

//CreateAccount ...
func CreateAccount(account *acc.Account, schema string) error {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	stmt, err := conn.PrepareContext(context.Background(), util.ParseQuerySchema(schema, CreateAccountQuery))

	if err != nil {

		fmt.Println("error preparing query", err)
		return errors.New("Error preparing Query")
	}

	account.UUID = uuid.NewV4().String()
	account.Name = strings.ToLower(account.Name)
	account.IsRoot = false

	accCustomBts, err := json.Marshal(account.CustomData)

	if err != nil {

		fmt.Println("error reading account custom data", err)
		return errors.New("error reading account custom data")
	}

	if _, err := stmt.ExecContext(context.Background(), account.UUID, account.Name, account.IsRoot, accCustomBts); err != nil {

		fmt.Println("Error creating account", err)
		return errors.New("account not created")
	}

	fmt.Println("Account created")
	return nil

}

//AssociateUserWithAccount ....
func AssociateUserWithAccount(UserUUID string, AccountUUID string, schema string) error {

	conn, err := GetDBInstance()

	if err != nil {

		fmt.Println("Error connecting with dabase", err)
		return fmt.Errorf("Error creating connection with the database %v", err)
	}

	defer conn.Close()

	var isInDefault bool
	//Check if the user belongs to the Organisation before linking him/her
	err = conn.QueryRowContext(context.Background(), util.ParseQuerySchema(schema, UserInDefaultAccountQuery), UserUUID).Scan(&isInDefault)

	if err != nil || !isInDefault {

		fmt.Println("Error User does not belong to this Organisation, need to relate the user first", err)
		return errors.New("NOT LINKED: Error User does not belong to this Organisation, need to relate the user first")
	}

	stmt, err := conn.PrepareContext(context.Background(), util.ParseQuerySchema(schema, AssociateUserToAccountQuery))

	if err != nil {

		fmt.Println("error preparing query", err)
		return errors.New("Error preparing Query")
	}

	if _, err := stmt.ExecContext(context.Background(), UserUUID, AccountUUID); err != nil {

		fmt.Println("Error associating user to the account", err)
		return errors.New("Error associating user to the account")
	}

	fmt.Println("User associated with account sucessfuly")
	return nil
}
