package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apirealms "github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg"
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
	realmList, err := h.dbClient.RealmManager().ListRealms(input.HTTPRequest.Context(), input.Offset, input.Limit)
	if err != nil {
		return apirealms.NewGetRealmDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	results := make([]*apimodels.GetRealmResponse, len(realmList.Realms))
	for i := 0; i < len(realmList.Realms); i++ {
		results[i] = realmToGetRealmResponse(&realmList.Realms[i])
	}
	payload := &apirealms.ListRealmsFoundRealmsBody{
		Results: results,
		Links: &apimodels.Links{
			Prev: prevToken(input.Offset, input.Limit),
			Next: nextToken(input.Offset, input.Limit, len(results)),
		},
		Total: pkg.Int64(int64(realmList.Page.Total)),
	}
	return apirealms.NewListRealmsFoundRealms().WithPayload(payload)
}
