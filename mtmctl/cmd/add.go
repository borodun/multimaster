package cmd

import (
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add cluster or node to cluster",
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
