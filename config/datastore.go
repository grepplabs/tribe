package config

import (
	"github.com/spf13/pflag"
)

type DatastoreConfig struct {
	flagBase

	Provider string
}

func NewDatastoreConfig() *DatastoreConfig {
	return &DatastoreConfig{}
}

func (c *DatastoreConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Provider, "datastore-provider", "db", "Datastore provider. One of: [db, minio]")
	}
	return c.flagSet
}
