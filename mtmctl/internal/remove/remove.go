package remove

import (
	"backup/internal/config"
	"backup/internal/connection"
	"backup/internal/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
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

	err = m.dropCleanUp(id)
	if err != nil {
		fmt.Printf("Error with cleaning up after drop: %s\n", err.Error())
		log.WithError(err).Fatal("clean up error")
	}
}

func (m *MtmRemoveNode) GetMtmNodes() []*connection.DB {
	var nodes []*connection.DB

	for _, node := range m.Connections {
		db := node.DB
		println("Checking status of", db.GetName())
		if stat := db.MtmStatus(); stat == "online" && db.GetName() != m.removeDb.GetName() {
			nodes = append(nodes, node.DB)
		}
	}

	return nodes
}

func (m *MtmRemoveNode) dropCleanUp(id string) error {
	nodes := m.GetMtmNodes()

	println("Got nodes for clean up:")
	for _, node := range nodes {
		print(node.GetName(), " ")
	}
	println("")

	initDoneQuery := fmt.Sprintf("DELETE FROM mtm.nodes_init_done WHERE id = %s", id)

	for _, node := range nodes {
		var info []string
		println("Cleaning init done for", node.GetName())
		err := node.Query(initDoneQuery, &info)
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Duration(2) * time.Second)
	println("Cleaning syncpoints")
	syncpointsQuery := fmt.Sprintf("DELETE FROM mtm.syncpoints WHERE receiver_node_id = %s OR origin_node_id = %s", id, id)
	var info []string
	err := m.initDb.Query(syncpointsQuery, &info)
	if err != nil {
		return err
	}

	return nil
}
