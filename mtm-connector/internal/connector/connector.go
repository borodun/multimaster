package connector

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"mtm-connector/internal/database"
	"net/http"
	"strings"
	"time"
)

type MtmConnector struct {
	Db       database.Database
	ConnInfo string
	Hosts    map[string]string
	Joined   map[string]bool
}

func (m *MtmConnector) Serve(port string) {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/add-node", m.AddNode).
		Queries("host", "{host}")
	r.HandleFunc("/api/v1/join-node", m.JoinNode).
		Queries("lsn", "{lsn}", "host", "{host}")
	r.HandleFunc("/api/v1/drop-node", m.DropNode).
		Queries("host", "{host}")

	srv := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
		IdleTimeout:  600 * time.Second,
		Handler:      r,
	}

	log.Infof("listening on localhost:%s", port)
	log.WithError(srv.ListenAndServe()).Fatal("error occurred while serving")
}

func (m *MtmConnector) mergeConnInfos(connInfos ...string) string {
	connInfo := make(map[string]string)

	for _, info := range connInfos {
		fields := strings.Split(info, " ")
		for _, field := range fields {
			keyValue := strings.Split(field, "=")
			connInfo[keyValue[0]] = keyValue[1]
		}
	}

	var ret []string
	for key, value := range connInfo {
		ret = append(ret, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(ret, " ")
}

func (m *MtmConnector) removeFromConnInfo(connInfo string, fieldsToRemove ...string) string {
	retConnInfo := make(map[string]string)

	connInfoFields := strings.Split(connInfo, " ")
	for _, f := range connInfoFields {
		keyValue := strings.Split(f, "=")
		if arrContains(fieldsToRemove, keyValue[0]) {
			continue
		}
		retConnInfo[keyValue[0]] = keyValue[1]
	}

	var ret []string
	for key, value := range retConnInfo {
		ret = append(ret, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(ret, " ")
}

func arrContains(arr []string, name string) bool {
	for _, value := range arr {
		if value == name {
			return true
		}
	}
	return false
}
