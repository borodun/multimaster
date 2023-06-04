package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/query"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	queryClusterCmd = &cobra.Command{
		Use:   "cluster node_name",
		Short: "Execute SQL query on all nodes of a cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("you need specify 1 node")
			}

			if queryStr == "" {
				log.Fatal("you need specify query")
			}

			mtmClusterQuery := &query.MtmQuery{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, queryClusterConnConf()),
				Node:        args[0],
				Query:       queryStr,
			}

			mtmClusterQuery.Run()
		},
	}
)

func init() {
	queryCmd.AddCommand(queryClusterCmd)
}

func queryClusterConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, node := range cfg.Toolbox.Connections {
		nodeConf := connection.Conf{
			ConnName:    node.Name,
			ConnectDb:   true,
			ConnectSsh:  true,
			SshRequired: true,
		}

		connConfs = append(connConfs, nodeConf)
	}

	return connConfs
}
