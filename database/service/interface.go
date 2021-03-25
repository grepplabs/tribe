package service

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
)

type API interface {
	CreateKMSKeyset(ctx context.Context, u *model.KMSKeyset) error
	DeleteKMSKeysetr(ctx context.Context, id string) error
	UpdateKMSKeyset(ctx context.Context, u *model.KMSKeyset) error
	GetKMSKeyset(ctx context.Context, id string) (*model.KMSKeyset, error)
	GetKMSKeysetsByName(ctx context.Context, name string) (*model.KMSKeysetList, error)
}
