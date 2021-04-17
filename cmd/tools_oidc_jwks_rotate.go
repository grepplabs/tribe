package cmd

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/grepplabs/tribe/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	oidcJwksCmd.AddCommand(newOidcJwksRotateCmd())
}

type oidcJwksRotateConfig struct {
	oidcJwksID string

	nextJwksID string
	alg        string

	currentJwksID string
	revoke        bool
}

func (c *oidcJwksRotateConfig) Validate() error {
	return nil
}

func newOidcJwksRotateCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	kmsConfig := config.NewKMSConfig(datastoreConfig)
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(oidcJwksRotateConfig)

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate OIDC JWKS",
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

			logger := log.NewLogger(logConfig.Configuration).WithName("oidc-jwks-rotate")
			result, err := runOidcJwksRotate(logger, datastoreConfig, kmsConfig, cmdConfig)
			if err != nil {
				log.Errorf("oidc jwks rotate command failed: %v", err)
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

	cmd.Flags().BoolVar(&cmdConfig.revoke, "revoke", false, "Rotate the currently used key")
	cmd.Flags().StringVar(&cmdConfig.oidcJwksID, "oidc-jwks-id", "", "Identifier of the oidc jwks")
	cmd.Flags().StringVar(&cmdConfig.nextJwksID, "next-jwks-id", "", "Next JWKS ID to use")
	cmd.Flags().StringVar(&cmdConfig.currentJwksID, "current-jwks-id", "", "Current JWKS ID used with revoke option")
	cmd.Flags().StringVar(&cmdConfig.alg, "alg", "RS256", "The specific asymmetric rfc7518 JWA algorithm to be used to generated the key. One of: [RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512]")

	_ = cmd.MarkFlagRequired("oidc-jwks-id")

	return cmd
}

func runOidcJwksRotate(logger log.Logger, datastoreConfig *config.DatastoreConfig, kmsConfig *config.KMSConfig, cmdConfig *oidcJwksRotateConfig) (interface{}, error) {
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
	err = validateRotateJwksIDs(cmdConfig, record)
	if err != nil {
		return nil, err
	}
	jwksCreate := NewJwksCreateCmd(logger, dsClient, kmsProvider)
	nextJwksID, currentJwksID, previousJwksID, err := rotateJwksIDs(jwksCreate, dsClient, cmdConfig, record)
	if err != nil {
		return nil, err
	}

	record.CurrentJwksID = currentJwksID
	record.NextJwksID = nextJwksID
	record.PreviousJwksID = previousJwksID
	record.LastRotated = time.Now()
	record.Version = record.Version + 1

	err = dsClient.API().UpdateOidcJWKS(context.Background(), record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func validateRotateJwksIDs(cmdConfig *oidcJwksRotateConfig, record *model.OidcJWKS) error {
	if cmdConfig.nextJwksID != "" && cmdConfig.currentJwksID == cmdConfig.nextJwksID {
		return errors.Errorf("OIDC jwksID must be different: next %s , current %s", cmdConfig.nextJwksID, cmdConfig.currentJwksID)
	}
	if cmdConfig.nextJwksID != "" && (cmdConfig.nextJwksID == record.NextJwksID || cmdConfig.nextJwksID == record.CurrentJwksID || cmdConfig.nextJwksID == utils.StringValue(record.PreviousJwksID)) {
		return errors.Errorf("Next OIDC jwksID must be different from stored keys : next %s , current %s, previous %s", cmdConfig.nextJwksID, record.NextJwksID, record.CurrentJwksID, utils.StringValue(record.PreviousJwksID))
	}
	if cmdConfig.currentJwksID != "" && (cmdConfig.currentJwksID == record.NextJwksID || cmdConfig.currentJwksID == record.CurrentJwksID || cmdConfig.currentJwksID == utils.StringValue(record.PreviousJwksID)) {
		return errors.Errorf("Current OIDC jwksID must be different from stored keys : next %s , current %s, previous %s", cmdConfig.currentJwksID, record.NextJwksID, record.CurrentJwksID, utils.StringValue(record.PreviousJwksID))
	}
	return nil
}

func rotateJwksIDs(jwksCreate *jwksCreateCmd, dsClient client.Client, cmdConfig *oidcJwksRotateConfig, record *model.OidcJWKS) (string, string, *string, error) {
	if cmdConfig.revoke {
		currentJwksID, err := oidcJwksCreateOrGet(cmdConfig.currentJwksID, cmdConfig.alg, jwksCreate, dsClient)
		if err != nil {
			return "", "", nil, err
		}
		nextJwksID, err := oidcJwksCreateOrGet(cmdConfig.nextJwksID, cmdConfig.alg, jwksCreate, dsClient)
		if err != nil {
			return "", "", nil, err
		}
		return nextJwksID, currentJwksID, nil, nil
	} else {
		nextJwksID, err := oidcJwksCreateOrGet(cmdConfig.nextJwksID, cmdConfig.alg, jwksCreate, dsClient)
		if err != nil {
			return "", "", nil, err
		}
		currentJwksID := record.NextJwksID
		previousJwksID := record.CurrentJwksID
		return nextJwksID, currentJwksID, utils.String(previousJwksID), nil
	}
}
