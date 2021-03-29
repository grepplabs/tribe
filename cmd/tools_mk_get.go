package cmd

import (
	"context"
	"os"

	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	dtomodel "github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	mkCmd.AddCommand(newMkGetCmd())
}

type mkGetCmdConfig struct {
	keysetID string
}

func newMkGetCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	dbConfig := config.NewDBConfig()
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(mkGetCmdConfig)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get master key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := outputConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			producer := outputConfig.MustGetProducer()

			logger := log.NewLogger(logConfig.Configuration).WithName("mk-get")
			result, err := runMkGet(logger, dbConfig, cmdConfig)
			if err != nil {
				log.Errorf("mk get command failed: %v", err)
				os.Exit(1)
			}
			if result == nil {
				log.Errorf("mk get command failed, not found keysetID %s", cmdConfig.keysetID)
				os.Exit(1)
			}
			err = producer.Produce(os.Stdout, result)
			if err != nil {
				log.Errorf("failed to write result: %v", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().AddFlagSet(logConfig.FlagSet())
	cmd.Flags().AddFlagSet(dbConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.keysetID, "keyset-id", "", "Identifier of the keyset, KeysetID")
	_ = cmd.MarkFlagRequired("keyset-id")

	return cmd
}

func runMkGet(logger log.Logger, dbConfig *config.DBConfig, cmdConfig *mkGetCmdConfig) (*dtomodel.KMSKeyset, error) {
	dbClient, err := client.NewSQLClient(logger, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create sql client failed")
	}
	return dbClient.API().GetKMSKeyset(context.Background(), cmdConfig.keysetID)
}
