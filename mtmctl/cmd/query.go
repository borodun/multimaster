package cmd

import (
	"github.com/spf13/cobra"
)

var (
	queryStr string

	queryCmd = &cobra.Command{
		Use:   "query {node | cluster}",
		Short: "Execute a query on a node or whole cluster",
	}
)

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.PersistentFlags().StringVarP(&queryStr, "query", "q", "", "query that will be executed")
}
