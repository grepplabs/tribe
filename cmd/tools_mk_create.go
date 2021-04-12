package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/grepplabs/tribe/config"
	dtomodel "github.com/grepplabs/tribe/database/model"
	"github.com/grepplabs/tribe/pkg/kms/masterkey"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	mkCmd.AddCommand(newMkCreateCmd())
}

type mkCreateCmdConfig struct {
	keysetID     string
	masterSecret string
}

func (c *mkCreateCmdConfig) Validate() error {
	return nil
}

func newMkCreateCmd() *cobra.Command {
	logConfig := config.NewLogConfig()
	datastoreConfig := config.NewDatastoreConfig()
	outputConfig := config.NewOutputConfig()
	cmdConfig := new(mkCreateCmdConfig)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create master key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdConfig.Validate(); err != nil {
				return err
			}
			if err := outputConfig.Validate(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			producer := outputConfig.MustGetProducer()

			logger := log.NewLogger(logConfig.Configuration).WithName("mk-create")
			result, err := runMkCreate(logger, datastoreConfig, cmdConfig)
			if err != nil {
				log.Errorf("mk create command failed: %v", err)
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
	cmd.Flags().AddFlagSet(datastoreConfig.FlagSet())
	cmd.Flags().AddFlagSet(outputConfig.FlagSet())

	cmd.Flags().StringVar(&cmdConfig.keysetID, "keyset-id", "", "Identifier of the keyset")
	cmd.Flags().StringVar(&cmdConfig.masterSecret, "master-secret", "", "Master secret")
	// flag will be optional when
	_ = cmd.MarkFlagRequired("master-secret")

	return cmd
}

func runMkCreate(logger log.Logger, datastoreConfig *config.DatastoreConfig, cmdConfig *mkCreateCmdConfig) (*dtomodel.KMSKeyset, error) {
	id := cmdConfig.keysetID
	if id == "" {
		id = uuid.NewString()
	}

	mk, err := masterkey.NewMasterKeyset([]byte(cmdConfig.masterSecret))
	if err != nil {
		return nil, errors.Wrap(err, "create master keyset failed")
	}
	encryptedKeyset, err := mk.EncryptKeyset()
	if err != nil {
		return nil, errors.Wrap(err, "encrypt master keyset failed")
	}

	dsClient, err := NewDatastoreClient(logger, datastoreConfig)
	if err != nil {
		return nil, err
	}
	keyset := dtomodel.KMSKeyset{
		ID:              id,
		CreatedAt:       time.Now(),
		EncryptedKeyset: base64.StdEncoding.EncodeToString(encryptedKeyset),
		Description:     fmt.Sprintf("Master keyset KeyId %d", mk.GetKeyset().KeysetInfo().PrimaryKeyId),
	}
	return &keyset, dsClient.API().CreateKMSKeyset(context.Background(), &keyset)
}
