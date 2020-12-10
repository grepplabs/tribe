package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg"
	"net/http"
)

func NewGetRealmHandler(dbClient client.Client) apirealms.GetRealmHandler {
	return &getRealmHandler{
		dbClient: dbClient,
	}
}

type getRealmHandler struct {
	dbClient client.Client
}

func (h *getRealmHandler) Handle(input apirealms.GetRealmParams) middleware.Responder {

	realm, err := h.dbClient.RealmManager().GetRealm(input.HTTPRequest.Context(), input.RealmID)
	if err != nil {
		return apirealms.NewGetRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	if realm == nil {
		return apirealms.NewGetRealmNotFound()
	}
	return apirealms.NewGetRealmFoundRealm().WithPayload(&apimodels.GetRealmResponse{
		RealmID:     pkg.String(realm.RealmID),
		CreatedAt:   strfmt.DateTime(realm.CreatedAt),
		Description: realm.Description,
	})
}
