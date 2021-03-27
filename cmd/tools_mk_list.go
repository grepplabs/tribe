package cmd

import (
	"context"
	"github.com/grepplabs/tribe/pkg/utils"
	"os"

	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	dtomodel "github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	mkCmd.AddCommand(newMkListCmd())
}

type mkListCmdConfig struct {
	Limit  int64
	Offset int64
}

func newMkListCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()
	outputConfig := config.NewOutputConfig()
	mkCmdConfig := new(mkListCmdConfig)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List master key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := outputConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			producer := outputConfig.MustGetProducer()

			logger := log.NewLogger(logConfig.Configuration).WithName("mk")
			result, err := runMkList(logger, dbConfig, mkCmdConfig)
			if err != nil {
				log.Errorf("mk list command failed: %v", err)
				os.Exit(1)
			}
			if result != nil {
				err = producer.Produce(os.Stdout, result)
				if err != nil {
					log.Errorf("failed to write result: %v", err)
					os.Exit(1)
				}
			}
		},
	}
	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().Int64Var(&mkCmdConfig.Limit, "limit", 0, "The numbers of entries to return")
	cmd.Flags().Int64Var(&mkCmdConfig.Offset, "offset", 0, "The number of items to skip before starting to collect the result set")

	return cmd
}

func runMkList(logger log.Logger, dbConfig *config.DBConfig, mkCmdConfig *mkListCmdConfig) (*dtomodel.KMSKeysetList, error) {
	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create sql client failed")
	}
	return dbClient.API().ListKMSKeysets(context.Background(), utils.Int64(mkCmdConfig.Offset), utils.Int64(mkCmdConfig.Limit))
}
