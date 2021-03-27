package cmd

import (
	"github.com/spf13/cobra"
)

const (
	defaultKeysetName = "master"
)

var mkCmd = &cobra.Command{
	Use:   "mk",
	Short: "Master keys tools",
}

func init() {
	toolsCmd.AddCommand(mkCmd)
}
