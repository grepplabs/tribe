package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	oidcJwksCmd.AddCommand(newOidcJwksCreateCmd())
}

type oidcJwksCreateConfig struct {
	oidcJwksID string

	currentJwksID string
	nextJwksID    string

	alg string
}

func (c *oidcJwksCreateConfig) Validate() error {
	return nil
}

func newOidcJwksCreateCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	kmsConfig := config.NewKMSConfig(datastoreConfig)
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(oidcJwksCreateConfig)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create OIDC JWKS",
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

			logger := log.NewLogger(logConfig.Configuration).WithName("oidc-jwks-create")
			result, err := runOidcJwksCreate(logger, datastoreConfig, kmsConfig, cmdConfig)
			if err != nil {
				log.Errorf("oidc jwks create command failed: %v", err)
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
	cmd.Flags().StringVar(&cmdConfig.currentJwksID, "current-jwks-id", "", "Current JWKS ID")
	cmd.Flags().StringVar(&cmdConfig.nextJwksID, "next-jwks-id", "", "Next JWKS ID to use")
	cmd.Flags().StringVar(&cmdConfig.alg, "alg", "RS256", "The specific asymmetric rfc7518 JWA algorithm to be used to generated the key. One of: [RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512]")

	return cmd
}

func runOidcJwksCreate(logger log.Logger, datastoreConfig *config.DatastoreConfig, kmsConfig *config.KMSConfig, cmdConfig *oidcJwksCreateConfig) (interface{}, error) {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return nil, err
	}
	kmsProvider, err := NewKMSProvider(logger, kmsConfig)
	if err != nil {
		return nil, err
	}
	jwksCreate := NewJwksCreateCmd(logger, dsClient, kmsProvider)
	currentJwksID, err := oidcJwksCreateOrGet(cmdConfig.currentJwksID, cmdConfig.alg, jwksCreate, dsClient)
	if err != nil {
		return nil, err
	}
	nextJwksID, err := oidcJwksCreateOrGet(cmdConfig.nextJwksID, cmdConfig.alg, jwksCreate, dsClient)
	if err != nil {
		return nil, err
	}
	if currentJwksID == nextJwksID {
		return nil, errors.Errorf("current and next OIDC jwksID must be different: %s", currentJwksID)
	}

	id := cmdConfig.oidcJwksID
	if id == "" {
		id = uuid.NewString()
	}
	now := time.Now()
	oidcJWKS := &model.OidcJWKS{
		ID:             id,
		CreatedAt:      now,
		CurrentJwksID:  currentJwksID,
		NextJwksID:     nextJwksID,
		RotationMode:   0,
		RotationPeriod: 0,
		LastRotated:    now,
		Version:        0,
	}
	err = dsClient.API().CreateOidcJWKS(context.Background(), oidcJWKS)
	if err != nil {
		return nil, err
	}
	return oidcJWKS, nil
}

func oidcJwksCreateOrGet(jwksID string, alg string, jwksCreate *jwksCreateCmd, dsClient client.Client) (string, error) {
	if jwksID == "" {
		if err := checkAllowedOidcJwks(alg); err != nil {
			return "", err
		}
		jwksID = uuid.NewString()
		_, err := jwksCreate.Run(&jwksCreateConfig{
			jwksID: jwksID,
			alg:    alg,
			use:    "sig",
		})
		if err != nil {
			return "", err
		}
		return jwksID, nil
	} else {
		key, err := getJwksByID(dsClient, jwksID)
		if err != nil {
			return "", err
		}
		if key.Use != "sig" {
			return "", errors.Errorf("OIDC JWKS requires use 'sig', but got '%s'", key.Use)
		}
		if err := checkAllowedOidcJwks(key.Alg); err != nil {
			return "", err
		}
		return jwksID, nil
	}
}

var (
	allowedOidcAlgos = map[string]struct{}{
		"RS256": {},
		"RS384": {},
		"RS512": {},
		"ES256": {},
		"ES384": {},
		"ES512": {},
		"PS256": {},
		"PS384": {},
		"PS512": {},
	}
)

func checkAllowedOidcJwks(alg string) error {
	if _, ok := allowedOidcAlgos[alg]; !ok {
		return errors.Errorf("OIDC JWKS requires default asymmetric alg, but '%s'", alg)
	}
	return nil
}
