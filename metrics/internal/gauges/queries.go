package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) QueriesSum() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_queries_sum",
			Help:        "Sum of all queries calls in the database",
			ConstLabels: g.labels,
		}, `SELECT sum(calls) FROM pg_stat_statements`)
}
