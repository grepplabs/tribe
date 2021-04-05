package client

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/service"
	"github.com/grepplabs/tribe/database/service/clientminio"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient struct {
	mc     *minio.Client
	logger log.Logger
	api    service.API
}

func NewMinioClient(logger log.Logger, config *config.MinioConfig) (Client, error) {
	var mc, err = minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mc.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{Region: config.BucketLocation})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := mc.BucketExists(ctx, config.BucketName)
		if !(errBucketExists == nil && exists) {
			return nil, errors.Wrapf(err, "Make bucket %s failed. Bucket exists %v, check error %v", config.BucketName, exists, errBucketExists)
		}
	}
	return &minioClient{
		mc:     mc,
		logger: logger,
		api:    clientminio.NewAPIImpl(mc, config),
	}, nil
}

func (c minioClient) API() service.API {
	return c.api
}
