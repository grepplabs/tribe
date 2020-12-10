package users

import (
	"context"
	"github.com/grepplabs/tribe/database/models"
	"github.com/grepplabs/tribe/pkg"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type sqlManager struct {
	DBS db.Session
}

func NewSQLManager(dbs db.Session) Manager {
	return &sqlManager{
		DBS: dbs,
	}
}

func (m sqlManager) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return pkg.ErrIllegalArgument{Reason: "Input parameter user is missing"}
	}
	_, err := m.DBS.WithContext(ctx).Collection(user.TableName()).Insert(user)
	if err != nil {
		return err
	}
	return nil
}
func (m sqlManager) GetUser(ctx context.Context, realmID string, username string) (*models.User, error) {
	if realmID == "" || username == "" {
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameter realmID and username must not be empty"}
	}
	var user models.User
	err := m.DBS.WithContext(ctx).Collection(user.TableName()).Find(db.Cond{"realm_id": realmID, "username": username}).One(&user)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find user")
	}
	return &user, nil
}
