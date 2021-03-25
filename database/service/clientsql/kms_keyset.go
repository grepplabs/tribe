package clientsql

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/upper/db/v4"
)

type KMSKeysetManager struct {
	dbs db.Session
}

func (m KMSKeysetManager) CreateKMSKeyset(ctx context.Context, u *model.KMSKeyset) error {
	return nil
}
func (m KMSKeysetManager) DeleteKMSKeysetr(ctx context.Context, id string) error {
	return nil

}
func (m KMSKeysetManager) UpdateKMSKeyset(ctx context.Context, u *model.KMSKeyset) error {
	return nil

}
func (m KMSKeysetManager) GetKMSKeyset(ctx context.Context, id string) (*model.KMSKeyset, error) {
	return nil, nil

}
func (m KMSKeysetManager) GetKMSKeysetsByName(ctx context.Context, name string) (*model.KMSKeysetList, error) {
	return nil, nil
}
