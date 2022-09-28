package node

import (
	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// ConnectedBackends returns the number of backends currently connected to database
func (n *Node) ConnectedBackends() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_total",
			Help:        "Number of backends currently connected to database",
			ConstLabels: n.Labels,
		},
		"SELECT numbackends FROM pg_stat_database WHERE datname = current_database()",
	)
}

// MaxBackends returns the maximum number of concurrent connections in the database
func (n *Node) MaxBackends() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_max_backends",
			Help:        "Maximum number of concurrent connections in the database",
			ConstLabels: n.Labels,
		},
		"SELECT setting::numeric FROM pg_settings WHERE name = 'max_connections'",
	)
}

// InstanceConnectedBackends returns the number of backends currently connected to all databases
func (n *Node) InstanceConnectedBackends() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_instance_connected_backends",
			Help:        "Current number of concurrent connections in all databases",
			ConstLabels: n.Labels,
		},
		"SELECT sum(numbackends) FROM pg_stat_database;",
	)
}

type backendsByState struct {
	Total float64 `db:"total"`
	State string  `db:"state"`
}

// BackendsByState returns the number of backends currently connected to database by state
func (n *Node) BackendsByState() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_state_total",
			Help:        "Number of backends currently connected to database by state",
			ConstLabels: n.Labels,
		},
		[]string{"state"},
	)

	if n.ctx.Err() != nil {
		return gauge
	}

	const backendsByStateQuery = `
		SELECT COUNT(*) AS total, COALESCE(state, 'null') as state
		FROM pg_stat_activity
		WHERE datname = current_database()
		GROUP BY state
	`

	go func() {
		for {
			if n.ctx.Err() != nil {
				return
			}

			gauge.Reset()
			var backendsByState []backendsByState
			if err := n.Db.Query(backendsByStateQuery, &backendsByState); err == nil {
				for _, row := range backendsByState {
					gauge.With(prometheus.Labels{
						"state": row.State,
					}).Set(row.Total)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}

type backendsByUserAndClientAddress struct {
	Total      float64 `db:"total"`
	User       string  `db:"usename"`
	ClientAddr string  `db:"client_addr"`
}

// BackendsByUserAndClientAddress returns the number of backends currently connected
// to database by user and client address
func (n *Node) BackendsByUserAndClientAddress() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_user_total",
			Help:        "Number of backends currently connected to database by user and client address",
			ConstLabels: n.Labels,
		},
		[]string{"user", "client_addr"},
	)

	if n.ctx.Err() != nil {
		return gauge
	}

	const backendsByUserAndClientAddressQuery = `
		SELECT
		  COUNT(*) AS total,
		  usename,
		  COALESCE(client_addr, '::1') AS client_addr
		FROM pg_stat_activity
		WHERE datname = current_database()
		GROUP BY usename, client_addr
	`

	go func() {
		for {
			if n.ctx.Err() != nil {
				return
			}

			gauge.Reset()
			var backendsByUserAndClientAddress []backendsByUserAndClientAddress
			if err := n.Db.Query(backendsByUserAndClientAddressQuery, &backendsByUserAndClientAddress); err == nil {
				for _, row := range backendsByUserAndClientAddress {
					gauge.With(prometheus.Labels{
						"user":        row.User,
						"client_addr": row.ClientAddr,
					}).Set(row.Total)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}

type backendsByWaitEventType struct {
	Total         float64 `db:"total"`
	WaitEventType string  `db:"wait_event_type"`
}

func (n *Node) backendsByWaitEventTypeQuery() string {
	if postgres.Version(n.Db.Version()).IsEqualOrGreaterThan96() {
		return `
			SELECT
			  COUNT(*) AS total,
			  wait_event_type
			FROM pg_stat_activity
			WHERE wait_event_type IS NOT NULL
			  AND datname = current_database()
			GROUP BY wait_event_type
		`
	}
	return `
		SELECT
		  COUNT(*) as total,
		  'Lock' as wait_event_type
		FROM pg_stat_activity
		WHERE datname = current_database()
		  AND waiting is true
	`
}

// BackendsByWaitEventType returns the number of backends currently waiting on some event
func (n *Node) BackendsByWaitEventType() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_wait_event_type_total",
			Help:        "Number of backends currently waiting on some event",
			ConstLabels: n.Labels,
		},
		[]string{"wait_event_type"},
	)

	if n.ctx.Err() != nil {
		return gauge
	}

	go func() {
		for {
			if n.ctx.Err() != nil {
				return
			}

			gauge.Reset()
			var backendsByWaitEventType []backendsByWaitEventType
			if err := n.Db.Query(n.backendsByWaitEventTypeQuery(), &backendsByWaitEventType); err == nil {
				for _, row := range backendsByWaitEventType {
					gauge.With(prometheus.Labels{
						"wait_event_type": row.WaitEventType,
					}).Set(row.Total)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}
