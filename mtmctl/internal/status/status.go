package status

import (
	"fmt"
	"mtmctl/internal/config"
	"mtmctl/internal/connection"
	"mtmctl/internal/utils"
)

type MtmStatus struct {
	Cfg         config.Config
	Connections connection.Connections
	Node        string
}

func (m *MtmStatus) Run() {
	if m.Node != "" {
		m.Connections = m.GetClusterNodes()
	}

	statuses := make(map[string]string)
	for _, name := range m.Connections.GetConnNames() {
		db := m.Connections[name].DB
		if db == nil || !db.Ping() {
			statuses[name] = "offline"
		}

		if !db.HasSharedPreloadLibrary("multimaster") {
			statuses[name] = "multimaster is not in shared_preload_libraries"
			if !db.HasExtension("multimaster") {
				statuses[name] = "multimaster extension is not installed"
			}
		}

		if stat := db.MtmStatus(); stat != "" {
			statuses[name] = stat
		}
	}

	fmt.Println("Multimaster status:")
	utils.PrintSortedMap(statuses)

	fmt.Println("Uptime:")
	m.PrintNodesLA()
}

func (m *MtmStatus) GetClusterNodes() connection.Connections {
	nodeDB := m.Connections[m.Node].DB
	nodeId := nodeDB.GetMtmNodeID()
	nodeConnInfo := ""
	clusterNodes := nodeDB.GetMtmNodes()

	for _, nodeTup := range clusterNodes {
		if nodeTup.Id == nodeId {
			nodeConnInfo = nodeTup.Conninfo
		}
	}

	nodes := make(connection.Connections)

	for nodeName, node := range m.Connections {
		db := node.DB
		if db == nil {
			continue
		}

		nodesTups := db.GetMtmNodes()
		for _, nodeTup := range nodesTups {
			if nodeTup.Id == nodeId && nodeTup.Conninfo == nodeConnInfo {
				nodes[nodeName] = node
			}
		}
	}

	return nodes
}

func (m *MtmStatus) PrintNodesLA() {
	uptimes := make(map[string]string)

	for nodeName, node := range m.Connections {
		nodeSSH := node.SSH

		up := nodeSSH.ExecInShell("uptime")
		uptimes[nodeName] = up
	}

	utils.PrintSortedMap(uptimes)
}
