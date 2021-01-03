package users

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
)

type Manager interface {
	CreateUser(ctx context.Context, u *models.User) error
	DeleteUser(ctx context.Context, realmID string, username string) error
	GetUser(ctx context.Context, realmID string, username string) (*models.User, error)
	ListUsers(ctx context.Context, realmID string, offset *int64, limit *int64) (*models.UserList, error)
	ExistsUser(ctx context.Context, realmID string, username string) (bool, error)
}
