package users

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
)

type Manager interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUser(ctx context.Context, realmID string, username string) (*models.User, error)
}
