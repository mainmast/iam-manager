package util

import (

	migration "gitlab.com/mainmast/microservices/iam/migrations.git/pkg"
)

//SetUpCustomer ...
func SetUpCustomer(action string, schema string) bool {

	return migration.Migrate(action, schema)
}
