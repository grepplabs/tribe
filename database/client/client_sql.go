package client

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/services/resources/realms"
	"github.com/grepplabs/tribe/pkg/services/resources/users"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"net/url"
)

type sqlClient struct {
	DBS db.Session
}

func NewSQLClient(config config.DBConfig) (Client, error) {
	dbs, err := newDBSession(config)
	if err != nil {
		return nil, err
	}
	return &sqlClient{
		DBS: dbs,
	}, nil
}

func (c sqlClient) UserManager() users.Manager {
	return users.NewSQLManager(c.DBS)
}

func (c sqlClient) RealmManager() realms.Manager {
	return realms.NewSQLManager(c.DBS)
}

func newDBSession(config config.DBConfig) (db.Session, error) {
	//TODO: dbx with tracing and profiling ()

	dsnUrl, err := url.Parse(config.ConnectionURL)
	if err != nil {
		return nil, err
	}
	// TODO: this is postgres only so far
	if dsnUrl.Scheme != postgresql.Adapter {
		return nil, errors.Errorf("'%s' DNS schema expected, but got '%s'", postgresql.Adapter, dsnUrl.Scheme)
	}
	//TODO: retry on connection timeout
	ctx := context.Background()
	dbx, err := sqlx.ConnectContext(ctx, "postgres", config.ConnectionURL)
	if err != nil {
		return nil, err
	}
	dbx.SetMaxIdleConns(config.MaxIdleConns)
	dbx.SetMaxOpenConns(config.MaxOpenConns)
	dbx.SetConnMaxLifetime(config.ConnMaxLifetime)

	return postgresql.New(dbx.DB)
}
