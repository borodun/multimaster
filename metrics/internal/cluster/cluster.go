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
	var dbLog = log.WithField("name", name)

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

	nd.StartMonitoring()
}

func (c *Cluster) RemoveNode(id int) {
	var newNodes []node.Node
	for _, n := range c.Nodes {
		if n.GetID() == id {
			n.StopMonitoring()
			break
		}
		newNodes = append(newNodes, n)
	}
	c.Nodes = newNodes
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

func (c *Cluster) getConnectedNodes() []node.MtmNode {
	var nodes []node.MtmNode
	for _, n := range c.Nodes {
		mn := node.MtmNode{
			Name:   n.Name,
			Id:     n.GetID(),
			Status: n.GetStatus(),
		}
		nodes = append(nodes, mn)
	}
	return nodes
}

func (c *Cluster) getActualNodes() (*[]node.MtmNode, error) {
	onlineNodes := false
	var err error
	for _, n := range c.Nodes {
		status := n.GetStatus()

		if status == "online" {
			onlineNodes = true
			nodes, err := n.GetMtmNodes()
			if err == nil {
				return &nodes, nil
			}
		}
	}

	if !onlineNodes {
		return nil, fmt.Errorf("no online nodes")
	}

	return nil, err
}
