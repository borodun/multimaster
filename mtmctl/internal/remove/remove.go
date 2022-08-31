package remove

import (
	"backup/internal/config"
	"backup/internal/connection"
	"backup/internal/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type MtmRemoveNode struct {
	Cfg         config.Config
	Connections connection.Connections
	InitNode    string
	RemoveNode  string

	initDb   *connection.DB
	removeDb *connection.DB
}

func (m *MtmRemoveNode) Run() {
	m.initDb = m.Connections[m.InitNode].DB
	m.removeDb = m.Connections[m.RemoveNode].DB

	log.RegisterExitHandler(func() {
		fmt.Println("Something went wrong, add --verbose flag to see logs")
	})

	fmt.Println("Checking connections")
	utils.CheckDatabases(m.initDb, m.removeDb)

	removeId := m.removeDb.GetMtmNodeID()
	fmt.Printf("Id of '%s': %s\n", m.RemoveNode, removeId)

	m.mtmRemoveNode(removeId)
	fmt.Printf("Successfully removed '%s' from cluster\n", m.RemoveNode)
}

func (m *MtmRemoveNode) mtmRemoveNode(id string) {
	mtmDropNodeQuery := fmt.Sprintf(`SELECT COALESCE(drop_node, 'null') FROM mtm.drop_node(%s)`, id)

	var info []string
	err := m.removeDb.Query(mtmDropNodeQuery, &info)

	if err != nil {
		fmt.Printf("Cannot drop node: %s\n", err.Error())
		log.WithError(err).Fatal("cannot drop node")
	}
}
