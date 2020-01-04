package main

import (
	"os"

	"github.com/mainmast/iam-manager/internal/handler"

	wb "github.com/mainmast/httpm/pkg/webserver"
)

func main() {

	port := os.Getenv("PORT")
	server := &wb.WebServer{}

	server.AddHandler("/iam/org", "POST", handler.CreateOrgHandler)
	server.AddHandler("/iam/acc", "POST", handler.CreateAccountHandler)
	server.AddHandler("/iam/usr", "POST", handler.CreateUserHandler)
	server.AddHandler("/iam/usr/api", "POST", handler.CreateUserAPIHandler)
	server.AddHandler("/iam/access/keys/:usr_uuid", "POST", handler.GenerateAccessKeysHandler)
	server.AddHandler("/iam/org/link/user/:org_uuid/:usr_uuid", "PATCH", handler.AssociateUserWithOrgHandler)
	server.AddHandler("/iam/acc/link/user/:acc_uuid", "PATCH", handler.AssociateUserWithAccHandler)

	if port == "" {

		port = "8080"
	}

	server.StartUp(port)
}
