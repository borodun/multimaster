package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/remove"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	removeNodes []string

	removeClusterCmd = &cobra.Command{
		Use:   "cluster node1,node2,node3",
		Short: "Remove multimaster cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("you need to specify all the nodes that make up the cluster")
			}

			removeNodes = strings.Split(args[0], ",")

			mtmRemoveCluster := &remove.MtmRemoveCluster{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, rmoveClusterConnConf()),
				Nodes:       removeNodes,
			}
			mtmRemoveCluster.Run()
		},
	}
)

func init() {
	removeCmd.AddCommand(removeClusterCmd)
}

func rmoveClusterConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, v := range removeNodes {
		conf := connection.Conf{
			ConnName:    v,
			ConnectDb:   true,
			ConnectSsh:  true,
			DbRequired:  true,
			SshRequired: true,
		}

		connConfs = append(connConfs, conf)
	}

	return connConfs
}
