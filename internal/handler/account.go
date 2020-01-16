package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"mainmast/iam-manager/internal/orm"
	"mainmast/iam-manager/internal/util"
)

// Create organisation platform accounts
func CreateAccountHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")

	schema := string(ctx.Request.Header.Peek("X-CMO-DB"))

	if schema == "" {
		ctx.SetBody([]byte(`{"message": "Could not determine organisation."}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Create an account object from the JSON request
	account, err := util.ParseAccount(ctx.PostBody())

	if err != nil || schema == "" {
		ctx.SetBody([]byte(`{"message": "Error reading account request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Create the account
	if err := orm.CreateAccount(account, schema); err != nil {
		// TODO: Better error reporting on why the account couldn't be created
		ctx.SetBody([]byte(`{"message": "Error creating the account"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	// Marshal the account object into a JSON response
	responseJson, err := json.Marshal(account)

	if err != nil {
		// TODO: Better error reporting
		ctx.SetBody([]byte(`{"message": "error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(responseJson)

	fmt.Println("LOG IAM-Manager/Account@Create: done")

	return
}

// Associate a user with a platform account.
func AssociateUserWithAccHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")

	accUUID := ctx.UserValue("acc_uuid")
	usrUUID := ctx.QueryArgs().Peek("usr_uuid")
	schema := string(ctx.Request.Header.Peek("X-CMO-DB"))

	if usrUUID == nil || accUUID == nil || schema == "" {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := orm.AssociateUserWithAccount(string(usrUUID), accUUID.(string), schema)

	if err != nil {

		if strings.HasPrefix(err.Error(), "NOT LINKED:") {

			ctx.SetBody([]byte(`{"message": "user does not belong to this organisation, go and associate him first"}`))
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return

		}

		ctx.SetBody([]byte(`{"message": "error associating user into the account"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	ctx.SetBody([]byte(`{"message": "user associated in the account succesfully"}`))
	fmt.Println("LOG IAM-Manager/Org@AssociateUsrInOrg: done")
	return
}
