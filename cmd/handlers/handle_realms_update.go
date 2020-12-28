package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	dtomodels "github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg"
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
	realm := &dtomodels.Realm{
		RealmID:     input.RealmID,
		Description: input.Realm.Description,
	}
	err := h.dbClient.RealmManager().UpdateRealm(input.HTTPRequest.Context(), realm)
	if err != nil {
		var notFound pkg.ErrNotFound
		if errors.As(err, &notFound) {
			return apirealms.NewUpdateRealmNotFound()
		}
		return apirealms.NewUpdateRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	return apirealms.NewUpdateRealmUpdatedRealm()
}
