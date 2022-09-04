package cluster

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"metrics/internal/node"
	"time"
)

type NodeStatus struct {
	Name   string
	Id     int
	Status string
}

func (c *Cluster) StartNodeDiscovery() {
	go func() {
		c.discoverNodes()
		time.Sleep(c.Interval)
	}()
}

func (c *Cluster) discoverNodes() {
	var nodeStatuses []NodeStatus
	var mtmNodes []node.MtmNodes

	for _, n := range c.Nodes {
		id := n.GetID()
		status := n.GetStatus()

		if status == "online" && len(mtmNodes) == 0 {
			nodes, err := n.GetMtmNodes()
			if err == nil {
				mtmNodes = nodes
			}
		}

		ns := NodeStatus{
			Name:   n.Name,
			Id:     id,
			Status: status,
		}
		nodeStatuses = append(nodeStatuses, ns)
	}

	connected := make(map[int]bool)

	for _, mtmNode := range mtmNodes {
		for _, status := range nodeStatuses {
			if mtmNode.Id == status.Id {
				connected[mtmNode.Id] = true
				break
			}
		}
		if connected[mtmNode.Id] {
			continue
		}

		log.Infof("discovered new node, id: %d", mtmNode.Id)
		c.AddNode(fmt.Sprintf("node%d", mtmNode.Id), mtmNode.ConnInfo)
	}
}

func (c *Cluster) getNodeName(nodeStatuses []NodeStatus, id int) string {
	name := fmt.Sprintf("nodeStatus%d", id)
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

	return fmt.Sprintf("nodeStatus%d", max+1)
}
