package clientminio

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/minio/minio-go/v7"
)

type jwksManager struct {
	mc         *minio.Client
	bucketName string
}

func (m jwksManager) CreateJWKS(ctx context.Context, jwks *model.JWKS) error {
	//TODO: implement me
	return nil
}

func (m jwksManager) GetJWKS(ctx context.Context, id string) (*model.JWKS, error) {
	//TODO: implement me
	return nil, nil
}

func (m jwksManager) GetJWKSByKidUse(ctx context.Context, kid string, use string) (*model.JWKS, error) {
	//TODO: implement me
	return nil, nil
}

func (m jwksManager) DeleteJWKS(ctx context.Context, id string) error {
	//TODO: implement me
	return nil
}

func (m jwksManager) DeleteJWKSByKidUse(ctx context.Context, kid string, use string) error {
	//TODO: implement me
	return nil
}
