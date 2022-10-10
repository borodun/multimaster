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
	if host == "" {
		log.Error("join node: empty host")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no host provided")
		return
	}

	lsn := vars["lsn"]
	if lsn == "" {
		log.WithField("host", host).Error("join node: empty host")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no lsn provided")
		return
	}

	id, ok := m.InProcess[host]
	if !ok {
		log.WithField("host", host).Errorf("join node: no id for '%s' host", host)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no id for %s found, you should add node firstly\n", host)
		return
	}
	_, ok = m.Joined[host]
	if !ok {
		log.WithField("host", host).Errorf("join node: '%s' already joined", host)
		w.WriteHeader(http.StatusAlreadyReported)
		fmt.Fprintf(w, "%s already joined\n", host)
		return
	}

	err := m.mtmJoinNode(id, lsn)
	if err != nil {
		log.WithField("host", host).WithError(err).Error("join node")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred: %s\n", err.Error())
		return
	}

	log.WithField("host", host).Infof("joined node: id: %s", id)

	delete(m.InProcess, host)
	m.Joined[host] = true

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
