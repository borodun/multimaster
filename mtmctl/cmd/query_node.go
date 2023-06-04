package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/query"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	queryNodeCmd = &cobra.Command{
		Use:   "node node_name,...",
		Short: "Execute SQL query on a node",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("you need specify comma-separated nodes")
			}
			nodes = strings.Split(args[0], ",")

			if queryStr == "" {
				log.Fatal("you need specify query")
			}

			mtmNodeQuery := &query.MtmQuery{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, queryNodeConnConf()),
				Node:        "",
				Query:       queryStr,
			}

			mtmNodeQuery.Run()
		},
	}
)

func init() {
	queryCmd.AddCommand(queryNodeCmd)
}

func queryNodeConnConf() []connection.Conf {
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
