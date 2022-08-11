package gauges

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

func (g *Gauges) checkPgStatStatements() {
	if !g.hasExtension("pg_stat_statements") {
		log.WithField("db", g.name).
			Warn("postgresql_slowest_queries disabled because pg_stat_statements extension is not installed")
		return
	}
	if !g.hasSharedPreloadLibrary("pg_stat_statements") {
		log.WithField("db", g.name).
			Warn("postgresql_slowest_queries disabled because pg_stat_statements is not on shared_preload_libraries")
		return
	}
	g.pgStatStatements = true
}

type slowQuery struct {
	Query string  `db:"query"`
	Time  float64 `db:"time_per_call"`
}

// SlowestQueries returns 10 slowest queries by average time per call in the database
func (g *Gauges) SlowestQueries() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_slowest_queries",
			Help:        "top 10 slowest queries by time per call",
			ConstLabels: g.labels,
		},
		[]string{"query"},
	)

	if !g.pgStatStatements {
		return gauge
	}

	var slowestQueriesQuery = fmt.Sprintf(`
					SELECT query, %[1]s / calls as time_per_call
					FROM pg_stat_statements
					WHERE dbid = (SELECT datid FROM pg_stat_database WHERE datname = current_database())
					ORDER BY %[1]s / calls desc limit 10`,
		postgres.Version(g.version()).PgStatStatementsTotalTimeColumn())

	go func() {
		for {
			var queries []slowQuery
			if err := g.query(slowestQueriesQuery, &queries, emptyParams); err == nil {
				for _, query := range queries {
					gauge.With(prometheus.Labels{
						"query": strings.Join(strings.Fields(query.Query), " "),
					}).Set(query.Time)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// QueriesSum returns total number of executed queries in the database
func (g *Gauges) QueriesSum() prometheus.Gauge {
	var gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "postgresql_queries_sum",
		Help:        "Number of all queries in the database",
		ConstLabels: g.labels,
	})

	if !g.pgStatStatements {
		return gauge
	}

	const queriesSumQuery = `SELECT sum(calls) FROM pg_stat_statements`

	go func() {
		for {
			var sum []float64
			if err := g.query(queriesSumQuery, &sum, emptyParams); err == nil {
				gauge.Set(sum[0])
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}
