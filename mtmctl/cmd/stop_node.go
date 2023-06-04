package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/pgctl"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stopNodeCmd = &cobra.Command{
		Use:   "node node_name,...",
		Short: "Stop specific nodes",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("you need specify comma-separated nodes")
			}
			nodes = strings.Split(args[0], ",")

			mtmNodePgctl := &pgctl.PGCtl{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, stopNodeConnConf()),
				Node:        "",
				Action:      "stop",
			}

			mtmNodePgctl.Run()
		},
	}
)

func init() {
	stopCmd.AddCommand(stopNodeCmd)
}

func stopNodeConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, v := range nodes {
		conf := connection.Conf{
			ConnName:    v,
			ConnectSsh:  true,
			SshRequired: true,
		}

		connConfs = append(connConfs, conf)
	}

	return connConfs
}
