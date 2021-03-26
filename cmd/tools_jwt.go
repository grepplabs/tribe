package cmd

import (
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"

	"github.com/spf13/cobra"
)

func init() {
	toolsCmd.AddCommand(newJwtCmd())
}

func newJwtCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()

	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "JWT tools",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("jwt")
			logger.Infof("jwt called")
		},
	}

	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())

	return cmd
}
