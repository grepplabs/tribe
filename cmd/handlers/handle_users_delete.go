package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	"net/http"
)

func NewDeleteRealmHandler(dbClient client.Client) apirealms.DeleteRealmHandler {
	return &deleteRealmHandler{
		dbClient: dbClient,
	}
}

type deleteRealmHandler struct {
	dbClient client.Client
}

func (h *deleteRealmHandler) Handle(input apirealms.DeleteRealmParams) middleware.Responder {
	exists, err := h.dbClient.RealmManager().ExistsRealm(input.HTTPRequest.Context(), input.RealmID)
	if err != nil {
		return h.newInternalError(err)
	}
	if !exists {
		return apirealms.NewDeleteRealmNotFound()
	}
	err = h.dbClient.RealmManager().DeleteRealm(input.HTTPRequest.Context(), input.RealmID)
	if err != nil {
		return h.newInternalError(err)
	}
	return apirealms.NewDeleteRealmOK()
}

func (h *deleteRealmHandler) newInternalError(err error) *apirealms.DeleteRealmDefault {
	return apirealms.NewDeleteRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Detail:  err.Error(),
	})
}
