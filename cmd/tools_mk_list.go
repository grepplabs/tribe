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

func newMkListCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()
	outputConfig := config.NewOutputConfig()
	paginationConfig := config.NewPaginationConfig()

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
			result, err := runMkList(logger, dbConfig, paginationConfig)
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
	cmd.Flags().AddFlagSet(paginationConfig.FlagSet())

	return cmd
}

func runMkList(logger log.Logger, dbConfig *config.DBConfig, paginationConfig *config.PaginationConfig) (*dtomodel.KMSKeysetList, error) {
	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create sql client failed")
	}
	return dbClient.API().ListKMSKeysets(context.Background(), utils.Int64(paginationConfig.Offset), utils.Int64(paginationConfig.Limit))
}
