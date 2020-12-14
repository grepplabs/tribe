package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
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
	users, err := h.dbClient.UserManager().ListUsers(input.HTTPRequest.Context(), input.RealmID, input.Offset, input.Limit)
	if err != nil {
		return apiusers.NewGetUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	payload := make([]*apimodels.GetUserResponse, len(users))
	for i := 0; i < len(users); i++ {
		payload[i] = userToGetUserResponse(&users[i])
	}
	return apiusers.NewListUsersFoundUsers().WithPayload(payload)
}
