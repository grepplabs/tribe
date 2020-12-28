package realms

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
)

type Manager interface {
	CreateRealm(ctx context.Context, r *models.Realm) error
	UpdateRealm(ctx context.Context, r *models.Realm) error
	DeleteRealm(ctx context.Context, realmID string) error
	GetRealm(ctx context.Context, realmID string) (*models.Realm, error)
	ListRealms(ctx context.Context, offset *int64, limit *int64) ([]models.Realm, error)
}
