package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
	"net/http"
)

func NewDeleteUserHandler(dbClient client.Client) apiusers.DeleteUserHandler {
	return &deleteUserHandler{
		dbClient: dbClient,
	}
}

type deleteUserHandler struct {
	dbClient client.Client
}

func (h *deleteUserHandler) Handle(input apiusers.DeleteUserParams) middleware.Responder {
	exists, err := h.dbClient.UserManager().ExistsUser(input.HTTPRequest.Context(), input.RealmID, input.Username)
	if err != nil {
		return h.newInternalError(err)
	}
	if !exists {
		return apiusers.NewDeleteUserNotFound()
	}
	err = h.dbClient.UserManager().DeleteUser(input.HTTPRequest.Context(), input.RealmID, input.Username)
	if err != nil {
		return h.newInternalError(err)
	}
	return apiusers.NewDeleteUserOK()
}

func (h *deleteUserHandler) newInternalError(err error) *apiusers.DeleteUserDefault {
	return apiusers.NewDeleteUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Detail:  err.Error(),
	})
}
