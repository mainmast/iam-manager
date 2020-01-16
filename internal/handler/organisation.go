package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"mainmast/iam-manager/internal/orm"
	"mainmast/iam-manager/internal/util"
)

// Creates organisations within the platform.
func CreateOrganisationHandler(ctx *fasthttp.RequestCtx) {

	// Set the content type to JSON
	ctx.SetContentType("application/json; charset=UTF-8")

	// Parse the body of the JSON request
	req, err := util.ParseOrg(ctx.PostBody())

	if err != nil {
		ctx.SetBody([]byte(`{"message": "Error reading the request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Create the organisation, and initial user / accounts / access keys.
	organisation, err := orm.CreateOrganisation(req)

	if err != nil {

		// The organisation name already exists.
		if strings.HasPrefix(err.Error(), "DUPLICATED") {
			ctx.SetBody([]byte(`{"message": "The organisation name already exists"}`))
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return
		}

		// Some other error occurred.
		// TODO: Fix this to show the real error.
		ctx.SetBody([]byte(`{"message": "An error occurred while trying to create the organisation."}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	// Build response JSON from the organisation object.
	responseBody, err := json.Marshal(organisation)

	if err != nil {
		ctx.SetBody([]byte(`{"message": "Error encoding the response"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(responseBody)

	fmt.Println("LOG IAM-Manager/Org@Create: done")

	return
}

// Associates a user with an existing organisation
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
