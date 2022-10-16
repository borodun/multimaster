package cmd

import (
	log "github.com/sirupsen/logrus"
	"mtm-joiner/internal/joiner"
	"os"

	"github.com/spf13/cobra"
)

var (
	url       string
	drop      bool
	localAddr string
	port      string
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "mtm-joiner --api-url api-url --local-addr local-addr --port port",
	Short: "API server for managing multimaster cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			log.Fatal("you need to specify URL for API")
		}

		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.FatalLevel)
		}

		j := &joiner.Joiner{
			URL:    url,
			PGDATA: "./db",
			Port:   port,
			Addr:   localAddr,
		}
		j.Start(drop)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "api-url", "u", "", "URL of API server (example: http://127.0.0.1:8080)")
	rootCmd.Flags().StringVarP(&localAddr, "local-addr", "a", "", "Local address of the device (will try to detect automatically)")
	rootCmd.Flags().StringVarP(&port, "port", "p", "15432", "Port for the database (default: 15432)")
	rootCmd.Flags().BoolVarP(&drop, "drop", "d", false, "If node needs to be dropped from cluster (default: false)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Add logs to output (default: false)")
}
