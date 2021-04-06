package cmd

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	jwksCmd.AddCommand(newJwksDeleteCmd())
}

type jwksDeleteConfig struct {
	jwksID string

	use string
	kid string
}

func (c *jwksDeleteConfig) Validate() error {
	if c.kid == "" && c.jwksID == "" {
		return errors.New("either jwks-id or kid is required")
	}
	return nil
}

func newJwksDeleteCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	cmdConfig := new(jwksDeleteConfig)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete JWKS",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("jwks-delete")
			err := runJwksDelete(logger, datastoreConfig, cmdConfig)
			if err != nil {
				log.Errorf("jwks delete command failed: %v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.jwksID, "jwks-id", "", "Identifier of the jwks, JWKSID")
	cmd.Flags().StringVar(&cmdConfig.use, "use", "sig", "How the key is meant to be used. One of: [sig, enc]")
	cmd.Flags().StringVar(&cmdConfig.kid, "kid", "", "Unique key identifier. The Key ID is generated if not specified.")

	return cmd
}

func runJwksDelete(logger log.Logger, datastoreConfig *config.DatastoreConfig, cmdConfig *jwksDeleteConfig) error {
	dsClient, err := getDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return err
	}
	if cmdConfig.jwksID != "" {
		return dsClient.API().DeleteJWKS(context.Background(), cmdConfig.jwksID)
	} else {
		return dsClient.API().DeleteJWKSByKidUse(context.Background(), cmdConfig.kid, cmdConfig.use)
	}
}
