package cmd

import (
	"github.com/spf13/cobra"
)

var oidcJwksCmd = &cobra.Command{
	Use:   "jwks",
	Short: "OIDC JWKS tools",
}

func init() {
	oidcCmd.AddCommand(oidcJwksCmd)
}
