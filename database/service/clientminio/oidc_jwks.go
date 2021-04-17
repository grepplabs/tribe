package clientminio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/database/service"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

type oidcJwksManager struct {
	mc         *minio.Client
	bucketName string
}

func (m oidcJwksManager) CreateOidcJWKS(ctx context.Context, oidcJWKS *model.OidcJWKS) error {
	if oidcJWKS == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter oidcJWKS is missing"}
	}
	objectName := m.objectNameForID(oidcJWKS.ID)
	exists, err := m.existsObjectWithName(ctx, objectName)
	if err != nil {
		return err
	}
	if exists {
		return service.ErrAlreadyExists{Reason: objectName}
	}
	data, err := json.Marshal(oidcJWKS)
	if err != nil {
		return errors.Wrap(err, "Marshal oidcJWKS failed")
	}
	_, err = m.mc.PutObject(ctx, m.bucketName, objectName, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return errors.Wrap(err, "PutObject failed")
	}
	return nil
}

func (m oidcJwksManager) objectNameForID(id string) string {
	return fmt.Sprintf("%s%s", m.objectPrefix(), id)
}

func (m oidcJwksManager) objectPrefix() string {
	var oidcJWKS model.OidcJWKS
	return fmt.Sprintf("%s/", oidcJWKS.TableName())
}

func (m oidcJwksManager) existsObjectWithName(ctx context.Context, objectName string) (bool, error) {
	_, err := m.mc.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
