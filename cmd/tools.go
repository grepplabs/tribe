package cmd

import (
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Tribe tools",
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}

func getDatastoreClient(logger log.Logger, datastoreConfig *config.DatastoreConfig) (client.Client, error) {
	switch strings.ToLower(datastoreConfig.Provider) {
	case "db":
		dbClient, err := client.NewSQLClient(logger, &datastoreConfig.DBConfig)
		if err != nil {
			return nil, errors.Wrap(err, "create sql client failed")
		}
		return dbClient, nil
	case "minio":
		minioClient, err := client.NewMinioClient(logger, &datastoreConfig.MinioConfig)
		if err != nil {
			return nil, errors.Wrap(err, "create minio client failed")
		}
		return minioClient, nil
	default:
		return nil, errors.Errorf("Unsupported datastore provider: %v", datastoreConfig.Provider)
	}
}
