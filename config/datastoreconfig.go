package config

import (
	"github.com/spf13/pflag"
)

type DatastoreConfig struct {
	flagBase

	Provider    string
	DBConfig    DBConfig
	MinioConfig MinioConfig
}

func NewDatastoreConfig() *DatastoreConfig {
	return &DatastoreConfig{}
}

func (c *DatastoreConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Provider, "datastore-provider", "db", "Datastore provider. One of: [db, minio]")
	}
	c.flagSet.AddFlagSet(c.DBConfig.FlagSet())
	c.flagSet.AddFlagSet(c.MinioConfig.FlagSet())
	return c.flagSet
}
