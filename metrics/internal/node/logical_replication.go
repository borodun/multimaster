package node

import (
	"fmt"
	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type slots struct {
	Name     string  `db:"slot_name"`
	Active   float64 `db:"active"`
	TotalLag float64 `db:"total_lag"`
}

// ReplicationSlotLagInBytes returns the total lag in bytes from the replication slots
func (n *Node) ReplicationSlotLagInBytes() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_lag_bytes",
			Help:        "Total lag of the replication slots in bytes",
			ConstLabels: n.Labels,
		},
		[]string{"slot_name"},
	)

	var replicationLagQuery = fmt.Sprintf(
		`
						SELECT
							slot_name,
							%s(%s(), confirmed_flush_lsn) AS total_lag
						FROM pg_replication_slots
						WHERE slot_type = 'logical'
						AND "database" = current_database();
					`,
		postgres.Version(n.Db.Version()).WalLsnDiffFunctionName(),
		postgres.Version(n.Db.Version()).CurrentWalLsnFunctionName())

	go func() {
		for {
			gauge.Reset()
			var slots []slots
			if err := n.Db.Query(replicationLagQuery, &slots); err == nil {
				for _, slot := range slots {
					gauge.With(prometheus.Labels{
						"slot_name": slot.Name,
					}).Set(slot.TotalLag)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}
