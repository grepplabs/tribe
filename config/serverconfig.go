package config

import (
	"time"
)

type ServerConfig struct {
	EnabledListeners []string
	CleanupTimeout   time.Duration
	GracefulTimeout  time.Duration
	MaxHeaderSize    int

	Host         string
	Port         int
	ListenLimit  int
	KeepAlive    time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	TLSHost           string
	TLSPort           int
	TLSCertificate    string
	TLSCertificateKey string
	TLSCACertificate  string
	TLSListenLimit    int
	TLSKeepAlive      time.Duration
	TLSReadTimeout    time.Duration
	TLSWriteTimeout   time.Duration
}
