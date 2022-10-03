package cmd

import (
	log "github.com/sirupsen/logrus"
	"mtm-joiner/internal/joiner"
	"os"

	"github.com/spf13/cobra"
)

var (
	url string
)

var rootCmd = &cobra.Command{
	Use:   "mtm-joiner --api-url api-url",
	Short: "API server for managing multimaster cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			log.Fatal("you need to specify URL for API")
		}

		j := &joiner.Joiner{
			URL:    url,
			PGDATA: "~/db",
		}
		j.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "api-url", "u", "", "URL of API server (example: 127.0.0.1:8080)")
}
