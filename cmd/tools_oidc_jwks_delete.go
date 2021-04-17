package cmd

import (
	"context"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	oidcJwksCmd.AddCommand(newOidcJwksDeleteCmd())
}

type oidcJwksDeleteConfig struct {
	oidcJwksID string
}

func (c *oidcJwksDeleteConfig) Validate() error {
	return nil
}

func newOidcJwksDeleteCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	cmdConfig := new(oidcJwksDeleteConfig)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete OIDC JWKS",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("oidc-jwks-delete")
			err := runOidcJwksDelete(logger, datastoreConfig, cmdConfig)
			if err != nil {
				log.Errorf("oidc jwks delete command failed: %v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.oidcJwksID, "oidc-jwks-id", "", "Identifier of the oidc jwks")
	_ = cmd.MarkFlagRequired("oidc-jwks-id")

	return cmd
}

func runOidcJwksDelete(logger log.Logger, datastoreConfig *config.DatastoreConfig, cmdConfig *oidcJwksDeleteConfig) error {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return err
	}
	return dsClient.API().DeleteOidcJWKS(context.Background(), cmdConfig.oidcJwksID)
}
