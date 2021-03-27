package config

import (
	"github.com/spf13/pflag"
)

type PaginationConfig struct {
	flagBase
	Limit  int64
	Offset int64
}

func NewPaginationConfig() *PaginationConfig {
	return &PaginationConfig{}
}

func (c *PaginationConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.Int64Var(&c.Limit, "limit", 0, "The numbers of entries to return")
		c.flagSet.Int64Var(&c.Offset, "offset", 0, "The number of items to skip before starting to collect the result set")
	}
	return c.flagSet
}
