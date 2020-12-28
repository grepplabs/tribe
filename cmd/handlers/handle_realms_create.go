package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	dtomodels "github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg"
	"net/http"
	"time"
)

func NewCreateRealmHandler(dbClient client.Client) apirealms.CreateRealmHandler {
	return &createRealmHandler{
		dbClient: dbClient,
	}
}

type createRealmHandler struct {
	dbClient client.Client
}

func (h *createRealmHandler) Handle(input apirealms.CreateRealmParams) middleware.Responder {
	realm := &dtomodels.Realm{
		RealmID:     pkg.StringValue(input.Realm.RealmID),
		CreatedAt:   time.Now(),
		Description: input.Realm.Description,
	}

	err := h.dbClient.RealmManager().CreateRealm(input.HTTPRequest.Context(), realm)
	if err != nil {
		return h.newInternalError(err)
	}
	return apirealms.NewCreateRealmCreated()
}

func (h *createRealmHandler) newInternalError(err error) *apirealms.CreateRealmDefault {
	return apirealms.NewCreateRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Detail:  err.Error(),
	})
}
