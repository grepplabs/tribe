package config

import "time"

type DBConfig struct {
	ConnectionURL   string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}
