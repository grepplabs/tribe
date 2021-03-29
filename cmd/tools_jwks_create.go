package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/tink/go/keyset"
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/jwk"
	"github.com/grepplabs/tribe/pkg/kms/dbkms"
	"github.com/grepplabs/tribe/pkg/kms/masterkey"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/grepplabs/tribe/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	jwksCmd.AddCommand(newJwksCreateCmd())
}

type jwksCreateConfig struct {
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
	dbConfig := config.NewDBConfig()
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
			result, err := runJwksCreate(logger, dbConfig, cmdConfig)
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
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.alg, "alg", "RS256", "The specific rfc7518 JWA algorithm to be used to generated the key. One of: [HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512]")
	cmd.Flags().StringVar(&cmdConfig.use, "use", "sig", "How the key is meant to be used. One of: [sig, enc]")
	cmd.Flags().StringVar(&cmdConfig.kid, "kid", "", "Unique key identifier. The Key ID is generated if not specified.")

	cmd.Flags().BoolVar(&cmdConfig.save, "save", true, "Save the JWKS in persistent store")
	cmd.Flags().StringVar(&cmdConfig.kmsKeysetURI, "kms-keyset-uri", "", "URI of the master KMS keyset to use to encrypt the JWKS")
	cmd.Flags().StringVar(&cmdConfig.masterSecret, "master-secret", "", "KMS master secret")

	return cmd
}

func runJwksCreate(logger log.Logger, dbConfig *config.DBConfig, cmdConfig *jwksCreateConfig) (interface{}, error) {
	kid := cmdConfig.kid
	if kid == "" {
		kid = uuid.NewString()
	}
	gen := jwk.NewJWKSGenerator()
	keys, err := gen.Generate(kid, cmdConfig.alg, cmdConfig.use)
	if err != nil {
		return nil, err
	}
	if cmdConfig.save {
		var mk masterkey.MasterKeyset
		mk, err = getMasterKey(logger, dbConfig, cmdConfig.kmsKeysetURI, cmdConfig.masterSecret)
		if err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(keys)
		if err != nil {
			return nil, err
		}
		a := dbkms.NewAEAD(func() (*keyset.Handle, error) {
			return mk.GetKeyset(), nil
		})
		encryptedKeys, err := a.Encrypt(bytes, []byte{})
		if err != nil {
			return nil, errors.Wrap(err, "AEAD keys encryption failed")
		}
		encodedEncryptedKeys := base64.StdEncoding.EncodeToString(encryptedKeys)
		dbClient, err := client.NewSQLClient(logger, dbConfig)
		if err != nil {
			return nil, errors.Wrap(err, "create sql client failed")
		}
		jwks := &model.JWKS{
			ID:            uuid.NewString(),
			Kid:           kid,
			Alg:           cmdConfig.alg,
			Use:           cmdConfig.use,
			KMSKeysetURI:  cmdConfig.kmsKeysetURI,
			EncryptedJwks: encodedEncryptedKeys,
			Description:   utils.EmptyToNullString(fmt.Sprintf("Used master keyset KeyId %d", mk.GetKeyset().KeysetInfo().PrimaryKeyId)),
		}
		err = dbClient.API().CreateJWKS(context.Background(), jwks)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func getMasterKey(logger log.Logger, dbConfig *config.DBConfig, keyURI string, masterSecret string) (masterkey.MasterKeyset, error) {
	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create sql client failed")
	}
	const dbPrefix = "db://"
	if !strings.HasPrefix(strings.ToLower(keyURI), dbPrefix) {
		return nil, fmt.Errorf("uriPrefix must start with %s, but got %s", dbPrefix, keyURI)
	}
	keysetID := strings.TrimPrefix(keyURI, dbPrefix)
	ks, err := dbClient.API().GetKMSKeyset(context.Background(), keysetID)
	if err != nil {
		return nil, errors.Wrapf(err, "get kms keyset failed: %s", keysetID)
	}
	if ks == nil {
		return nil, errors.Errorf("kms keyset not found: %s", keysetID)
	}
	if masterSecret == "" {
		return nil, errors.Errorf("master-secret is required for kms-keyset-uri: %s", keyURI)
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
