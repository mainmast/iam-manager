package util

import (
	"encoding/json"
	"fmt"
	"strings"

	acc "gitlab.com/mainmast/microservices/iam/iam-models.git/pkg/account"
	org "gitlab.com/mainmast/microservices/iam/iam-models.git/pkg/organisation"
	usr "gitlab.com/mainmast/microservices/iam/iam-models.git/pkg/user"
)

//ParseOrg ...
func ParseOrg(payload []byte) (*org.CreateOrgRQ, error) {

	request := &org.CreateOrgRQ{}

	if err := json.Unmarshal(payload, request); err != nil {

		fmt.Println("Err: error decoding organisation request", err)
		return nil, err
	}

	return request, nil
}

//ParseAccount ...
func ParseAccount(payload []byte) (*acc.Account, error) {

	request := &acc.Account{}

	if err := json.Unmarshal(payload, request); err != nil {

		fmt.Println("Err: error decoding account request", err)
		return nil, err
	}

	return request, nil
}

//ParseUser ...
func ParseUser(payload []byte) (*usr.User, error) {

	request := &usr.User{}

	if err := json.Unmarshal(payload, request); err != nil {

		fmt.Println("Err: error decoding user request", err)
		return nil, err
	}

	return request, nil
}

//ParseQuerySchema ...
func ParseQuerySchema(schema string, query string) string {

	return strings.Replace(query, "#*", schema, -1)
}
