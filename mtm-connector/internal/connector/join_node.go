package connector

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (m *MtmConnector) JoinNode(w http.ResponseWriter, r *http.Request) {
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

	lsn := vars["lsn"]
	if lsn == "" {
		log.WithField("addr", addr).Error("join node: empty lsn")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no lsn provided")
		return
	}

	id, ok := m.Hosts[addr]
	if !ok {
		log.WithField("addr", addr).Errorf("join node: no id for '%s' host", addr)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no id for %s found, you should add node firstly\n", addr)
		return
	}

	joined := m.Joined[addr]
	if joined {
		log.WithField("addr", addr).Errorf("join node: '%s' already joined", addr)
		w.WriteHeader(http.StatusAlreadyReported)
		fmt.Fprintf(w, "%s already joined\n", addr)
		return
	}

	err := m.mtmJoinNode(id, lsn)
	if err != nil {
		log.WithField("addr", addr).WithError(err).Error("join node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred: %s\n", err.Error())
		return
	}

	log.WithField("addr", addr).Infof("joined node: id: %s", id)

	m.Joined[addr] = true

	fmt.Fprintln(w, "joined node")
}

func (m *MtmConnector) mtmJoinNode(id, lsn string) error {
	mtmJoinNodeQuery := fmt.Sprintf(`SELECT mtm.join_node(%s, '%s')`, id, lsn)

	var info []string
	err := m.Db.Query(mtmJoinNodeQuery, &info)

	if err != nil {
		return err
	}

	return nil
}
