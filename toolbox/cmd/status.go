package cmd

import (
	"backup/internal/connection"
	"backup/internal/status"
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status [connection_names]",
		Short: "Get status nodes",
		Run: func(cmd *cobra.Command, args []string) {
			mtmStatus := &status.MtmStatus{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, getConnConf(args)),
			}

			mtmStatus.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getConnConf(args []string) []connection.Conf {
	if len(args) != 0 {
		return fromArgs(args)
	} else {
		return fromConfig()
	}
}

func fromArgs(args []string) []connection.Conf {
	connConfs := make([]connection.Conf, len(args))

	for _, arg := range args {
		conf := connection.Conf{
			ConnName:    arg,
			ConnectDb:   true,
			ConnectSsh:  false,
			DbRequired:  true,
			SshRequired: false,
		}
		connConfs = append(connConfs, conf)
	}

	return connConfs
}

func fromConfig() []connection.Conf {
	connConfs := make([]connection.Conf, len(cfg.Toolbox.Connections))

	for _, conn := range cfg.Toolbox.Connections {
		conf := connection.Conf{
			ConnName:    conn.Name,
			ConnectDb:   true,
			ConnectSsh:  false,
			DbRequired:  false,
			SshRequired: false,
		}
		connConfs = append(connConfs, conf)
	}

	return connConfs
}
