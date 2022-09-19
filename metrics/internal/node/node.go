package node

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"metrics/internal/database"
	"strings"
	"time"
)

type Node struct {
	Name     string
	Db       database.Database
	Interval time.Duration
	Labels   map[string]string

	mmts    bool
	removed bool
}

func NewNode(name string, db database.Database, interval time.Duration) Node {
	return Node{
		Name:     name,
		Db:       db,
		Interval: interval,
		Labels:   map[string]string{"connection_name": name},
	}
}

func (n *Node) StartMonitoring() {

	reg := prometheus.DefaultRegisterer

	n.CheckMultimaster()

	//ctx, cancel := context.WithCancel(context.Background())

	reg.MustRegister(n.MtmStatus())
	reg.MustRegister(n.MtmGenNum())
	reg.MustRegister(n.Latency())
	reg.MustRegister(n.ConnectedBackends())
	reg.MustRegister(n.Up())
	reg.MustRegister(n.ReplicationSlotLagInBytes())
	reg.MustRegister(n.DatabaseReadingUsage())
	reg.MustRegister(n.DatabaseWritingUsage())
	reg.MustRegister(n.Size())
	reg.MustRegister(n.BackendsByState())
	reg.MustRegister(n.BackendsByUserAndClientAddress())
	reg.MustRegister(n.TransactionsSum())
	reg.MustRegister(n.TransactionsCommitSum())
	reg.MustRegister(n.TransactionsRollbackSum())

	log.WithField("conn", n.Name).Info("started monitoring")
}

func (n *Node) StopMonitoring() {
	reg := prometheus.DefaultRegisterer

	reg.Unregister(n.MtmStatus())
	reg.Unregister(n.MtmGenNum())
	reg.Unregister(n.Latency())
	reg.Unregister(n.ConnectedBackends())
	reg.Unregister(n.Up())
	reg.Unregister(n.ReplicationSlotLagInBytes())
	reg.Unregister(n.DatabaseReadingUsage())
	reg.Unregister(n.DatabaseWritingUsage())
	reg.Unregister(n.Size())
	reg.Unregister(n.BackendsByState())
	reg.Unregister(n.BackendsByUserAndClientAddress())
	reg.Unregister(n.TransactionsSum())
	reg.Unregister(n.TransactionsCommitSum())
	reg.Unregister(n.TransactionsRollbackSum())

	n.removed = true
	n.Db.Disconnect()

	log.WithField("conn", n.Name).Info("stopped monitoring")
}

func paramsFix(params []string) []interface{} {
	iparams := make([]interface{}, len(params))
	for i, v := range params {
		iparams[i] = v
	}
	return iparams
}

func (n *Node) new(opts prometheus.GaugeOpts, query string, params ...string) prometheus.Gauge {
	var gauge = prometheus.NewGauge(opts)
	if n.removed {
		return gauge
	}

	go n.observe(gauge, query, paramsFix(params))
	return gauge
}

func (n *Node) observeOnce(gauge prometheus.Gauge, query string, params []interface{}) {
	log.WithField("db", n.Name).Debugf("collecting")
	var result []float64
	if err := n.Db.QueryWithParams(query, &result, params); err == nil {
		gauge.Set(result[0])
	}
}

func (n *Node) observe(gauge prometheus.Gauge, query string, params []interface{}) {
	for {
		if n.removed {
			return
		}

		n.observeOnce(gauge, query, params)
		time.Sleep(n.Interval)
	}
}

func (n *Node) CheckMultimaster() {
	if !n.Db.HasExtension("multimaster") {
		log.WithField("db", n.Name).
			Warn("metrics monitoring is disabled because multimaster extension is not installed")
		return
	}
	if !n.Db.HasSharedPreloadLibrary("multimaster") {
		log.WithField("db", n.Name).
			Warn("metrics monitoring is disabled because multimaster is not on shared_preload_libraries")
		return
	}
	n.mmts = true
}

const mtmNodeStatusQuery = `
		SELECT my_node_id, status FROM mtm.status()
	`

func (n *Node) GetStatus() string {
	var statusTypes = []string{
		"online",
		"recovery",
		"catchup",
		"disabled",
		"isolated",
	}

	var status []statusTup
	err := n.Db.Query(mtmNodeStatusQuery, &status)

	if err != nil {
		if strings.Contains(err.Error(), "pq: [MTM] multimaster node is not online: current status") {
			stat := strings.Split(err.Error(), ":")[2]
			for _, statusType := range statusTypes {
				if strings.Contains(stat, statusType) {
					return statusType
				}
			}
		} else if strings.Contains(err.Error(), "multimaster is not configured") {
			return "multimaster is not configured"
		} else if strings.Contains(err.Error(), "connection refused") {
			return "offline"
		} else {
			return fmt.Sprintf("unknown error: %s", err.Error())
		}
	}

	for _, statusType := range statusTypes {
		if status[0].Status == statusType {
			return statusType
		}
	}

	return fmt.Sprintf("unkown status: %s", status[0].Status)
}

func (n *Node) GetID() int {
	var status []statusTup
	err := n.Db.Query(mtmNodeStatusQuery, &status)

	if err != nil {
		return -1
	}
	return status[0].Id
}

type MtmNode struct {
	Id        int    `db:"id"`
	ConnInfo  string `db:"conninfo"`
	Enabled   bool   `db:"enabled"`
	Connected bool   `db:"connected"`
	Name      string
	Status    string
}

func (n *Node) GetMtmNodes() ([]MtmNode, error) {
	const mtmNodesQuery = `SELECT id, conninfo, enabled, connected FROM mtm.nodes()`

	var nodes []MtmNode
	err := n.Db.Query(mtmNodesQuery, &nodes)

	return nodes, err
}
