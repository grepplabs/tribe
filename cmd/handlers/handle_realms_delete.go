package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg"
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
	err := h.dbClient.RealmManager().DeleteRealm(input.HTTPRequest.Context(), input.RealmID)
	if err != nil {
		var notFound pkg.ErrNotFound
		if errors.As(err, &notFound) {
			return apirealms.NewDeleteRealmNotFound()
		}
		return apirealms.NewDeleteRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	return apirealms.NewDeleteRealmDeletedRealm()
}
