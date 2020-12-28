package realms

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

func (m sqlManager) CreateRealm(ctx context.Context, realm *models.Realm) error {
	if realm == nil {
		return pkg.ErrIllegalArgument{Reason: "Input parameter realm is missing"}
	}
	_, err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Insert(realm)
	if err != nil {
		return errors.Wrap(err, "insert realm")
	}
	return nil
}

func (m sqlManager) GetRealm(ctx context.Context, realmID string) (*models.Realm, error) {
	if realmID == "" {
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameter realmID is missing"}
	}
	var realm models.Realm
	err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Find(db.Cond{"realm_id": realmID}).One(&realm)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find realm")
	}
	return &realm, nil
}

func (m sqlManager) ListRealms(ctx context.Context, offset *int64, limit *int64) ([]models.Realm, error) {
	if limit != nil && int(*limit) < 0 {
		return nil, pkg.ErrIllegalArgument{Reason: "Input parameter limit must not be negative"}
	}

	var realm models.Realm
	result := m.DBS.WithContext(ctx).Collection(realm.TableName()).Find().OrderBy("created_at")
	if offset != nil {
		result = result.Offset(int(*offset))
	}
	if limit != nil && int(*limit) > 0 {
		// limit 0 all elements
		result = result.Limit(int(*limit))
	}

	var realms []models.Realm
	err := result.All(&realms)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "list realms")
	}
	return realms, nil
}

func (m sqlManager) ExistsRealm(ctx context.Context, realmID string) (bool, error) {
	if realmID == "" {
		return false, pkg.ErrIllegalArgument{Reason: "Input parameter realmID is missing"}
	}
	var realm models.Realm
	exists, err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Find(db.Cond{"realm_id": realmID}).Exists()
	return exists, errors.Wrap(err, "exists realm")
}

func (m sqlManager) UpdateRealm(ctx context.Context, realm *models.Realm) error {
	if realm == nil {
		return pkg.ErrIllegalArgument{Reason: "Input parameter realm is missing"}
	}
	err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Find(db.Cond{"realm_id": realm.RealmID}).Update(realm)
	return errors.Wrap(err, "update realm")
}

func (m sqlManager) DeleteRealm(ctx context.Context, realmID string) error {
	if realmID == "" {
		return pkg.ErrIllegalArgument{Reason: "Input parameter realmID is missing"}
	}
	var realm models.Realm
	err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Find(db.Cond{"realm_id": realmID}).Delete()
	return errors.Wrap(err, "delete realm")
}
