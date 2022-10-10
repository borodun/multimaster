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
	if host == "" {
		log.Error("add node: empty host")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no host provided")
		return
	}

	id, err := m.mtmAddNodeAndGetID(host)
	if err != nil {
		log.WithField("host", host).WithError(err).Error("add node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred %s\n", err.Error())
		return
	}
	log.WithField("host", host).Infof("added node: id: %s", id)

	m.InProcess[host] = id
	m.Joined[host] = false

	fmt.Fprintln(w, m.removeFromConnInfo(m.ConnInfo, "sslmode"))
}

func (m *MtmConnector) mtmAddNodeAndGetID(host string) (string, error) {
	mtmAddNodeQuery := fmt.Sprintf(`SELECT mtm.add_node('%s')`, m.mergeConnInfos(m.removeFromConnInfo(m.ConnInfo, "sslmode"), "host="+host, "port=15432"))

	var id []string
	err := m.Db.Query(mtmAddNodeQuery, &id)

	if err != nil {
		return "", err
	}

	return id[0], nil
}
