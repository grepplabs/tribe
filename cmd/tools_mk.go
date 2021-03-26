package cmd

import (
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

func init() {
	toolsCmd.AddCommand(newMkCmd())
}

func newMkCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()

	cmd := &cobra.Command{
		Use:   "mk",
		Short: "Master keys tools",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("mk")
			logger.Infof("mk called")

			err := runMk(logger, dbConfig)
			if err != nil {
				log.Errorf("mk command failed: %v", err)
			}
		},
	}
	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())

	return cmd
}

type MkCommand struct {
	dbConfig  *config.DBConfig
	logConfig *config.LogConfig
}

func runMk(logger log.Logger, dbConfig *config.DBConfig) error {

	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return errors.Wrap(err, "create sql client failed")
	}
	_ = dbClient
	return nil
}
