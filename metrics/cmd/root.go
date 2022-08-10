package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metrics/internal/config"
	"metrics/internal/metrics"
)

var (
	configPath string
	cfg        config.Config

	rootCmd = &cobra.Command{
		Use:   "mmts-metrics [--config config.yaml]",
		Short: "Metrics server for Postgres Multimaster ",
		Run: func(cmd *cobra.Command, args []string) {
			metrics.Run(cfg)
		},
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is ./config.yaml)")
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetDefault("spec.listenPort", 8080)
	viper.SetDefault("spec.interval", 10)
	viper.SetDefault("spec.queryTimeout", 5)
	viper.SetDefault("spec.connectionPoolMaxSize", 5)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	log.Infof("using config file: %s", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	if err := validate(); err != nil {
		log.Fatal(err)
	}
}

func validate() error {
	names := make(map[string]bool)
	for _, conf := range cfg.Spec.Databases {
		if conf.Name == "" {
			return fmt.Errorf("failed to validate configuration. Database name cannot be empty")
		}
		if conf.URL == "" {
			return fmt.Errorf("failed to validate configuration. URL cannot be empty in the '%s' database", conf.Name)
		}
		if names[conf.Name] {
			return fmt.Errorf("failed to validate configuration. A database named '%s' has already been declared", conf.Name)
		}
		names[conf.Name] = true
	}

	return nil
}
