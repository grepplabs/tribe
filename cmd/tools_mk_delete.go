package cmd

import (
	"context"
	"os"

	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	mkCmd.AddCommand(newMkDeleteCmd())
}

type mkDeleteCmdConfig struct {
	keysetID string
}

func newMkDeleteCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()
	cmdConfig := new(mkDeleteCmdConfig)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete master key",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("mk-delete")
			err := runMkDelete(logger, dbConfig, cmdConfig)
			if err != nil {
				log.Errorf("mk delete command failed: %v", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.keysetID, "keyset-id", "", "Identifier of the keyset")
	_ = cmd.MarkFlagRequired("keyset-id")

	return cmd
}

func runMkDelete(logger log.Logger, dbConfig *config.DBConfig, cmdConfig *mkDeleteCmdConfig) error {
	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return errors.Wrap(err, "create sql client failed")
	}
	return dbClient.API().DeleteKMSKeyset(context.Background(), cmdConfig.keysetID)
}
