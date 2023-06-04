package cmd

import (
	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop {node | cluster}",
		Short: "Stop specific nodes or whole cluster",
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
}
