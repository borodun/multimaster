package node

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type writingUsage struct {
	TuplesInserted float64 `db:"tup_inserted"`
	TuplesUpdated  float64 `db:"tup_updated"`
	TuplesDeleted  float64 `db:"tup_deleted"`
}

func (n *Node) DatabaseWritingUsage() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_database_writing_usage",
			Help:        "Number of inserted, updated and deleted rows per database",
			ConstLabels: n.Labels,
		},
		[]string{"stat"},
	)

	const databaseWritingUsageQuery = `
		SELECT COALESCE(tup_inserted, 0) as tup_inserted, 
			COALESCE(tup_updated, 0) as tup_updated,
			COALESCE(tup_deleted, 0) as tup_deleted
		  FROM pg_stat_database 
			WHERE datname = current_database()
	`

	go func() {
		for {
			var writingUsage []writingUsage
			if err := n.Db.Query(databaseWritingUsageQuery, &writingUsage); err == nil {
				for _, database := range writingUsage {
					gauge.With(prometheus.Labels{
						"stat": "tup_inserted",
					}).Set(database.TuplesInserted)
					gauge.With(prometheus.Labels{
						"stat": "tup_updated",
					}).Set(database.TuplesUpdated)
					gauge.With(prometheus.Labels{
						"stat": "tup_deleted",
					}).Set(database.TuplesDeleted)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}
