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
	"io/ioutil"
	"strconv"
	"strings"
)

type jwksManager struct {
	mc         *minio.Client
	bucketName string
}

func (m jwksManager) CreateJWKS(ctx context.Context, jwks *model.JWKS) error {
	if jwks == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter jwks is missing"}
	}
	objectName := m.objectNameForID(jwks.ID)
	exists, err := m.existsObjectWithName(ctx, objectName)
	if err != nil {
		return err
	}
	if exists {
		return service.ErrAlreadyExists{Reason: objectName}
	}
	exists, otherObjectName, err := m.existsObjectWithKidUse(ctx, jwks.Kid, jwks.Use)
	if err != nil {
		return err
	}
	if exists {
		return service.ErrAlreadyExists{Reason: fmt.Sprintf("object '%s' with kid '%s' and use '%s'", otherObjectName, jwks.Kid, jwks.Use)}
	}
	data, err := json.Marshal(jwks)
	if err != nil {
		return errors.Wrap(err, "Marshal jwks failed")
	}
	_, err = m.mc.PutObject(ctx, m.bucketName, objectName, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return errors.Wrap(err, "PutObject failed")
	}
	return nil
}

func (m jwksManager) GetJWKS(ctx context.Context, id string) (*model.JWKS, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	return m.getObject(ctx, m.objectNameForID(id))
}

func (m jwksManager) getObject(ctx context.Context, objectName string) (*model.JWKS, error) {
	reader, err := m.mc.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "GetObject failed")
	}
	defer reader.Close()
	var jwks model.JWKS
	err = json.NewDecoder(reader).Decode(&jwks)
	if err != nil {
		return nil, err
	}
	return &jwks, nil
}

func (m jwksManager) DeleteJWKS(ctx context.Context, id string) error {
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

func (m jwksManager) ListJWKS(ctx context.Context, offset *int64, limit *int64) (*model.JWKSList, error) {
	var maxKeys int
	if limit != nil && *limit > 0 {
		// limit 0 all elements
		maxKeys = int(*limit)
	}
	list := make([]model.JWKS, 0)
	for object := range m.mc.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{Prefix: m.objectPrefix(), Recursive: true, MaxKeys: maxKeys}) {
		jwks, err := m.getObject(ctx, object.Key)
		if err != nil {
			return nil, err
		}
		list = append(list, *jwks)
	}
	return &model.JWKSList{List: list, Page: model.Page{
		Offset: nil, // TODO: offset was not used
		Limit:  limit,
		Total:  0, // TODO: Total is unknown
	}}, nil
}

func (m jwksManager) GetJWKSByKidUse(ctx context.Context, kid string, use string) (*model.JWKS, error) {
	exists, objectName, err := m.existsObjectWithKidUse(ctx, kid, use)
	if err != nil {
		return nil, err
	}
	if exists {
		return m.getObject(ctx, objectName)
	}
	return nil, nil
}

func (m jwksManager) DeleteJWKSByKidUse(ctx context.Context, kid string, use string) error {
	exists, objectName, err := m.existsObjectWithKidUse(ctx, kid, use)
	if err != nil {
		return err
	}
	if exists {
		err := m.mc.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			return errors.Wrap(err, "RemoveObject failed")
		}
	}
	return nil
}

func (m jwksManager) objectNameForID(id string) string {
	return fmt.Sprintf("%s%s", m.objectPrefix(), id)
}

func (m jwksManager) objectPrefix() string {
	var jwks model.JWKS
	return fmt.Sprintf("%s/", jwks.TableName())
}

func (m jwksManager) existsObjectWithName(ctx context.Context, objectName string) (bool, error) {
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

func (m jwksManager) existsObjectWithKidUse(ctx context.Context, kid string, use string) (bool, string, error) {
	// make it parallel ?
	for object := range m.mc.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{Prefix: m.objectPrefix(), Recursive: true}) {
		exists, err := m.hasKidUse(object.Key, kid, use)
		if err != nil {
			return false, "", err
		}
		if exists {
			return true, object.Key, nil
		}
	}
	return false, "", nil
}

func (m jwksManager) hasKidUse(objectName string, kid string, use string) (bool, error) {
	opts := minio.SelectObjectOptions{
		Expression:     fmt.Sprintf("select count(*) from s3object where kid='%s' and use='%s'", kid, use),
		ExpressionType: minio.QueryExpressionTypeSQL,
		InputSerialization: minio.SelectObjectInputSerialization{
			CompressionType: minio.SelectCompressionNONE,
			JSON: &minio.JSONInputOptions{
				Type: minio.JSONDocumentType,
			},
		},
		OutputSerialization: minio.SelectObjectOutputSerialization{
			CSV: &minio.CSVOutputOptions{},
		},
	}
	reader, err := m.mc.SelectObjectContent(context.Background(), m.bucketName, objectName, opts)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return false, err
	}
	i, err := strconv.Atoi(strings.Replace(string(data), "\n", "", -1))
	if err != nil {
		return false, err
	}
	return i > 0, nil
}
