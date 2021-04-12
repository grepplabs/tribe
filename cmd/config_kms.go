package cmd

import (
	"fmt"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/hcvault"
	"github.com/google/tink/go/tink"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/kms/dbkms"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

const (
	vaultRefKeyURIPrefix = "hcvault://vault"
)

func registerKMSClient(logger log.Logger, kmsConfig *config.KMSConfig, masterSecret string) error {
	switch kmsConfig.Provider {
	case "db":
		dsClient, err := NewDatastoreClient(logger, kmsConfig.DatastoreConfig)
		if err != nil {
			return err
		}
		err = dbkms.RegisterKMSClient(logger, dsClient, masterSecret, kmsConfig.DatastoreConfig.Provider)
		if err != nil {
			return err
		}
	case "vault", "hcvault":
		vurl, err := url.Parse(kmsConfig.VaultConfig.Address)
		if err != nil {
			return err
		}
		if vurl.Scheme != "" && vurl.Scheme != "https" {
			return errors.Errorf("vault address with https schema expected, but got %s", vurl.Scheme)
		}
		if vurl.Host == "" {
			return errors.Errorf("vault address is empty, address %s", kmsConfig.VaultConfig.Address)
		}
		tlsConfig, err := kmsConfig.VaultConfig.TLSConfig.NewClientConfig()
		if err != nil {
			return err
		}
		keyURI := fmt.Sprintf("hcvault://%s", vurl.Host)
		vaultClient, err := hcvault.NewClient(keyURI, tlsConfig, kmsConfig.VaultConfig.Token)
		registry.RegisterKMSClient(vaultClient)
	}
	return nil
}

type kmsProvider struct {
	logger       log.Logger
	kmsConfig    *config.KMSConfig
	masterSecret string
}

type KMSProvider interface {
	AEADFromKeyURI(refKeyURI string) (aead tink.AEAD, err error)
	NewAEAD(jwksID string) (aead tink.AEAD, refKeyURI string, err error)
}

func NewKMSProvider(logger log.Logger, kmsConfig *config.KMSConfig, masterSecret string) (KMSProvider, error) {
	err := registerKMSClient(logger, kmsConfig, masterSecret)
	if err != nil {
		return nil, err
	}
	return &kmsProvider{
		logger:       logger,
		kmsConfig:    kmsConfig,
		masterSecret: masterSecret,
	}, nil
}

func (p kmsProvider) AEADFromKeyURI(refKeyURI string) (aead tink.AEAD, err error) {
	var keyURI = ""
	switch p.kmsConfig.Provider {
	case "db":
		keyURI = refKeyURI
	case "vault", "hcvault":
		vurl, err := url.Parse(p.kmsConfig.VaultConfig.Address)
		if err != nil {
			return nil, err
		}
		if vurl.Host == "" {
			return nil, errors.Errorf("vault address is empty, address %s", p.kmsConfig.VaultConfig.Address)
		}
		keyURI = strings.Replace(refKeyURI, vaultRefKeyURIPrefix, fmt.Sprintf("hcvault://%s", vurl.Host), 1)
	default:
		return nil, errors.Errorf("unsupported kms provider %s", p.kmsConfig.Provider)
	}
	kmsClient, err := registry.GetKMSClient(keyURI)
	if err != nil {
		return nil, err
	}
	aead, err = kmsClient.GetAEAD(keyURI)
	if err != nil {
		return nil, err
	}
	return aead, nil
}

func (p kmsProvider) NewAEAD(jwksID string) (aead tink.AEAD, refKeyURI string, err error) {
	switch p.kmsConfig.Provider {
	case "db":
		keyURI := fmt.Sprintf("db://%s?kms-keyset-id=%s", p.kmsConfig.DatastoreConfig.Provider, p.kmsConfig.KeysetId)
		kmsClient, err := registry.GetKMSClient(keyURI)
		if err != nil {
			return nil, "", err
		}
		aead, err := kmsClient.GetAEAD(keyURI)
		if err != nil {
			return nil, "", err
		}
		return aead, keyURI, nil
	case "vault", "hcvault":
		vurl, err := url.Parse(p.kmsConfig.VaultConfig.Address)
		if err != nil {
			return nil, "", err
		}
		if vurl.Host == "" {
			return nil, "", errors.Errorf("vault address is empty, address %s", p.kmsConfig.VaultConfig.Address)
		}
		keyURI := fmt.Sprintf("hcvault://%s/transit/keys/tribe-jwks-%s", vurl.Host, jwksID)
		kmsClient, err := registry.GetKMSClient(keyURI)
		if err != nil {
			return nil, "", err
		}
		aead, err := kmsClient.GetAEAD(keyURI)
		if err != nil {
			return nil, "", err
		}
		return aead, strings.Replace(keyURI, fmt.Sprintf("hcvault://%s", vurl.Host), vaultRefKeyURIPrefix, 1), nil
	default:
		return nil, "", errors.Errorf("unsupported kms provider %s", p.kmsConfig.Provider)
	}
}
