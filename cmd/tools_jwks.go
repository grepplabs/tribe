package cmd

import (
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"

	"github.com/spf13/cobra"
)

func init() {
	toolsCmd.AddCommand(newJwksCmd())
}

func newJwksCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()

	cmd := &cobra.Command{
		Use:   "jwks",
		Short: "JWKS tools",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewLogger(logConfig.Configuration).WithName("jwks")
			logger.Infof("jwks called")
		},
	}

	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())

	return cmd
}
