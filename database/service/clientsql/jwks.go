package clientsql

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/database/service"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type JWKSManager struct {
	dbs db.Session
}

func (m JWKSManager) CreateJWKS(ctx context.Context, jwks *model.JWKS) error {
	if jwks == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter jwks is missing"}
	}
	_, err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Insert(jwks)
	if err != nil {
		return errors.Wrap(err, "insert jwks")
	}
	return nil
}

func (m JWKSManager) GetJWKS(ctx context.Context, id string) (*model.JWKS, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"id": id}).One(&jwks)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find jwks")
	}
	return &jwks, nil
}

func (m JWKSManager) GetJWKSByKidUse(ctx context.Context, kid string, use string) (*model.JWKS, error) {
	if kid == "" || use == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter kid/use is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"kid": kid, "use": use}).One(&jwks)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find jwks by kid and use")
	}
	return &jwks, nil
}

func (m JWKSManager) DeleteJWKS(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"id": id}).Delete()
	return errors.Wrap(err, "delete id")
}

func (m JWKSManager) DeleteJWKSByKidUse(ctx context.Context, kid string, use string) error {
	if kid == "" || use == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter kid/use is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"kid": kid, "use": use}).Delete()
	return errors.Wrap(err, "delete kid/use")
}
