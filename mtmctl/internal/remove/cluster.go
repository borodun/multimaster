package remove

import (
	"fmt"
	"mtmctl/internal/config"
	"mtmctl/internal/connection"
	"mtmctl/internal/utils"
	"sync"
)

type MtmRemoveCluster struct {
	Cfg         config.Config
	Connections connection.Connections
	Nodes       []string
}

func (m *MtmRemoveCluster) Run() {
	fmt.Print("Checking all nodes: ")
	if !m.checkAllNodes() {
		return
	}

	clearData := utils.AskUser("Do you want to clear databases data?")

	var wg sync.WaitGroup

	for _, nodename := range m.Nodes {
		wg.Add(1)

		go func(node string) {
			defer wg.Done()

			nodeSSH := m.Connections[node].SSH

			nodeSSH.PgCtlStop()
			fmt.Printf("Stopped database on '%s'\n", node)

			if clearData {
				nodeSSH.RemovePGDATA()
				fmt.Printf("Removed database on '%s'\n", node)
			}
		}(nodename)
	}

	wg.Wait()

	fmt.Println("Cluster removed")
}

func (m *MtmRemoveCluster) checkAllNodes() bool {
	nodesCheck := make(map[string]string)
	nodesCount := 0

	for _, node := range m.Nodes {
		nodeDB := m.Connections[node].DB

		if stat := nodeDB.MtmStatus(); stat != "online" {
			fmt.Printf("'%s' is not online, current status: %s \n", node, stat)
			return false
		}

		nodes := nodeDB.GetMtmNodes()

		if nodesCount == 0 {
			nodesCount = len(nodes)
		} else {
			if nodesCount != len(nodes) {
				fmt.Printf("nodes count doesnt match on '%s'\n", node)
				return false
			}
		}

		for _, nodeTup := range nodes {
			id, conninfo := nodeTup.Id, nodeTup.Conninfo
			val, ok := nodesCheck[id]

			if !ok {
				nodesCheck[id] = conninfo
			} else {
				if conninfo != val {
					fmt.Printf("node '%s' have diffent conninfo '%s' that doesn't match with others '%s'\n", node, conninfo, val)
					return false
				}
			}
		}
	}

	if nodesCount != len(m.Nodes) {
		fmt.Printf("wrong nodes specified\n")
		return false
	}

	fmt.Println("ok")
	return true
}
