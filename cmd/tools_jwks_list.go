package cmd

import (
	"context"
	"github.com/grepplabs/tribe/pkg/utils"
	"os"

	"github.com/grepplabs/tribe/config"
	dtomodel "github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/spf13/cobra"
)

func init() {
	jwksCmd.AddCommand(newJwksListCmd())
}

func newJwksListCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	outputConfig := config.NewOutputConfig()
	paginationConfig := config.NewPaginationConfig()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List JWKS",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := outputConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			producer := outputConfig.MustGetProducer()

			logger := log.NewLogger(logConfig.Configuration).WithName("jwks-list")
			result, err := runJwksList(logger, datastoreConfig, paginationConfig)
			if err != nil {
				log.Errorf("jwks list command failed: %v", err)
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
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())
	cmd.Flags().AddFlagSet(paginationConfig.FlagSet())

	return cmd
}

func runJwksList(logger log.Logger, datastoreConfig *config.DatastoreConfig, paginationConfig *config.PaginationConfig) (*dtomodel.JWKSList, error) {
	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return nil, err
	}
	return dsClient.API().ListJWKS(context.Background(), utils.Int64(paginationConfig.Offset), utils.Int64(paginationConfig.Limit))
}
