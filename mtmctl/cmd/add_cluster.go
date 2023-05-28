package cmd

import (
	"mtmctl/internal/add"
	"mtmctl/internal/connection"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	nodes []string

	addClusterCmd = &cobra.Command{
		Use:   "cluster node1,node2,node3",
		Short: "Create cluster from nodes",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("you need to specify nodes for creating cluster")
			}

			nodes = strings.Split(args[0], ",")

			if len(nodes) < 3 {
				log.Fatal("you need specify at least 3 nodes for creating cluster")
			}

			mtmAddCluster := &add.MtmAddCluster{
				Cfg:            cfg,
				SSHConnections: connection.Connect(cfg, addClusterConnConf()),
				Nodes:          nodes,
			}
			mtmAddCluster.Run()
		},
	}
)

func init() {
	addCmd.AddCommand(addClusterCmd)
}

func addClusterConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, v := range nodes {
		conf := connection.Conf{
			ConnName:    v,
			ConnectDb:   false,
			ConnectSsh:  true,
			DbRequired:  false,
			SshRequired: true,
		}

		connConfs = append(connConfs, conf)
	}

	return connConfs
}
