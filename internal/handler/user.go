package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"mainmast/iam-manager/internal/orm"
	"mainmast/iam-manager/internal/util"
)

// HTTP Handler for creating a platform user
// TODO: Better error responses
func CreateUserHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=UTF-8")

	orgUUID := string(ctx.QueryArgs().Peek("org_uuid"))
	user, err := util.ParseUser(ctx.PostBody())

	if err != nil {
		ctx.SetBody([]byte(`{"message": "Error reading user request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := orm.CreateIamUser(user, orgUUID); err != nil {
		if strings.HasPrefix(err.Error(), "DUPLICATED") {

			ctx.SetBody([]byte(`{"message": "User login already exists"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return

		}
		ctx.SetBody([]byte(`{"message": "Error creating the user"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	responseJson, err := json.Marshal(user)

	if err != nil {
		ctx.SetBody([]byte(`{"message": "Error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(responseJson)
	fmt.Println("LOG IAM-Manager/User@Create: done")
	return
}

// Create Organisation Platform User
func CreateUserAPIHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")
	schema := string(ctx.Request.Header.Peek("X-CMO-DB"))
	req, err := util.ParseUser(ctx.PostBody())

	if err != nil || schema == "" {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := orm.CreateUserAPI(req, schema); err != nil {

		if strings.HasPrefix(err.Error(), "DUPLICATED") {

			ctx.SetBody([]byte(`{"message": "user login already exists"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return

		}
		ctx.SetBody([]byte(`{"message": "error creating the account"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	bt, err := json.Marshal(req)
	if err != nil {

		ctx.SetBody([]byte(`{"message": "error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(bt)
	fmt.Println("LOG IAM-Manager/User@CreateUserAPI: done")
	return
}

//GenerateAccessKeysHandler ...
func GenerateAccessKeysHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")
	schema := string(ctx.Request.Header.Peek("X-CMO-DB"))
	usrUUID := ctx.UserValue("usr_uuid")

	if schema == "" || usrUUID == nil {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	access, secret, err := orm.GenerateAccessKeys(usrUUID.(string), schema)

	if err != nil {

		if strings.HasPrefix(err.Error(), "BAD TYPE") {

			ctx.SetBody([]byte(`{"message": "user not allowed to have access keys"}`))
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return

		} else if strings.HasPrefix(err.Error(), "CONFLICT") {

			ctx.SetBody([]byte(`{"message": "user with active access keys"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return

		}

		ctx.SetBody([]byte(`{"message": "error generating access keys"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	rs := struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	}{
		AccessKey: access,
		SecretKey: secret,
	}

	bt, err := json.Marshal(rs)
	if err != nil {

		ctx.SetBody([]byte(`{"message": "error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bt)
	fmt.Println("LOG IAM-Manager/User@GenerateAccessKeys: done")
	return
}
