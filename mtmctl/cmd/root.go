package cmd

import (
	"fmt"
	"mtmctl/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	verbose    bool
	cfg        config.Config

	rootCmd = &cobra.Command{
		Use:   "mtmctl <cmd> [--config config.yaml]",
		Short: "Toolbox for Postgres multimaster",
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(start)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "turn on verbose output")
}

func start() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}

	initConfig()
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

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
	if cfg.Toolbox.PGDATA == "" {
		return fmt.Errorf("failed to validate configuration. Path to PGDATA cannot be empty")
	}
	if cfg.Toolbox.PGBIN == "" {
		return fmt.Errorf("failed to validate configuration. Path to PGBIN cannot be empty")
	}

	names := make(map[string]bool)
	for _, conn := range cfg.Toolbox.Connections {
		if conn.Name == "" {
			return fmt.Errorf("failed to validate configuration. Connection name cannot be empty")
		}
		if conn.Ssh.Host == "" {
			return fmt.Errorf("failed to validate configuration. SSH Host cannot be empty in the '%s' connection", conn.Name)
		}
		if names[conn.Name] {
			return fmt.Errorf("failed to validate configuration. A connection named '%s' has already been used", conn.Name)
		}
		names[conn.Name] = true
	}

	return nil
}
