package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type readingUsage struct {
	TuplesReturned float64 `db:"tup_returned"`
	TuplesFetched  float64 `db:"tup_fetched"`
}

func (g *Gauges) DatabaseReadingUsage() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_database_reading_usage",
			Help:        "Database reading usage statistics",
			ConstLabels: g.labels,
		},
		[]string{"stat"},
	)

	const databaseReadingUsageQuery = `
		SELECT COALESCE(tup_returned, 0) as tup_returned, 
			COALESCE(tup_fetched, 0) as tup_fetched
		  FROM pg_stat_database 
		 WHERE datname = current_database()
	`

	go func() {
		for {
			var readingUsage []readingUsage
			if err := g.query(databaseReadingUsageQuery, &readingUsage, emptyParams); err == nil {
				for _, database := range readingUsage {
					gauge.With(prometheus.Labels{
						"stat": "tup_fetched",
					}).Set(database.TuplesFetched)
					gauge.With(prometheus.Labels{
						"stat": "tup_returned",
					}).Set(database.TuplesReturned)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
