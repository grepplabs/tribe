package client

import (
	"context"
	"github.com/grepplabs/tribe/pkg/log"
	"net/url"
	"time"

	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/service"
	"github.com/grepplabs/tribe/database/service/clientsql"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/strategy"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type sqlClient struct {
	dbs db.Session
	api service.API
}

func NewSQLClient(logger log.Logger, config *config.DBConfig) (Client, error) {
	dbs, err := newDBSession(logger, config)
	if err != nil {
		return nil, err
	}
	clientsql.NewAPIImpl(dbs)
	return &sqlClient{
		dbs: dbs,
		api: clientsql.NewAPIImpl(dbs),
	}, nil
}

func (c sqlClient) API() service.API {
	return c.api
}

func newDBSession(logger log.Logger, config *config.DBConfig) (db.Session, error) {
	//TODO: dbx with tracing and profiling ()
	dsnUrl, err := url.Parse(config.ConnectionURL)
	if err != nil {
		return nil, err
	}
	// TODO: this is postgres only so far
	if dsnUrl.Scheme != postgresql.Adapter {
		return nil, errors.Errorf("'%s' DNS schema expected, but got '%s'", postgresql.Adapter, dsnUrl.Scheme)
	}

	dbx, err := connectWithRetry(logger, "postgres", config.ConnectionURL)
	if err != nil {
		return nil, err
	}
	dbx.SetMaxIdleConns(config.MaxIdleConns)
	dbx.SetMaxOpenConns(config.MaxOpenConns)
	dbx.SetConnMaxLifetime(config.ConnMaxLifetime)

	return postgresql.New(dbx.DB)
}

func connectWithRetry(logger log.Logger, driverName, dataSourceName string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var dbx *sqlx.DB
	action := func(ctx context.Context) (err error) {
		dbx, err = sqlx.ConnectContext(ctx, driverName, dataSourceName)
		if err != nil {
			logger.Warnf("connection to database failed: %v", err)
		}
		return err
	}
	how := retry.How{
		strategy.Limit(10),
		strategy.Wait(1 * time.Second),
	}
	err := retry.Do(ctx, action, how...)
	return dbx, err
}
