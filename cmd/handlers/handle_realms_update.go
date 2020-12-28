package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	dtomodels "github.com/grepplabs/tribe/database/models"
	"net/http"
)

func NewUpdateRealmHandler(dbClient client.Client) apirealms.UpdateRealmHandler {
	return &updateRealmHandler{
		dbClient: dbClient,
	}
}

type updateRealmHandler struct {
	dbClient client.Client
}

func (h *updateRealmHandler) Handle(input apirealms.UpdateRealmParams) middleware.Responder {
	exists, err := h.dbClient.RealmManager().ExistsRealm(input.HTTPRequest.Context(), input.RealmID)
	if err != nil {
		return h.newInternalError(err)
	}
	if !exists {
		return apirealms.NewUpdateRealmNotFound()
	}

	realm := &dtomodels.Realm{
		RealmID:     input.RealmID,
		Description: input.Realm.Description,
	}

	err = h.dbClient.RealmManager().UpdateRealm(input.HTTPRequest.Context(), realm)
	if err != nil {
		return h.newInternalError(err)
	}
	return apirealms.NewUpdateRealmUpdatedRealm()
}

func (h *updateRealmHandler) newInternalError(err error) *apirealms.UpdateRealmDefault {
	return apirealms.NewUpdateRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Detail:  err.Error(),
	})
}
