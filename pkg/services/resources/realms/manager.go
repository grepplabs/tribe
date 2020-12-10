package realms

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
)

type Manager interface {
	CreateRealm(ctx context.Context, r *models.Realm) error
	GetRealm(ctx context.Context, realmID string) (*models.Realm, error)
}
