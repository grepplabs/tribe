package cmd

import (
	"github.com/spf13/cobra"
)

var jwksCmd = &cobra.Command{
	Use:   "jwks",
	Short: "JWKS tools",
}

func init() {
	toolsCmd.AddCommand(jwksCmd)
}
