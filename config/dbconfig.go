package config

import (
	"github.com/spf13/pflag"
	"time"
)

type DBConfig struct {
	flagBase

	ConnectionURL   string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func NewDBConfig() *DBConfig {
	return &DBConfig{}
}

func (c *DBConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.ConnectionURL, "db-connection-url", "postgresql://tribe:secret@localhost:5432/tribe?sslmode=disable", "data source name as connection URI e.g. postgresql://user:password@localhost:5432/dbname?sslmode=disable")
		c.flagSet.IntVar(&c.MaxIdleConns, "db-max-idle-conns", 2, "The maximum number of connections in the idle connection pool")
		c.flagSet.IntVar(&c.MaxOpenConns, "db-max-open-conns", 25, "The maximum number of open connections to the database")
		c.flagSet.DurationVar(&c.ConnMaxLifetime, "db-conn-max-lifetime", 0, "The maximum amount of time a connection may be reused")
	}
	return c.flagSet
}
