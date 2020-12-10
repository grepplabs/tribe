package realms

import (
	"context"
	"fmt"
	"github.com/grepplabs/tribe/database/models"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type ErrIllegalArgument struct {
	reason string
}

func (e ErrIllegalArgument) Error() string {
	return fmt.Sprintf("Illegal argument: %q", e.reason)
}

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
		return ErrIllegalArgument{reason: "Input parameter realm is missing"}
	}
	_, err := m.DBS.WithContext(ctx).Collection(realm.TableName()).Insert(realm)
	if err != nil {
		return errors.Wrap(err, "insert realm")
	}
	return nil
}

func (m sqlManager) GetRealm(ctx context.Context, realmID string) (*models.Realm, error) {
	if realmID == "" {
		return nil, ErrIllegalArgument{reason: "Input parameter realmID is missing"}
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
