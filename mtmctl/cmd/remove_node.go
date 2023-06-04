package cmd

import (
	"mtmctl/internal/connection"
	"mtmctl/internal/remove"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	initNode string

	removeNodeCmd = &cobra.Command{
		Use:   "node connection_name -i connection_name",
		Short: "Remove node from multimaster cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("you need to specify node to remove in args")
			}
			removeNode := args[0]

			if initNode == "" {
				log.Fatal("you need to specify init node connection")
			}

			mtmRemoveNode := remove.MtmRemoveNode{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, removeConnConf()),
				InitNode:    initNode,
				RemoveNode:  removeNode,
			}

			mtmRemoveNode.Run()
		},
	}
)

func init() {
	removeCmd.AddCommand(removeNodeCmd)

	removeNodeCmd.Flags().StringVarP(&initNode, "init-node", "i", "", "specify node that will initiate removal")
}

func removeConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	for _, node := range cfg.Toolbox.Connections {
		nodeConf := connection.Conf{
			ConnName:  node.Name,
			ConnectDb: true,
		}

		connConfs = append(connConfs, nodeConf)
	}

	return connConfs
}
