package connector

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"mtm-connector/internal/database"
	"net/http"
)

func (m *MtmConnector) DropNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	host := vars["host"]
	port := vars["port"]
	addr := host + ":" + port
	if host == "" || port == "" {
		log.Error("add node: empty host or port")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "empty host or port provided")
		return
	}

	id, ok := m.Hosts[addr]
	if !ok {
		log.WithField("addr", addr).Errorf("drop node: no id for '%s' host", addr)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no id for %s found, you should add node firstly\n", addr)
		return
	}

	//joined := m.Joined[addr]
	//if !joined {
	//	log.WithField("addr", addr).Errorf("drop node: node '%s' not joined", addr)
	//	w.WriteHeader(http.StatusNotFound)
	//	fmt.Fprintf(w, "node %s not joined, you should join node firstly\n", addr)
	//	return
	//}

	err := m.mtmDropNode(id)
	if err != nil {
		log.WithField("addr", addr).WithError(err).Error("drop node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred %s\n", err.Error())
		delete(m.Hosts, addr)
		delete(m.Joined, addr)
		return
	}

	log.WithField("addr", addr).Infof("dropped node: id: %s", id)

	delete(m.Hosts, addr)
	delete(m.Joined, addr)

	log.Infof("starting clean up after dropping node %s", id)

	err = m.dropCleanUp(id)
	if err != nil {
		log.WithField("addr", addr).WithError(err).Error("clean up")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "clean up: error occurred %s\n", err.Error())
		return
	}

	log.Infof("clean up succesfull")

	fmt.Fprintln(w, "node dropped")
}

func (m *MtmConnector) mtmDropNode(id string) error {
	mtmDropNodeQuery := fmt.Sprintf(`SELECT COALESCE(drop_node, 'null') FROM mtm.drop_node(%s)`, id)

	var info []string
	err := m.Db.Query(mtmDropNodeQuery, &info)

	if err != nil {
		return err
	}

	return nil
}

func (m *MtmConnector) dropCleanUp(id string) error {
	nodes, err := m.GetMtmNodes()
	if err != nil {
		return err
	}

	var conns []database.Database

	for _, node := range nodes {
		//if !(node.Connected && node.Enabled) {
		//	continue
		//}

		conns = append(conns, database.NewDatabase(m.mergeConnInfos(node.ConnInfo, "sslmode=disable")))
	}

	initDoneQuery := fmt.Sprintf("DELETE FROM mtm.nodes_init_done WHERE id = %s", id)

	for _, conn := range conns {
		var info []string
		err = conn.Query(initDoneQuery, &info)
		if err != nil {
			return err
		}
	}

	syncpointsQuery := fmt.Sprintf("DELETE FROM mtm.syncpoints WHERE receiver_node_id = %s OR origin_node_id = %s", id, id)
	var info []string
	err = m.Db.Query(syncpointsQuery, &info)
	if err != nil {
		return err
	}

	return nil
}

type MtmNode struct {
	Id        int    `db:"id"`
	ConnInfo  string `db:"conninfo"`
	Enabled   bool   `db:"enabled"`
	Connected bool   `db:"connected"`
	Name      string
	Status    string
}

func (m *MtmConnector) GetMtmNodes() ([]MtmNode, error) {
	const mtmNodesQuery = `SELECT id, conninfo, enabled, connected FROM mtm.nodes()`

	var nodes []MtmNode
	err := m.Db.Query(mtmNodesQuery, &nodes)

	return nodes, err
}
