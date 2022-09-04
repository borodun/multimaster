package node

import "github.com/prometheus/client_golang/prometheus"

func (n *Node) TransactionsSum() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_sum",
			Help:        "Sum of all transactions in the database",
			ConstLabels: n.Labels,
		},
		`
			SELECT xact_commit + xact_rollback
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func (n *Node) TransactionsCommitSum() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_commit_sum",
			Help:        "Sum of commit transactions in the database",
			ConstLabels: n.Labels,
		},
		`
			SELECT xact_commit
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func (n *Node) TransactionsRollbackSum() prometheus.Gauge {
	return n.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_rollback_sum",
			Help:        "Sum of rollback transactions in the database",
			ConstLabels: n.Labels,
		},
		`
			SELECT xact_rollback
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}
