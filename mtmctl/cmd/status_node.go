package cmd

import (
	"fmt"
	"mtmctl/internal/connection"
	"mtmctl/internal/status"
	"strings"

	"github.com/spf13/cobra"
)

var (
	statusNodeCmd = &cobra.Command{
		Use:   "node [connection_names,...]",
		Short: "Get status nodes",
		Run: func(cmd *cobra.Command, args []string) {
			mtmNodeStatus := &status.MtmStatus{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, getConnConf(args)),
				Node:        "",
			}

			mtmNodeStatus.Run()
		},
	}
)

func init() {
	statusCmd.AddCommand(statusNodeCmd)
}

func getConnConf(args []string) []connection.Conf {
	if len(args) == 1 {
		return fromArgs(strings.Split(args[0], ","))
	} else if len(args) > 1 {
		fmt.Println("Warning: using first argument only")
		return fromArgs(strings.Split(args[0], ","))
	} else {
		return fromConfig()
	}
}

func fromArgs(args []string) []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, arg := range args {
		conf := connection.Conf{
			ConnName:   arg,
			ConnectDb:  true,
			ConnectSsh: true,
		}
		connConfs = append(connConfs, conf)
	}

	return connConfs
}

func fromConfig() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, conn := range cfg.Toolbox.Connections {
		conf := connection.Conf{
			ConnName:   conn.Name,
			ConnectDb:  true,
			ConnectSsh: true,
		}
		connConfs = append(connConfs, conf)
	}

	return connConfs
}
