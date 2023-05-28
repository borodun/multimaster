package cmd

import (
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status [connection_names]",
		Short: "Get whole cluster or specific node status",
	}
)

func init() {
	rootCmd.AddCommand(statusCmd)
}
