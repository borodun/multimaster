package utils

import (
	"fmt"
	"mtmctl/internal/connection"

	log "github.com/sirupsen/logrus"
)

func CheckDatabase(db *connection.DB) error {
	if err := db.TruePing(); err != nil {
		return fmt.Errorf("cannot ping '%s' database: %s", db.GetName(), err.Error())
	}
	if stat := db.MtmStatus(); stat != "online" {
		return fmt.Errorf("'%s' node is not online, current status: %s", db.GetName(), stat)
	}
	return nil
}

func CheckDatabases(dbs ...*connection.DB) {
	for _, db := range dbs {
		if err := CheckDatabase(db); err != nil {
			println(err.Error())
			log.WithError(err).Fatal("bad connection to db")
		}
	}
}

func GetClusterNodes(conns connection.Connections, name string) connection.Connections {
	nodeDB := conns[name].DB
	nodeId := nodeDB.GetMtmNodeID()
	nodeConnInfo := ""
	clusterNodes := nodeDB.GetMtmNodes()

	for _, nodeTup := range clusterNodes {
		if nodeTup.Id == nodeId {
			nodeConnInfo = nodeTup.Conninfo
		}
	}

	nodes := make(connection.Connections)
	nodes[name] = conns[name]

	for nodeName, node := range conns {
		db := node.DB
		if db == nil {
			continue
		}

		id := db.GetMtmNodeID()
		if id == nodeId {
			continue
		}

		nodesTups := db.GetMtmNodes()
		for _, nodeTup := range nodesTups {
			if nodeTup.Id == nodeId && nodeTup.Conninfo == nodeConnInfo {
				fmt.Printf("Discovered node '%s' in the same cluster as node '%s'\n", nodeName, name)
				nodes[nodeName] = node
			}
		}
	}

	return nodes
}
