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

func (m oidcJwksManager) CreateOidcJWKS(ctx context.Context, oidcJWKS *model.OidcJWKS) error {
	if oidcJWKS == nil {
		return service.ErrIllegalArgument{Reason: "Input parameter oidcJWKS is missing"}
	}
	_, err := m.dbs.WithContext(ctx).Collection(oidcJWKS.TableName()).Insert(oidcJWKS)
	if err != nil {
		return errors.Wrap(err, "insert oidcJWKS")
	}
	return nil
}
