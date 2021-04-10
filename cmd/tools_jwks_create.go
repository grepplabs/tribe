package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/jwk"
	"github.com/grepplabs/tribe/pkg/kms/dbkms"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	kid string

	save         bool
	kmsKeysetURI string
	masterSecret string
}

func (c *jwksCreateConfig) Validate() error {
	if c.save {
		if c.kmsKeysetURI == "" {
			return errors.New("kms-keyset-uri is required when save is enabled")
		}
	}
	return nil
}

func newJwksCreateCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
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
			result, err := runJwksCreate(logger, datastoreConfig, cmdConfig)
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
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.jwksID, "jwks-id", "", "Identifier of the jwks")
	cmd.Flags().StringVar(&cmdConfig.alg, "alg", "RS256", "The specific rfc7518 JWA algorithm to be used to generated the key. One of: [HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512]")
	cmd.Flags().StringVar(&cmdConfig.use, "use", "sig", "How the key is meant to be used. One of: [sig, enc]")
	cmd.Flags().StringVar(&cmdConfig.kid, "kid", "", "Unique key identifier. The Key ID is generated if not specified.")

	cmd.Flags().BoolVar(&cmdConfig.save, "save", true, "Save the JWKS in persistent store")
	cmd.Flags().StringVar(&cmdConfig.kmsKeysetURI, "kms-keyset-uri", "", "URI of the master KMS keyset to use to encrypt the JWKS")
	cmd.Flags().StringVar(&cmdConfig.masterSecret, "master-secret", "", "KMS master secret")

	return cmd
}

func runJwksCreate(logger log.Logger, datastoreConfig *config.DatastoreConfig, cmdConfig *jwksCreateConfig) (interface{}, error) {
	id := cmdConfig.jwksID
	if id == "" {
		id = uuid.NewString()
	}
	kid := cmdConfig.kid
	if kid == "" {
		kid = uuid.NewString()
	}
	keys, err := jwk.NewJWKSGenerator().Generate(kid, cmdConfig.alg, cmdConfig.use)
	if err != nil {
		return nil, err
	}
	if cmdConfig.save {
		dsClient, err := getDatastoreClient(logger, datastoreConfig)
		if err != nil {
			return nil, err
		}
		dbkmsClient, err := dbkms.NewClient(dbkms.WithMasterSecret(cmdConfig.masterSecret), dbkms.WithLogger(logger), dbkms.WithDBClient(dsClient))
		if err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(keys)
		if err != nil {
			return nil, err
		}
		aead, err := dbkmsClient.GetAEAD(cmdConfig.kmsKeysetURI)
		if err != nil {
			return nil, errors.Wrap(err, "Get AEAD failed")
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
			KMSKeysetURI:  cmdConfig.kmsKeysetURI,
			EncryptedJwks: base64.StdEncoding.EncodeToString(encryptedKeys),
		}
		err = dsClient.API().CreateJWKS(context.Background(), jwks)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
