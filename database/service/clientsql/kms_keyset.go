package clientsql

import (
	"context"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/database/service"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type kmsKeysetManager struct {
	dbs db.Session
}

func (m kmsKeysetManager) CreateKMSKeyset(ctx context.Context, kmsKeyset *model.KMSKeyset) error {
	if kmsKeyset == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter kmsKeyset is missing"}
	}
	_, err := m.dbs.WithContext(ctx).Collection(kmsKeyset.TableName()).Insert(kmsKeyset)
	if err != nil {
		return errors.Wrap(err, "insert kmsKeyset")
	}
	return nil
}
func (m kmsKeysetManager) GetKMSKeyset(ctx context.Context, id string) (*model.KMSKeyset, error) {
	if id == "" {
		return nil, service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var record model.KMSKeyset
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": id}).One(&record)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find KMSKeyset")
	}
	return &record, nil
}

func (m kmsKeysetManager) DeleteKMSKeyset(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrIllegalArgument{Reason: "Input parameter id is missing"}
	}
	var record model.KMSKeyset
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": id}).Delete()
	return errors.Wrap(err, "delete KMSKeyset")
}

func (m kmsKeysetManager) UpdateKMSKeyset(ctx context.Context, record *model.KMSKeyset) error {
	if record == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter record is missing"}
	}
	err := m.dbs.WithContext(ctx).Collection(record.TableName()).Find(db.Cond{"id": record.ID}).Update(record)
	return errors.Wrap(err, "update KMSKeyset")
}

func (m kmsKeysetManager) ListKMSKeysets(ctx context.Context, offset *int64, limit *int64) (*model.KMSKeysetList, error) {
	var kmsKeyset model.KMSKeyset
	result := m.dbs.WithContext(ctx).Collection(kmsKeyset.TableName()).Find().OrderBy("created_at")
	if offset != nil && *offset > 0 {
		result = result.Offset(int(*offset))
	}
	if limit != nil && *limit > 0 {
		// limit 0 all elements
		result = result.Limit(int(*limit))
	}

	var list []model.KMSKeyset
	err := result.All(&list)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "list kmsKeyset")
	}
	// this executes additional query
	total, err := result.TotalEntries()
	if err != nil {
		return nil, errors.Wrap(err, "list realms total entries")
	}
	return &model.KMSKeysetList{List: list, Page: model.Page{
		Offset: offset,
		Limit:  limit,
		Total:  total,
	}}, nil
}
