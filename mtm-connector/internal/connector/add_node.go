package connector

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (m *MtmConnector) AddNode(w http.ResponseWriter, r *http.Request) {
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

	_, ok := m.Hosts[addr]
	if ok {
		log.WithField("addr", addr).Error("add node: node already added")
		w.WriteHeader(http.StatusAlreadyReported)
		fmt.Fprintf(w, "%s already added\n", addr)
		return
	}

	id, err := m.mtmAddNodeAndGetID(host, port)
	if err != nil {
		log.WithField("addr", addr).WithError(err).Error("add node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred %s\n", err.Error())
		return
	}
	log.WithField("addr", addr).Infof("added node: id: %s", id)

	m.Hosts[addr] = id
	m.Joined[addr] = false

	fmt.Fprintln(w, m.removeFromConnInfo(m.ConnInfo, "sslmode"))
}

func (m *MtmConnector) mtmAddNodeAndGetID(host, port string) (string, error) {
	mtmAddNodeQuery := fmt.Sprintf(`SELECT mtm.add_node('%s')`, m.mergeConnInfos(m.removeFromConnInfo(m.ConnInfo, "sslmode"), "host="+host, "port="+port))

	var id []string
	err := m.Db.Query(mtmAddNodeQuery, &id)

	if err != nil {
		return "", err
	}

	return id[0], nil
}
