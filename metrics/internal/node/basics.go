package node

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Status struct {
	Id     int
	Status string
}

// Up returns if database is up and accepting connections
func (n *Node) Up() prometheus.Gauge {
	var gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "postgresql_up",
			Help:        "Database is up and accepting connections",
			ConstLabels: n.Labels,
		},
	)
	if n.ctx.Err() != nil {
		return gauge
	}

	go func() {
		for {
			if n.ctx.Err() != nil {
				return
			}

			var up = 1.0
			if err := n.Db.Ping(); err != nil {
				up = 0.0
			}
			gauge.Set(up)
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}

// Latency returns round trip time to db in ms
func (n *Node) Latency() prometheus.Gauge {
	var gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "postgresql_latency",
			Help:        "Database RTT to db in ms",
			ConstLabels: n.Labels,
		},
	)
	if n.ctx.Err() != nil {
		return gauge
	}

	go func() {
		for {
			if n.ctx.Err() != nil {
				return
			}

			start := time.Now()
			var ret []string
			err := n.Db.Query("SELECT 1", &ret)
			end := time.Now()
			if err == nil {
				gauge.Set(float64(end.Sub(start).Milliseconds()))
			} else {
				gauge.Set(0)
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}

// Size returns the database size in bytes
func (n *Node) Size() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_size_bytes",
			Help:        "Database size in bytes",
			ConstLabels: n.Labels,
		},
		"SELECT pg_database_size(current_database())",
	)
}
