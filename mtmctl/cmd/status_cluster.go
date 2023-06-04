package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/status"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	statusClusterCmd = &cobra.Command{
		Use:   "cluster node_name",
		Short: "Get status nodes",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("you need specify 1 node")
			}

			mtmClusterStatus := &status.MtmStatus{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, statusClusterConnConf()),
				Node:        args[0],
			}

			mtmClusterStatus.Run()
		},
	}
)

func init() {
	statusCmd.AddCommand(statusClusterCmd)
}

func statusClusterConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, node := range cfg.Toolbox.Connections {
		nodeConf := connection.Conf{
			ConnName:   node.Name,
			ConnectDb:  true,
			ConnectSsh: true,
		}

		connConfs = append(connConfs, nodeConf)
	}

	return connConfs
}
