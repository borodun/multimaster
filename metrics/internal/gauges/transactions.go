package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) TransactionsSum() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_sum",
			Help:        "Sum of all transactions in the database",
			ConstLabels: g.labels,
		},
		`
			SELECT xact_commit + xact_rollback
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func (g *Gauges) TransactionsCommitSum() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_commit_sum",
			Help:        "Sum of commit transactions in the database",
			ConstLabels: g.labels,
		},
		`
			SELECT xact_commit
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func (g *Gauges) TransactionsRollbackSum() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_rollback_sum",
			Help:        "Sum of rollback transactions in the database",
			ConstLabels: g.labels,
		},
		`
			SELECT xact_rollback
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}
