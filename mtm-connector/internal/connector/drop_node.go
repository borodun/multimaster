package connector

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	joined := m.Joined[host]
	if !joined {
		log.WithField("addr", addr).Errorf("drop node: node '%s' not joined", addr)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "node %s not joined, you should join node firstly\n", addr)
		return
	}

	err := m.mtmDropNode(id)
	if err != nil {
		log.WithField("addr", addr).WithError(err).Error("drop node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred %s\n", err.Error())
		delete(m.Hosts, addr)
		delete(m.Joined, addr)
		return
	}

	log.WithField("host", addr).Infof("dropped node: id: %s", id)

	delete(m.Hosts, addr)
	delete(m.Joined, addr)

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
