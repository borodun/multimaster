package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Up returns if database is up and accepting connections
func (g *Gauges) Up() prometheus.Gauge {
	var gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "postgresql_up",
			Help:        "Database is up and accepting connections",
			ConstLabels: g.labels,
		},
	)
	go func() {
		for {
			var up = 1.0
			if err := g.db.Ping(); err != nil {
				up = 0.0
			}
			gauge.Set(up)
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// Latency returns round trip time to db in ms
func (g *Gauges) Latency() prometheus.Gauge {
	var gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "postgresql_latency",
			Help:        "Database RTT to db in ms",
			ConstLabels: g.labels,
		},
	)

	go func() {
		for {
			start := time.Now()
			var ret []string
			err := g.query("SELECT 1", &ret, emptyParams)
			end := time.Now()
			if err == nil {
				gauge.Set(float64(end.Sub(start).Milliseconds()))
			} else {
				gauge.Set(0)
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

// Size returns the database size in bytes
func (g *Gauges) Size() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_size_bytes",
			Help:        "Dabatase size in bytes",
			ConstLabels: g.labels,
		},
		"SELECT pg_database_size(current_database())",
	)
}

// TempSize returns the database total amount of data written to temporary files in bytes
func (g *Gauges) TempSize() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_temp_bytes",
			Help:        "Database total amount of data written to temporary files in bytes",
			ConstLabels: g.labels,
		},
		"SELECT temp_bytes FROM pg_stat_database WHERE datname = current_database()",
	)
}

// TempFiles returns the number of temporary files created by queries in this database.
func (g *Gauges) TempFiles() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_temp_files",
			Help:        "Number of temporary files created by queries in this database.",
			ConstLabels: g.labels,
		},
		"SELECT temp_files FROM pg_stat_database WHERE datname = current_database()",
	)
}
