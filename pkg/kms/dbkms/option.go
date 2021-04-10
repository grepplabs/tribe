package dbkms

import (
	"github.com/grepplabs/tribe/config"
	dbClient "github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/log"
)

type Option func(client *client) error

func WithDBClient(dbClient dbClient.Client) Option {
	return func(c *client) error {
		c.dbClient = dbClient
		return nil
	}
}

func WithLogger(logger log.Logger) Option {
	return func(c *client) error {
		c.logger = logger
		return nil
	}
}

func WithDBConfig(config *config.DBConfig) Option {
	return func(c *client) error {
		c.dbConfig = config
		return nil
	}
}

func WithMasterSecret(masterSecret string) Option {
	return func(c *client) error {
		c.masterSecret = masterSecret
		return nil
	}
}
