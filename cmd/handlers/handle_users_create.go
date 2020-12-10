package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	apimodels "github.com/grepplabs/tribe/api/v1/models"
	apiusers "github.com/grepplabs/tribe/api/v1/server/restapi/users"
	"github.com/grepplabs/tribe/database/client"
	dtomodels "github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg"
	"github.com/grepplabs/tribe/pkg/crypto"
	"net/http"
	"time"
)

func NewCreateUserHandler(dbClient client.Client, bcryptCost int) apiusers.CreateUserHandler {
	return &createUserHandler{
		dbClient:   dbClient,
		bcryptCost: bcryptCost,
	}
}

type createUserHandler struct {
	dbClient   client.Client
	bcryptCost int
}

func (h *createUserHandler) Handle(input apiusers.CreateUserParams) middleware.Responder {
	hasher := crypto.NewPasswordHasher()
	if h.bcryptCost > 0 {
		hasher = hasher.WithBCryptCost(h.bcryptCost)
	}
	passwordHash, err := hasher.HashPassword(input.User.Password.String())
	if err != nil {
		return apiusers.NewCreateUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	// generate user UUID
	userID := uuid.New().String()

	user := &dtomodels.User{
		UserID:            userID,
		CreatedAt:         time.Now(),
		RealmID:           input.RealmID,
		Username:          pkg.StringValue(input.User.Username),
		EncryptedPassword: passwordHash,
		Enabled:           input.User.Enabled,
		Email:             input.User.Email.String(),
		EmailVerified:     input.User.EmailVerified,
	}
	if user.Email == "" {
		user.EmailVerified = false
	}

	err = h.dbClient.UserManager().CreateUser(input.HTTPRequest.Context(), user)
	if err != nil {
		return apiusers.NewCreateUserDefault(http.StatusInternalServerError).WithPayload(&apimodels.Problem{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Detail:  err.Error(),
		})
	}
	return apiusers.NewCreateUserCreatedUser().WithPayload(
		&apimodels.CreateUserResponse{
			RealmID:  &input.RealmID,
			UserID:   &userID,
			Username: input.User.Username,
		},
	)
}
