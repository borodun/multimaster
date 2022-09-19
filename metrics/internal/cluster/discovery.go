package cluster

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"metrics/internal/node"
	"time"
)

func (c *Cluster) StartNodeDiscovery() {
	go func() {
		for {
			c.discoverNodes()
			time.Sleep(c.Interval)
		}
	}()
}

func (c *Cluster) discoverNodes() {
	connectedNodes := c.getConnectedNodes()
	actualNodes, err := c.getActualNodes()
	if err != nil {
		log.WithError(err).Warn("discovery: couldn't retrieve nodes in cluster")
		return
	}
	if len(*actualNodes) == 0 {
		log.Warn("discovery: no nodes retrieved")
		return
	}

	// Adding nodes
	for _, mtmNode := range *actualNodes {
		found := findById(connectedNodes, mtmNode.Id)
		if found {
			continue
		}
		if !(mtmNode.Connected && mtmNode.Enabled) {
			continue
		}

		name := c.getNodeName(connectedNodes, mtmNode.Id)
		log.WithField("id", mtmNode.Id).
			WithField("name", name).Infof("discovered new node")
		c.AddNode(name, mtmNode.ConnInfo)
	}

	// Removing nodes
	for _, connectedNode := range connectedNodes {
		found := findById(*actualNodes, connectedNode.Id)
		if found {
			continue
		}

		log.Infof("node '%s' not in the cluster, removing", connectedNode.Name)
		c.RemoveNode(connectedNode.Id)
	}
}

func findById(nodes []node.MtmNode, id int) bool {
	for _, n := range nodes {
		if n.Id == id {
			return true
		}
	}
	return false
}

func (c *Cluster) getNodeName(nodeStatuses []node.MtmNode, id int) string {
	name := fmt.Sprintf("node%d", id)
	nameTaken := false
	max := 0
	for _, nodeStatus := range nodeStatuses {
		if nodeStatus.Name == name {
			nameTaken = true
		}
		if nodeStatus.Id > max {
			max = nodeStatus.Id
		}
	}

	if !nameTaken {
		return name
	}

	newName := fmt.Sprintf("node%d", max+1)
	log.Infof("node name '%s' is already taken, new node will have '%s' name", name, newName)

	return newName
}
