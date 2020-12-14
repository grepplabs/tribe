package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg"
	"net/http"
)

func NewGetUserHandler(dbClient client.Client) apiusers.GetUserHandler {
	return &getUserHandler{
		dbClient: dbClient,
	}
}

type getUserHandler struct {
	dbClient client.Client
}

func (h *getUserHandler) Handle(input apiusers.GetUserParams) middleware.Responder {

	user, err := h.dbClient.UserManager().GetUser(input.HTTPRequest.Context(), input.RealmID, input.Username)
	if err != nil {
		return apiusers.NewGetUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	if user == nil {
		return apiusers.NewGetUserNotFound()
	}
	return apiusers.NewGetUserFoundUser().WithPayload(userToGetUserResponse(user))
}

func userToGetUserResponse(user *models.User) *apimodels.GetUserResponse {
	return &apimodels.GetUserResponse{
		UserID:        pkg.String(user.UserID),
		CreatedAt:     strfmt.DateTime(user.CreatedAt),
		RealmID:       pkg.String(user.RealmID),
		Username:      pkg.String(user.Username),
		Enabled:       user.Enabled,
		Email:         strfmt.Email(user.Email),
		EmailVerified: user.EmailVerified,
	}
}
