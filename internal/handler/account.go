package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/mainmast/iam-manager/internal/orm"
	"github.com/mainmast/iam-manager/internal/util"
)

//CreateAccountHandler ...
func CreateAccountHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")
	schema := string(ctx.Request.Header.Peek("X-CMO-DB"))
	req, err := util.ParseAccount(ctx.PostBody())

	if err != nil || schema == "" {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := orm.CreateAccount(req, schema); err != nil {

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
	fmt.Println("LOG IAM-Manager/Account@Create: done")
	return
}

//AssociateUserWithAccHandler ...
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
