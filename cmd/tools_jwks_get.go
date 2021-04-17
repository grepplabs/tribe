package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
	"os"
)

func init() {
	jwksCmd.AddCommand(newJwksGetCmd())
}

type jwksGetConfig struct {
	jwksID string

	use string
	kid string
}

func (c *jwksGetConfig) Validate() error {
	if c.kid == "" && c.jwksID == "" {
		return errors.New("either jwks-id or kid is required")
	}
	return nil
}

func newJwksGetCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	kmsConfig := config.NewKMSConfig(datastoreConfig)
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(jwksGetConfig)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get JWKS",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdConfig.Validate(); err != nil {
				return err
			}
			if err := outputConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			producer := outputConfig.MustGetProducer()

			logger := log.NewLogger(logConfig.Configuration).WithName("jwks-get")
			result, err := runJwksGet(logger, datastoreConfig, kmsConfig, cmdConfig)
			if err != nil {
				log.Errorf("jwks get command failed: %v", err)
				os.Exit(1)
			}
			if result == nil {
				log.Errorf("jwks get command failed, not found jwksID %s", cmdConfig.jwksID)
				os.Exit(1)
			}
			err = producer.Produce(os.Stdout, result)
			if err != nil {
				log.Errorf("failed to write result: %v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())
	cmd.Flags().AddFlagSet(kmsConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.jwksID, "jwks-id", "", "Identifier of the jwks, JWKSID")
	cmd.Flags().StringVar(&cmdConfig.use, "use", "sig", "How the key is meant to be used. One of: [sig, enc]")
	cmd.Flags().StringVar(&cmdConfig.kid, "kid", "", "Unique key identifier. The Key ID is generated if not specified.")

	return cmd
}

func runJwksGet(logger log.Logger, datastoreConfig *config.DatastoreConfig, kmsConfig *config.KMSConfig, cmdConfig *jwksGetConfig) (interface{}, error) {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return nil, err
	}
	jwks, err := getJwks(dsClient, cmdConfig)
	if err != nil {
		return nil, err
	}
	encryptedKeys, err := base64.StdEncoding.DecodeString(jwks.EncryptedJwks)
	if err != nil {
		return nil, errors.Wrapf(err, "base64 decode of JWKS ID failed: %s", cmdConfig.jwksID)
	}
	kmsProvider, err := NewKMSProvider(logger, kmsConfig)
	if err != nil {
		return nil, err
	}
	aead, err := kmsProvider.AEADFromKeyURI(jwks.KMSKeyURI)
	if err != nil {
		return nil, err
	}
	decryptedKeys, err := aead.Decrypt(encryptedKeys, []byte{})
	if err != nil {
		return nil, errors.Wrap(err, "AEAD keys decryption failed")
	}
	var result jose.JSONWebKeySet
	err = json.Unmarshal(decryptedKeys, &result)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal JSONWebKeySet failed")
	}
	return &result, nil
}

func getJwks(dsClient client.Client, cmdConfig *jwksGetConfig) (*model.JWKS, error) {
	if cmdConfig.jwksID != "" {
		return getJwksByID(dsClient, cmdConfig.jwksID)
	} else {
		jwks, err := dsClient.API().GetJWKSByKidUse(context.Background(), cmdConfig.kid, cmdConfig.use)
		if err != nil {
			return nil, err
		}
		if jwks == nil {
			return nil, errors.Errorf("not found JWKS kid, sig: %s, %s", cmdConfig.kid, cmdConfig.use)
		}
		return jwks, nil
	}
}

func getJwksByID(dsClient client.Client, jwksID string) (*model.JWKS, error) {
	jwks, err := dsClient.API().GetJWKS(context.Background(), jwksID)
	if err != nil {
		return nil, err
	}
	if jwks == nil {
		return nil, errors.Errorf("not found JWKS ID: %s", jwksID)
	}
	return jwks, nil
}
