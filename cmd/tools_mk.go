package cmd

import (
	"github.com/spf13/cobra"
)

var mkCmd = &cobra.Command{
	Use:   "mk",
	Short: "Master keys tools",
}

func init() {
	toolsCmd.AddCommand(mkCmd)
}
