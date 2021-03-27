package service

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
)

type API interface {
	CreateKMSKeyset(ctx context.Context, kmsKeyset *model.KMSKeyset) error
	DeleteKMSKeyset(ctx context.Context, keysetID string) error
	UpdateKMSKeyset(ctx context.Context, kmsKeyset *model.KMSKeyset) error
	GetKMSKeyset(ctx context.Context, keysetID string) (*model.KMSKeyset, error)
	GetKMSKeysetsByName(ctx context.Context, name string) (*model.KMSKeysetList, error)
	ListKMSKeysets(ctx context.Context, offset *int64, limit *int64) (*model.KMSKeysetList, error)
}
