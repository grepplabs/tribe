package clientsql

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/database/service"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type jwksManager struct {
	dbs db.Session
}

func (m jwksManager) CreateJWKS(ctx context.Context, jwks *model.JWKS) error {
	if jwks == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter jwks is missing"}
	}
	_, err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Insert(jwks)
	if err != nil {
		return errors.Wrap(err, "insert jwks")
	}
	return nil
}

func (m jwksManager) GetJWKS(ctx context.Context, id string) (*model.JWKS, error) {
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

func (m jwksManager) GetJWKSByKidUse(ctx context.Context, kid string, use string) (*model.JWKS, error) {
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

func (m jwksManager) DeleteJWKS(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"id": id}).Delete()
	return errors.Wrap(err, "delete id")
}

func (m jwksManager) DeleteJWKSByKidUse(ctx context.Context, kid string, use string) error {
	if kid == "" || use == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter kid/use is missing"}
	}
	var jwks model.JWKS
	err := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find(db.Cond{"kid": kid, "use": use}).Delete()
	return errors.Wrap(err, "delete kid/use")
}

func (m jwksManager) ListJWKS(ctx context.Context, offset *int64, limit *int64) (*model.JWKSList, error) {
	var jwks model.JWKS
	result := m.dbs.WithContext(ctx).Collection(jwks.TableName()).Find().OrderBy("created_at")
	if offset != nil && *offset > 0 {
		result = result.Offset(int(*offset))
	}
	if limit != nil && *limit > 0 {
		// limit 0 all elements
		result = result.Limit(int(*limit))
	}

	var list []model.JWKS
	err := result.All(&list)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "list jwks")
	}
	// this executes additional query
	total, err := result.TotalEntries()
	if err != nil {
		return nil, errors.Wrap(err, "list realms total entries")
	}
	return &model.JWKSList{List: list, Page: model.Page{
		Offset: offset,
		Limit:  limit,
		Total:  total,
	}}, nil
}
