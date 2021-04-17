package clientsql

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/database/service"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type oidcJwksManager struct {
	dbs db.Session
}

func (m oidcJwksManager) CreateOidcJWKS(ctx context.Context, id *model.OidcJWKS) error {
	if id == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter id for OidcJWKS is missing"}
	}
	_, err := m.dbs.WithContext(ctx).Collection(id.TableName()).Insert(id)
	if err != nil {
		return errors.Wrap(err, "insert OidcJWKS")
	}
	return nil
}

func (m oidcJwksManager) GetOidcJWKS(ctx context.Context, id string) (*model.OidcJWKS, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var record model.OidcJWKS
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": id}).One(&record)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find OidcJWKS")
	}
	return &record, nil
}

func (m oidcJwksManager) DeleteOidcJWKS(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var record model.OidcJWKS
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": id}).Delete()
	return errors.Wrap(err, "delete OidcJWKS")
}

func (m oidcJwksManager) UpdateOidcJWKS(ctx context.Context, record *model.OidcJWKS) error {
	if record == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter record is missing"}
	}
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": record.ID}).Update(record)
	return errors.Wrap(err, "update OidcJWKS")
}
