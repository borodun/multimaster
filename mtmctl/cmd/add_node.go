package cmd

import (
	"mtmctl/internal/add"
	"mtmctl/internal/connection"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	srcNode string
	dstNode string
	connStr string

	addNodeCmd = &cobra.Command{
		Use:   "node -s source-node -d destination-node",
		Short: "Add node to multimaster cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if srcNode == "" {
				log.Fatal("you need to specify source node connection")
			}
			if dstNode == "" {
				log.Fatal("you need to specify destination node connection")
			}

			mtmAddNode := &add.MtmAddNode{
				Cfg:         cfg,
				Connections: connection.Connect(cfg, addConnConf()),
				SrcNode:     srcNode,
				DstNode:     dstNode,
				ConnStr:     connStr,
			}
			mtmAddNode.Run()
		},
	}
)

func init() {
	addCmd.AddCommand(addNodeCmd)

	addNodeCmd.Flags().StringVarP(&srcNode, "source-node", "s", "", "specify source node")
	addNodeCmd.Flags().StringVarP(&dstNode, "destination-node", "d", "", "specify destination node")
	addNodeCmd.Flags().StringVarP(&connStr, "conn-str", "n", "", "connection string of destination node for mtm")
}

func addConnConf() []connection.Conf {
	connConfs := make([]connection.Conf, 0)

	srcConf := connection.Conf{
		ConnName:    srcNode,
		ConnectDb:   true,
		ConnectSsh:  false,
		DbRequired:  true,
		SshRequired: false,
	}
	dstConf := connection.Conf{
		ConnName:    dstNode,
		ConnectDb:   false,
		ConnectSsh:  true,
		DbRequired:  false,
		SshRequired: true,
	}

	connConfs = append(connConfs, srcConf, dstConf)
	return connConfs
}
