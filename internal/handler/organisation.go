package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"mainmast/iam-manager/internal/orm"
	"mainmast/iam-manager/internal/util"
)

//CreateOrgHandler ...
func CreateOrgHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")

	req, err := util.ParseOrg(ctx.PostBody())

	if err != nil {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	rs, err := orm.CreateOrganisation(req)

	if err != nil {

		if strings.HasPrefix(err.Error(), "DUPLICATED") {

			ctx.SetBody([]byte(`{"message": "organisation name already exists"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return
		}

		ctx.SetBody([]byte(`{"message": "error creating Organisation"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	bt, err := json.Marshal(rs)
	if err != nil {

		ctx.SetBody([]byte(`{"message": "error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(bt)
	fmt.Println("LOG IAM-Manager/Org@Create: done")
	return
}

//AssociateUserWithOrgHandler ...
func AssociateUserWithOrgHandler(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json; charset=UTF-8")

	orgUUID := ctx.UserValue("org_uuid")
	usrUUID := ctx.UserValue("usr_uuid")

	if usrUUID == nil || orgUUID == nil {

		ctx.SetBody([]byte(`{"message": "error reading request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := orm.AssociateUserWithOrg(usrUUID.(string), orgUUID.(string))

	if err != nil {

		if strings.HasPrefix(err.Error(), "DUPLICATED") {

			ctx.SetBody([]byte(`{"message": "user already associated to this org"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return

		}

		if strings.HasPrefix(err.Error(), "NOT FOUND") {

			ctx.SetBody([]byte(`{"message": "user or organisation not found"}`))
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return

		}

		ctx.SetBody([]byte(`{"message": "error associating user in the Organisation"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody([]byte(`{"message": "user associated in the organisation succesfully"}`))
	fmt.Println("LOG IAM-Manager/Org@Create: done")
	return
}
