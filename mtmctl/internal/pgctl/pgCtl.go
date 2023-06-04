package pgctl

import (
	"fmt"
	"mtmctl/internal/config"
	"mtmctl/internal/connection"
	"mtmctl/internal/utils"
)

type PGCtl struct {
	Cfg         config.Config
	Connections connection.Connections
	Node        string
	Action      string
}

func (m *PGCtl) Run() {
	if m.Node != "" {
		m.Connections = utils.GetClusterNodes(m.Connections, m.Node)
	}

	switch m.Action {
	case "start":
		m.startNodes()
	case "stop":
		m.stopNodes()
	}

}

func (m *PGCtl) startNodes() {
	for _, name := range m.Connections.GetConnNames() {
		nodeSSH := m.Connections[name].SSH
		nodeSSH.PgCtlStart()
		fmt.Printf("Node '%s' started\n", name)
	}
}

func (m *PGCtl) stopNodes() {
	for _, name := range m.Connections.GetConnNames() {
		nodeSSH := m.Connections[name].SSH
		nodeSSH.PgCtlStop()
		fmt.Printf("Node '%s' stopped\n", name)
	}
}
