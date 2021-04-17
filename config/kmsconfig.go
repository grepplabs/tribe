package config

import (
	"github.com/spf13/pflag"
)

type KMSConfig struct {
	flagBase

	Provider     string
	KeysetId     string
	MasterSecret string

	DatastoreConfig *DatastoreConfig
	VaultConfig     *VaultConfig
}

func NewKMSConfig(datastoreConfigd *DatastoreConfig) *KMSConfig {
	return &KMSConfig{
		DatastoreConfig: datastoreConfigd,
		VaultConfig:     NewVaultConfig(),
	}
}

func (c *KMSConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Provider, "kms-provider", "db", "KMS provider. One of: [db, vault]")
		c.flagSet.StringVar(&c.KeysetId, "kms-keyset-id", "", "Identifier of the keyset")
		c.flagSet.StringVar(&c.MasterSecret, "kms-master-secret", "", "Master secret")
	}
	c.flagSet.AddFlagSet(c.DatastoreConfig.FlagSet())
	c.flagSet.AddFlagSet(c.VaultConfig.FlagSet())
	return c.flagSet
}
