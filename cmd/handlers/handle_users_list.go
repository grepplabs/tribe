package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg"
	"net/http"
)

func NewListUsersHandler(dbClient client.Client) apiusers.ListUsersHandler {
	return &listUsersHandler{
		dbClient: dbClient,
	}
}

type listUsersHandler struct {
	dbClient client.Client
}

func (h *listUsersHandler) Handle(input apiusers.ListUsersParams) middleware.Responder {
	userList, err := h.dbClient.UserManager().ListUsers(input.HTTPRequest.Context(), input.RealmID, input.Offset, input.Limit)
	if err != nil {
		return apiusers.NewGetUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	results := make([]*apimodels.GetUserResponse, len(userList.Users))
	for i := 0; i < len(userList.Users); i++ {
		results[i] = userToGetUserResponse(&userList.Users[i])
	}
	payload := &apiusers.ListUsersFoundUsersBody{
		Results: results,
		Links: &apimodels.Links{
			Prev: prevToken(input.Offset, input.Limit),
			Next: nextToken(input.Offset, input.Limit, len(results)),
		},
		Total: pkg.Int64(int64(userList.Page.Total)),
	}
	return apiusers.NewListUsersFoundUsers().WithPayload(payload)
}
