package query

import (
	"fmt"
	"mtmctl/internal/config"
	"mtmctl/internal/connection"
	"mtmctl/internal/utils"
	"sync"
)

type MtmQuery struct {
	Cfg         config.Config
	Connections connection.Connections
	Node        string
	Query       string
}

func (m *MtmQuery) Run() {
	if m.Node != "" {
		m.Connections = utils.GetClusterNodes(m.Connections, m.Node)
	}

	var wg sync.WaitGroup

	for nodeName, node := range m.Connections {
		wg.Add(1)

		go func(name string, node connection.Connection) {
			defer wg.Done()

			nodeSSH := node.SSH
			out := nodeSSH.RunPSQL_MTM_Full(m.Query)
			fmt.Printf("Node '%s' result: \n%s\n\n", name, out)
		}(nodeName, node)
	}

	wg.Wait()
}
