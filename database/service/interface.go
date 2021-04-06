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
	ListKMSKeysets(ctx context.Context, offset *int64, limit *int64) (*model.KMSKeysetList, error)

	CreateJWKS(ctx context.Context, jwks *model.JWKS) error
	GetJWKS(ctx context.Context, id string) (*model.JWKS, error)
	GetJWKSByKidUse(ctx context.Context, kid string, use string) (*model.JWKS, error)
	DeleteJWKS(ctx context.Context, id string) error
	DeleteJWKSByKidUse(ctx context.Context, kid string, use string) error
	ListJWKS(ctx context.Context, offset *int64, limit *int64) (*model.JWKSList, error)
}
