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

type kmsKeysetManager struct {
	mc         *minio.Client
	bucketName string
}

func (m kmsKeysetManager) CreateKMSKeyset(ctx context.Context, record *model.KMSKeyset) error {
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
		return errors.Wrap(err, "Marshal KMSKeyset failed")
	}
	_, err = m.mc.PutObject(ctx, m.bucketName, objectName, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return errors.Wrap(err, "PutObject failed")
	}
	return nil
}

func (m kmsKeysetManager) GetKMSKeyset(ctx context.Context, id string) (*model.KMSKeyset, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	return m.getObject(ctx, m.objectNameForID(id))
}

func (m kmsKeysetManager) getObject(ctx context.Context, objectName string) (*model.KMSKeyset, error) {
	reader, err := m.mc.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "GetObject failed")
	}
	defer reader.Close()
	var record model.KMSKeyset
	err = json.NewDecoder(reader).Decode(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (m kmsKeysetManager) DeleteKMSKeyset(ctx context.Context, id string) error {
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

func (m kmsKeysetManager) UpdateKMSKeyset(ctx context.Context, record *model.KMSKeyset) error {
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

func (m kmsKeysetManager) ListKMSKeysets(ctx context.Context, offset *int64, limit *int64) (*model.KMSKeysetList, error) {
	var maxKeys int
	if limit != nil && *limit > 0 {
		// limit 0 all elements
		maxKeys = int(*limit)
	}
	list := make([]model.KMSKeyset, 0)
	for object := range m.mc.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{Prefix: m.objectPrefix(), Recursive: true, MaxKeys: maxKeys}) {
		kmsKeyset, err := m.getObject(ctx, object.Key)
		if err != nil {
			return nil, err
		}
		list = append(list, *kmsKeyset)
	}
	return &model.KMSKeysetList{List: list, Page: model.Page{
		Offset: nil, // TODO: offset was not used
		Limit:  limit,
		Total:  0, // TODO: Total is unknown
	}}, nil
}

func (m kmsKeysetManager) objectNameForID(id string) string {
	return fmt.Sprintf("%s%s", m.objectPrefix(), id)
}

func (m kmsKeysetManager) objectPrefix() string {
	var kmsKeyset model.KMSKeyset
	return fmt.Sprintf("%s/", kmsKeyset.TableName())
}

func (m kmsKeysetManager) existsObjectWithName(ctx context.Context, objectName string) (bool, error) {
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
