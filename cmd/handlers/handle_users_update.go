package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
	dtomodels "github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg/crypto"
	"net/http"
)

func NewUpdateUserHandler(dbClient client.Client, bcryptCost int) apiusers.UpdateUserHandler {
	return &updateUserHandler{
		dbClient:   dbClient,
		bcryptCost: bcryptCost,
	}
}

type updateUserHandler struct {
	dbClient   client.Client
	bcryptCost int
}

func (h *updateUserHandler) Handle(input apiusers.UpdateUserParams) middleware.Responder {
	oldUser, err := h.dbClient.UserManager().GetUser(input.HTTPRequest.Context(), input.RealmID, input.Username)
	if err != nil {
		return h.newInternalError(err)
	}
	if oldUser == nil {
		return apiusers.NewUpdateUserNotFound()
	}
	passwordHash := oldUser.EncryptedPassword
	if input.User.Password != "" {
		hasher := crypto.NewPasswordHasher()
		if h.bcryptCost > 0 {
			hasher = hasher.WithBCryptCost(h.bcryptCost)
		}
		passwordHash, err = hasher.HashPassword(input.User.Password.String())
		if err != nil {
			return apiusers.NewUpdateUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Detail:  err.Error(),
			})
		}
	}
	user := &dtomodels.User{
		UserID:            oldUser.UserID,
		RealmID:           input.RealmID,
		Username:          input.Username,
		EncryptedPassword: passwordHash,
		Enabled:           input.User.Enabled,
		Email:             input.User.Email.String(),
		EmailVerified:     input.User.EmailVerified,
	}

	if user.Email == "" {
		user.EmailVerified = false
	}

	err = h.dbClient.UserManager().UpdateUser(input.HTTPRequest.Context(), user)
	if err != nil {
		return h.newInternalError(err)
	}
	return apiusers.NewUpdateUserOK()
}

func (h *updateUserHandler) newInternalError(err error) *apiusers.UpdateUserDefault {
	return apiusers.NewUpdateUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Detail:  err.Error(),
	})
}
