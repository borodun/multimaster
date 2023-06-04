package cmd

import (
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start node",
		Short: "Start specific nodes",
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
}
