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
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameters realmID and username must not be empty"}
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

func (m sqlManager) ListUsers(ctx context.Context, realmID string, offset *int64, limit *int64) (*models.UserList, error) {
	if realmID == "" {
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameter realmID must not be empty"}
	}
	if limit != nil && int(*limit) < 0 {
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameter limit must not be negative"}
	}

	var user models.User
	result := m.DBS.WithContext(ctx).Collection(user.TableName()).Find(db.Cond{"realm_id": realmID}).OrderBy("realm_id", "created_at")
	if offset != nil {
		result = result.Offset(int(*offset))
	}
	if limit != nil && int(*limit) > 0 {
		// limit 0 all elements
		result = result.Limit(int(*limit))
	}

	var users []models.User
	err := result.All(&users)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "list users")
	}
	// this executes additional query
	total, err := result.TotalEntries()
	if err != nil {
		return nil, errors.Wrap(err, "list users total entries")
	}

	return &models.UserList{Users: users, Page: models.Page{
		Offset: offset,
		Limit:  offset,
		Total:  total,
	}}, nil
}

func (m sqlManager) ExistsUser(ctx context.Context, realmID string, username string) (bool, error) {
	if realmID == "" || username == "" {
		return false, pkg.ErrIllegalArgument{Reason: "Input parameters realmID and username must not be empty"}
	}
	var user models.User
	exists, err := m.DBS.WithContext(ctx).Collection(user.TableName()).Find(db.Cond{"realm_id": realmID, "username": username}).Exists()
	return exists, errors.Wrap(err, "exists user")
}

func (m sqlManager) DeleteUser(ctx context.Context, realmID string, username string) error {
	if realmID == "" || username == "" {
		return pkg.ErrIllegalArgument{Reason: "Input parameters realmID and username must not be empty"}
	}
	var user models.User
	err := m.DBS.WithContext(ctx).Collection(user.TableName()).Find(db.Cond{"realm_id": realmID, "username": username}).Delete()
	return errors.Wrap(err, "delete user")
}
