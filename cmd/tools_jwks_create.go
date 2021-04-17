package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/jwk"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
	"os"
	"time"
)

func init() {
	jwksCmd.AddCommand(newJwksCreateCmd())
}

type jwksCreateConfig struct {
	jwksID string

	alg string
	use string
}

func (c *jwksCreateConfig) Validate() error {
	return nil
}

func newJwksCreateCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	kmsConfig := config.NewKMSConfig(datastoreConfig)
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(jwksCreateConfig)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create JWKS",
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

			logger := log.NewLogger(logConfig.Configuration).WithName("jwks-create")
			dsClient, err := NewDatastoreClient(logger, datastoreConfig)
			if err != nil {
				log.Errorf("create datastore client failed: %v", err)
				os.Exit(1)
			}
			kmsProvider, err := NewKMSProvider(logger, kmsConfig)
			if err != nil {
				log.Errorf("create kms provider failed: %v", err)
				os.Exit(1)
			}
			result, err := NewJwksCreateCmd(logger, dsClient, kmsProvider).Run(cmdConfig)
			if err != nil {
				log.Errorf("jwks create command failed: %v", err)
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

	cmd.Flags().StringVar(&cmdConfig.jwksID, "jwks-id", "", "Identifier of the jwks used also a kid")
	cmd.Flags().StringVar(&cmdConfig.alg, "alg", "RS256", "The specific rfc7518 JWA algorithm to be used to generated the key. One of: [HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512]")
	cmd.Flags().StringVar(&cmdConfig.use, "use", "sig", "How the key is meant to be used. One of: [sig, enc]")

	return cmd
}

type jwksCreateCmd struct {
	logger      log.Logger
	dsClient    client.Client
	kmsProvider KMSProvider
}

func NewJwksCreateCmd(logger log.Logger, dsClient client.Client, kmsProvider KMSProvider) *jwksCreateCmd {
	return &jwksCreateCmd{
		logger:      logger,
		dsClient:    dsClient,
		kmsProvider: kmsProvider,
	}
}

func (c *jwksCreateCmd) Run(cmdConfig *jwksCreateConfig) (*jose.JSONWebKeySet, error) {
	id := cmdConfig.jwksID
	if id == "" {
		id = uuid.NewString()
	}
	kid := id
	keys, err := jwk.NewJWKSGenerator().Generate(kid, cmdConfig.alg, cmdConfig.use)
	if err != nil {
		return nil, err
	}
	// persist generated keys
	aead, keyURI, err := c.kmsProvider.NewAEAD(id)
	if err != nil {
		return nil, errors.Wrap(err, "Get AEAD failed")
	}
	bytes, err := json.Marshal(keys)
	if err != nil {
		return nil, err
	}
	encryptedKeys, err := aead.Encrypt(bytes, []byte{})
	if err != nil {
		return nil, errors.Wrap(err, "AEAD keys encryption failed")
	}
	jwks := &model.JWKS{
		ID:            id,
		CreatedAt:     time.Now(),
		Kid:           kid,
		Alg:           cmdConfig.alg,
		Use:           cmdConfig.use,
		KMSKeyURI:     keyURI,
		EncryptedJwks: base64.StdEncoding.EncodeToString(encryptedKeys),
	}
	err = c.dsClient.API().CreateJWKS(context.Background(), jwks)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
