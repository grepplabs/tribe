package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	"net/http"
)

func NewListRealmsHandler(dbClient client.Client) apirealms.ListRealmsHandler {
	return &listRealmsHandler{
		dbClient: dbClient,
	}
}

type listRealmsHandler struct {
	dbClient client.Client
}

func (h *listRealmsHandler) Handle(input apirealms.ListRealmsParams) middleware.Responder {
	users, err := h.dbClient.RealmManager().ListRealms(input.HTTPRequest.Context(), input.Offset, input.Limit)
	if err != nil {
		return apirealms.NewGetRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	results := make([]*apimodels.GetRealmResponse, len(users))
	for i := 0; i < len(users); i++ {
		results[i] = realmToGetRealmResponse(&users[i])
	}
	payload := &apirealms.ListRealmsFoundRealmsBody{
		Results: results,
		Links: &apimodels.Links{
			Prev: prevToken(input.Offset, input.Limit),
			Next: nextToken(input.Offset, input.Limit, len(results)),
		},
	}
	return apirealms.NewListRealmsFoundRealms().WithPayload(payload)
}
