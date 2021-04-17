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

func (m oidcJwksManager) CreateOidcJWKS(ctx context.Context, record *model.OidcJWKS) error {
	if record == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter record is missing"}
	}
	objectName := m.objectNameForID(record.ID)
	exists, err := m.existsObjectWithName(ctx, objectName)
	if err != nil {
		return err
	}
	if exists {
		return service.ErrAlreadyExists{Reason: objectName}
	}
	data, err := json.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "Marshal OidcJWKS failed")
	}
	_, err = m.mc.PutObject(ctx, m.bucketName, objectName, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return errors.Wrap(err, "PutObject failed")
	}
	return nil
}

func (m oidcJwksManager) GetOidcJWKS(ctx context.Context, id string) (*model.OidcJWKS, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id for OidcJWKS is missing"}
	}
	return m.getObject(ctx, m.objectNameForID(id))
}

func (m oidcJwksManager) getObject(ctx context.Context, objectName string) (*model.OidcJWKS, error) {
	reader, err := m.mc.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "GetObject failed")
	}
	defer reader.Close()
	var record model.OidcJWKS
	err = json.NewDecoder(reader).Decode(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (m oidcJwksManager) DeleteOidcJWKS(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	objectName := m.objectNameForID(id)
	err := m.mc.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "RemoveObject failed")
	}
	return nil
}

func (m oidcJwksManager) UpdateOidcJWKS(ctx context.Context, record *model.OidcJWKS) error {
	if record == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter record is missing"}
	}
	objectName := m.objectNameForID(record.ID)
	data, err := json.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "Marshal KMSKeyset failed")
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
