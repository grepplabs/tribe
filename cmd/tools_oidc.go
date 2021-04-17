package cmd

import (
	"github.com/spf13/cobra"
)

var oidcCmd = &cobra.Command{
	Use:   "oidc",
	Short: "OIDC tools",
}

func init() {
	toolsCmd.AddCommand(oidcCmd)
}
