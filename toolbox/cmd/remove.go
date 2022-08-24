package cmd

import (
	"backup/internal/connection"
	"backup/internal/remove"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	initNode string

	removeCmd = &cobra.Command{
		Use:   "remove connection_name -i connection_name",
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
				Connections: connection.Connect(cfg, removeConnConf(removeNode)),
				InitNode:    initNode,
				RemoveNode:  removeNode,
			}

			mtmRemoveNode.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(&initNode, "init-node", "i", "", "specify node that will initiate removal")
}

func removeConnConf(node string) []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	nodeConf := connection.Conf{
		ConnName:   node,
		ConnectDb:  true,
		DbRequired: true,
	}
	initConf := connection.Conf{
		ConnName:   initNode,
		ConnectDb:  true,
		DbRequired: true,
	}
	connConfs = append(connConfs, nodeConf, initConf)
	return connConfs
}
