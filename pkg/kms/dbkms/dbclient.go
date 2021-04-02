package dbkms

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"github.com/grepplabs/tribe/config"
	dbClient "github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/kms/masterkey"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"strings"
)

type client struct {
	logger   log.Logger
	dbClient dbClient.Client
	dbConfig *config.DBConfig
}

func (c *client) GetAEAD(keyURI string, masterSecret string) tink.AEAD {
	return NewAEAD(func() (*keyset.Handle, error) {
		// add some caching
		mk, err := c.getMasterKey(keyURI, masterSecret)
		if err != nil {
			return nil, err
		}
		return mk.GetKeyset(), nil
	})
}

func NewClient(options ...Option) (*client, error) {
	c := &client{
		logger: log.DefaultLogger.WithName("dbkms-client"),
	}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	if c.dbClient == nil {
		if c.dbConfig == nil {
			return nil, errors.New("either dbClient or dbConfig must be provided")
		}
		dbc, err := dbClient.NewSQLClient(c.logger, c.dbConfig)
		if err != nil {
			return nil, errors.Wrap(err, "create sql client failed")
		}
		c.dbClient = dbc
	}
	return c, nil
}

func (c *client) getMasterKey(keyURI string, masterSecret string) (masterkey.MasterKeyset, error) {
	if masterSecret == "" {
		return nil, errors.Errorf("master-secret is required for kms-keyset-uri: %s", keyURI)
	}
	const dbPrefix = "db://"
	if !strings.HasPrefix(strings.ToLower(keyURI), dbPrefix) {
		return nil, fmt.Errorf("uriPrefix must start with %s, but got %s", dbPrefix, keyURI)
	}
	keysetID := strings.TrimPrefix(keyURI, dbPrefix)
	ks, err := c.dbClient.API().GetKMSKeyset(context.Background(), keysetID)
	if err != nil {
		return nil, errors.Wrapf(err, "get kms keyset failed: %s", keysetID)
	}
	if ks == nil {
		return nil, errors.Errorf("kms keyset not found: %s", keysetID)
	}
	encryptedKeyset, err := base64.StdEncoding.DecodeString(ks.EncryptedKeyset)
	if err != nil {
		return nil, errors.Wrap(err, "base64 decode of encrypted keyset failed")
	}
	mk, err := masterkey.DecryptKeyset(encryptedKeyset, []byte(masterSecret))
	if err != nil {
		return nil, errors.Wrap(err, "decrypt master keyset failed")
	}
	return mk, nil
}
