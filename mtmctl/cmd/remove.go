package cmd

import (
	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove cluster or drop node from cluster",
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)
}
