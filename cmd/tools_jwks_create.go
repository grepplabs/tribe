package cmd

import (
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/jwk"
	"github.com/grepplabs/tribe/pkg/log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	jwksCmd.AddCommand(newJwksCreateCmd())
}

type jwksCreateConfig struct {
	alg string
	use string
	kid string
}

func (c *jwksCreateConfig) Validate() error {
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
				log.Errorf("mk create command failed: %v", err)
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
	return keys, nil
}
