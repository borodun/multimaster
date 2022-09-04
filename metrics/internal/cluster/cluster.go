package cluster

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"metrics/internal/database"
	"metrics/internal/node"
	"strings"
	"time"
)

type Cluster struct {
	Nodes            []node.Node
	Interval         time.Duration
	Timeout          time.Duration
	DbMaxConnections int

	commonConnInfo string
}

func (c *Cluster) AddNode(name, connInfo string) {
	var dbLog = log.WithField("conn", name)

	if c.commonConnInfo == "" {
		c.commonConnInfo = connInfo
	}

	db, err := sql.Open("postgres", c.mergeConnInfos(c.commonConnInfo, connInfo))
	if err != nil {
		dbLog.WithError(err).Warn("failed to open url, disabling database from monitoring")
		return
	}

	if err = db.Ping(); err != nil {
		dbLog.WithError(err).Warn("failed to ping database")
	}
	db.SetMaxOpenConns(c.DbMaxConnections)

	dbx := database.NewDatabase(name, db, c.Timeout)
	nd := node.NewNode(name, dbx, c.Interval)
	c.Nodes = append(c.Nodes, nd)

	dbLog.Info("node added")

	nd.StartMonitoring()
}

func (c *Cluster) getFieldFromConnInfo(connInfo, field string) string {
	fields := strings.Split(connInfo, "")
	for _, f := range fields {
		if strings.Contains(f, field) {
			return strings.Split(f, "=")[1]
		}
	}
	return ""
}

func (c *Cluster) mergeConnInfos(connInfos ...string) string {
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
