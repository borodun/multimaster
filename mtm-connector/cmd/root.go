package cmd

import (
	log "github.com/sirupsen/logrus"
	"mtm-connector/internal/serve"
	"os"

	"github.com/spf13/cobra"
)

var (
	port string
	url  string
)

var rootCmd = &cobra.Command{
	Use:   "mtm-connector --conn-url connection-string",
	Short: "API server for managing multimaster cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			log.Fatal("you need to specify connection URL")
		}
		serve.Run(port, url)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "port that server will listen on (default 8080)")
	rootCmd.Flags().StringVarP(&url, "conn-url", "u", "", "connection URL for multimaster node (example: postgresql://mtmuser@node1:5432/mydb)")
}
