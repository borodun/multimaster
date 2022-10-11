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
	if host == "" {
		log.Error("drop node: empty host")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no host provided")
		return
	}

	id, ok := m.Hosts[host]
	if !ok {
		log.WithField("host", host).Errorf("drop node: no id for '%s' host", host)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no id for %s found, you should add node firstly\n", host)
		return
	}

	joined := m.Joined[host]
	if !joined {
		log.WithField("host", host).Errorf("drop node: node '%s' not joined", host)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "node %s not joined, you should join node firstly\n", host)
		return
	}

	err := m.mtmDropNode(id)
	if err != nil {
		log.WithField("host", host).WithError(err).Error("drop node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred %s\n", err.Error())
		delete(m.Hosts, host)
		delete(m.Joined, host)
		return
	}

	log.WithField("host", host).Infof("dropped node: id: %s", id)

	delete(m.Hosts, host)
	delete(m.Joined, host)

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
