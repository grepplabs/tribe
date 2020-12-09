package users

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

var ErrUserParamMissing = errors.New("user param is missing")

type sqlManager struct {
	DBS db.Session
}

func NewSQLManager(dbs db.Session) Manager {
	return &sqlManager{
		DBS: dbs,
	}
}

func (m sqlManager) CreateUser(ctx context.Context, u *models.User) error {
	if u == nil {
		return ErrUserParamMissing
	}
	_, err := m.DBS.WithContext(ctx).Collection(u.TableName()).Insert(u)
	if err != nil {
		return err
	}
	return nil
}
