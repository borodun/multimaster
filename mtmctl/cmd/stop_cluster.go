package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/pgctl"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stopClusterCmd = &cobra.Command{
		Use:   "cluster node_name",
		Short: "Stop all nodes in cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("you need specify 1 node")
			}

			mtmClusterPgctl := &pgctl.PGCtl{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, stopClusterConnConf()),
				Node:        args[0],
				Action:      "stop",
			}

			mtmClusterPgctl.Run()
		},
	}
)

func init() {
	stopCmd.AddCommand(stopClusterCmd)
}

func stopClusterConnConf() []connection.Conf {
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
