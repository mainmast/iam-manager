package util

import (

	migration "github.com/mainmast/iam-migrations/pkg"
)

//SetUpCustomer ...
func SetUpCustomer(action string, schema string) bool {

	return migration.Migrate(action, schema)
}
