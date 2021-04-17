package clientminio

import (
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/service"
	"github.com/minio/minio-go/v7"
)

type APIImpl struct {
	kmsKeysetManager
	jwksManager
	oidcJwksManager
}

var _ service.API = (*APIImpl)(nil)

func NewAPIImpl(mc *minio.Client, config *config.MinioConfig) *APIImpl {
	return &APIImpl{
		kmsKeysetManager{mc, config.BucketName},
		jwksManager{mc, config.BucketName},
		oidcJwksManager{mc, config.BucketName},
	}
}
