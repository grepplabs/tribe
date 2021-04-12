package cmd

import (
	"context"
	"os"

	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"
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
	datastoreConfig := config.NewDatastoreConfig()
	cmdConfig := new(mkDeleteCmdConfig)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete master key",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("mk-delete")
			err := runMkDelete(logger, datastoreConfig, cmdConfig)
			if err != nil {
				log.Errorf("mk delete command failed: %v", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())
	cmd.Flags().StringVar(&cmdConfig.keysetID, "keyset-id", "", "Identifier of the keyset")
	_ = cmd.MarkFlagRequired("keyset-id")

	return cmd
}

func runMkDelete(logger log.Logger, datastoreConfig *config.DatastoreConfig, cmdConfig *mkDeleteCmdConfig) error {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return err
	}
	return dsClient.API().DeleteKMSKeyset(context.Background(), cmdConfig.keysetID)
}
