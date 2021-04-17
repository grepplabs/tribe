package cmd

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/jwk"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
	"os"
)

func init() {
	oidcJwksCmd.AddCommand(newOidcjwksGetCmd())
}

type oidcjwksGetConfig struct {
	oidcJwksID string
}

func (c *oidcjwksGetConfig) Validate() error {
	return nil
}

func newOidcjwksGetCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	kmsConfig := config.NewKMSConfig(datastoreConfig)
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(oidcjwksGetConfig)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get OIDC JWKS public keys",
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

			logger := log.NewLogger(logConfig.Configuration).WithName("oidc-jwks-get")
			result, err := runOidcjwksGet(logger, datastoreConfig, kmsConfig, cmdConfig)
			if err != nil {
				log.Errorf("oidc jwks get command failed: %v", err)
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

	cmd.Flags().StringVar(&cmdConfig.oidcJwksID, "oidc-jwks-id", "", "Identifier of the oidc jwks")
	_ = cmd.MarkFlagRequired("oidc-jwks-id")

	return cmd
}

func runOidcjwksGet(logger log.Logger, datastoreConfig *config.DatastoreConfig, kmsConfig *config.KMSConfig, cmdConfig *oidcjwksGetConfig) (interface{}, error) {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return nil, err
	}
	kmsProvider, err := NewKMSProvider(logger, kmsConfig)
	if err != nil {
		return nil, err
	}
	record, err := dsClient.API().GetOidcJWKS(context.Background(), cmdConfig.oidcJwksID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.Errorf("OIDC jwksID not found: %s", cmdConfig.oidcJwksID)
	}

	jwksGetCmd := NewJwksGetCmd(logger, dsClient, kmsProvider)
	result := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0),
	}
	if err := appendPublicKey(record.CurrentJwksID, jwksGetCmd, result); err != nil {
		return nil, err
	}
	if err := appendPublicKey(record.NextJwksID, jwksGetCmd, result); err != nil {
		return nil, err
	}
	if record.PreviousJwksID != nil {
		if err := appendPublicKey(*record.PreviousJwksID, jwksGetCmd, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}
func appendPublicKey(jwksID string, jwksGetCmd *jwksCreateGet, result *jose.JSONWebKeySet) error {
	jsoNWebKeySet, err := jwksGetCmd.Run(&jwksGetConfig{
		jwksID: jwksID,
	})
	if err != nil {
		return err
	}
	if len(jsoNWebKeySet.Keys) != 2 {
		return errors.Errorf("JWKS ID %s should have private and public key: %v", jwksID, len(jsoNWebKeySet.Keys))
	}
	if jwk.IsPublic(&jsoNWebKeySet.Keys[0]) && jwk.IsPrivate(&jsoNWebKeySet.Keys[1]) {
		result.Keys = append(result.Keys, jsoNWebKeySet.Keys[0])
	} else if jwk.IsPublic(&jsoNWebKeySet.Keys[1]) && jwk.IsPrivate(&jsoNWebKeySet.Keys[0]) {
		result.Keys = append(result.Keys, jsoNWebKeySet.Keys[1])
	} else {
		return errors.Errorf("JWKS ID %s with invalid private and public keys set", jwksID)
	}
	return nil
}
