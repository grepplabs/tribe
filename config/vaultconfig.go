package config

import (
	"crypto/tls"
	tlsconfig "github.com/grepplabs/tribe/pkg/tls"
	"github.com/spf13/pflag"
)

type VaultConfig struct {
	flagBase

	Address string
	Token   string

	TLSConfig VaultTLSConfig
}

type VaultTLSConfig struct {
	Cert               string
	Key                string
	CaCert             string
	ServerName         string
	UseSystemCertPool  bool
	InsecureSkipVerify bool
}

func NewVaultConfig() *VaultConfig {
	return &VaultConfig{}
}

func (c *VaultConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Address, "vault-addr", "https://localhost:8201", "Vault server address")
		c.flagSet.StringVar(&c.Token, "vault-token", "tribe-root-token", "Vault authentication token")
		c.flagSet.StringVar(&c.TLSConfig.Cert, "vault-tls-cert", "", "Client cert file")
		c.flagSet.StringVar(&c.TLSConfig.Key, "vault-tls-key", "", "Client key file")
		c.flagSet.StringVar(&c.TLSConfig.CaCert, "vault-tls-ca-cert", "", "CA cert file")
		c.flagSet.StringVar(&c.TLSConfig.ServerName, "vault-tls-servername", "", "Server name to verify the hostname on the returned certificates")
		c.flagSet.BoolVar(&c.TLSConfig.UseSystemCertPool, "vault-tls-use-system-cert-pool", true, "Use system cert pool which system default locations can be overridden with SSL_CERT_FILE/SSL_CERT_DIR env")
		c.flagSet.BoolVar(&c.TLSConfig.InsecureSkipVerify, "vault-tls-insecure-skip-verify", false, "Disables SSL certificate verification")
	}
	return c.flagSet
}

func (c *VaultTLSConfig) NewClientConfig() (*tls.Config, error) {
	return tlsconfig.NewClientConfig(c.Cert, c.Key, c.CaCert, c.ServerName, c.UseSystemCertPool, c.InsecureSkipVerify)
}
