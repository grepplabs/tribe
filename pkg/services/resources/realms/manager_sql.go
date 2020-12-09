package realms

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

var ErrRealmParamMissing = errors.New("realm param is missing")

type sqlManager struct {
	DBS db.Session
}

func NewSQLManager(dbs db.Session) Manager {
	return &sqlManager{
		DBS: dbs,
	}
}

func (m sqlManager) CreateRealm(ctx context.Context, r *models.Realm) error {
	if r == nil {
		return ErrRealmParamMissing
	}
	_, err := m.DBS.WithContext(ctx).Collection(r.TableName()).Insert(r)
	if err != nil {
		return err
	}
	return nil
}
