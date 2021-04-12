package dbkms

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"github.com/grepplabs/tribe/config"
	dbClient "github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/kms/masterkey"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

const (
	dbPrefix       = "db://"
	keyKmsKeysetId = "kms-keyset-id"
)

var _ registry.KMSClient = (*client)(nil)

type client struct {
	keyURIPrefix string
	masterSecret string
	logger       log.Logger
	dbClient     dbClient.Client
	dbConfig     *config.DBConfig
}

// Implements KMSClient Supported methods
func (c *client) Supported(keyURI string) bool {
	return strings.HasPrefix(keyURI, c.keyURIPrefix)
}

// Implements KMSClient GetAEAD methods
func (c *client) GetAEAD(keyURI string) (tink.AEAD, error) {
	mk, err := c.getMasterKey(keyURI, c.masterSecret)
	if err != nil {
		return nil, err
	}
	return NewAEAD(func() (*keyset.Handle, error) {
		return mk.GetKeyset(), nil
	}), nil
}

func NewClient(options ...Option) (*client, error) {
	c := &client{
		keyURIPrefix: dbPrefix,
		logger:       log.DefaultLogger.WithName("dbkms-client"),
	}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	if !strings.HasPrefix(strings.ToLower(c.keyURIPrefix), dbPrefix) {
		return nil, fmt.Errorf("keyURIPrefix must start with %s, but got %s", dbPrefix, c.keyURIPrefix)
	}
	if c.masterSecret == "" {
		return nil, errors.New("masterSecret must not be empty")
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
		return nil, errors.Errorf("master-secret is required for key-uri: %s", keyURI)
	}
	if !strings.HasPrefix(strings.ToLower(keyURI), c.keyURIPrefix) {
		return nil, fmt.Errorf("uriPrefix must start with %s, but got %s", c.keyURIPrefix, keyURI)
	}
	u, err := url.Parse(keyURI)
	if err != nil {
		return nil, errors.Wrapf(err, "url parse failed: %s", keyURI)
	}
	keysetID := u.Query().Get(keyKmsKeysetId)
	if keysetID == "" {
		return nil, errors.Errorf("query param %s not found: %s", keyKmsKeysetId, keysetID)
	}
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

func RegisterKMSClient(logger log.Logger, dsClient dbClient.Client, masterSecret string, provider string) error {
	dbkmsClient, err := NewClient(WithMasterSecret(masterSecret), WithLogger(logger), WithDBClient(dsClient), WithKeyURIPrefix(fmt.Sprintf("%s%s", dbPrefix, provider)))
	if err != nil {
		return err
	}
	registry.RegisterKMSClient(dbkmsClient)
	return nil
}
