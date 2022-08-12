package metrics

import (
	"database/sql"
	"metrics/internal/config"
	"metrics/internal/gauges"
	"net/http"
	nurl "net/url"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
)

func Run(cfg config.Config) {
	var server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Spec.ListenPort),
	}

	http.Handle("/metrics", promhttp.Handler())

	for _, con := range cfg.Spec.Databases {
		var dbLog = log.WithField("db", con.Name)

		url, err := pgpass.UpdateURL(con.URL)
		if err != nil {
			dbLog.WithError(err).Error("failed to add password from ~/.pgpass")
		}
		if !checkPass(url) {
			dbLog.Error("no password provided for user")
		}

		db, err := sql.Open("postgres", url)
		if err != nil {
			dbLog.WithError(err).Warn("failed to open url, disabling database from monitoring")
			continue
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			dbLog.WithError(err).Warn("failed to ping database")
		}
		db.SetMaxOpenConns(cfg.Spec.ConnectionPoolMaxSize)

		labels := mergeLabels(cfg.Spec.Labels, con.Labels, map[string]string{"connection_name": con.Name})

		watch(db, prometheus.DefaultRegisterer, con.Name, cfg.Spec.Interval, cfg.Spec.QueryTimeout, labels)
		dbLog.Info("started monitoring")
	}

	log.WithField("port", strconv.Itoa(cfg.Spec.ListenPort)).Info("server started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
}

func mergeLabels(labelsToMerge ...map[string]string) map[string]string {
	labels := map[string]string{}

	for i := range labelsToMerge {
		for k, v := range labelsToMerge[i] {
			labels[k] = v
		}
	}

	return labels
}

func checkPass(url string) bool {
	u, err := nurl.Parse(url)
	if err != nil {
		return false
	}
	if user := u.User; user != nil {
		if _, ok := user.Password(); !ok {
			return false
		}
	}
	return true
}

func watch(db *sql.DB, reg prometheus.Registerer, conName string, interval, timeout int, labels map[string]string) {
	var postgresGauges = gauges.New(db, conName, time.Duration(interval)*time.Second, time.Duration(timeout)*time.Second, labels)
	reg.MustRegister(postgresGauges.Errs)

	reg.MustRegister(postgresGauges.ConnectedBackends())
	reg.MustRegister(postgresGauges.MaxBackends())
	reg.MustRegister(postgresGauges.InstanceConnectedBackends())
	reg.MustRegister(postgresGauges.BackendsByState())
	reg.MustRegister(postgresGauges.BackendsByUserAndClientAddress())
	reg.MustRegister(postgresGauges.BackendsByWaitEventType())
	reg.MustRegister(postgresGauges.RequestedCheckpoints())
	reg.MustRegister(postgresGauges.ScheduledCheckpoints())
	reg.MustRegister(postgresGauges.BuffersMaxWrittenClean())
	reg.MustRegister(postgresGauges.BuffersWritten())
	reg.MustRegister(postgresGauges.HeapBlocksHit())
	reg.MustRegister(postgresGauges.HeapBlocksRead())
	reg.MustRegister(postgresGauges.IndexScans())
	reg.MustRegister(postgresGauges.UnusedIndexes())
	reg.MustRegister(postgresGauges.IndexBlocksReadBySchema())
	reg.MustRegister(postgresGauges.IndexBlocksHitBySchema())
	reg.MustRegister(postgresGauges.IndexBloat())
	reg.MustRegister(postgresGauges.Locks())
	reg.MustRegister(postgresGauges.NotGrantedLocks())
	reg.MustRegister(postgresGauges.DeadLocks())
	reg.MustRegister(postgresGauges.ReplicationDelayInSeconds())
	reg.MustRegister(postgresGauges.ReplicationDelayInBytes())
	reg.MustRegister(postgresGauges.ReplicationStatus())
	reg.MustRegister(postgresGauges.Size())
	reg.MustRegister(postgresGauges.StreamingWALs())
	reg.MustRegister(postgresGauges.TableBloat())
	reg.MustRegister(postgresGauges.TableUsage())
	reg.MustRegister(postgresGauges.TempFiles())
	reg.MustRegister(postgresGauges.TempSize())
	reg.MustRegister(postgresGauges.TransactionsSum())
	reg.MustRegister(postgresGauges.Up())
	reg.MustRegister(postgresGauges.TableScans())
	reg.MustRegister(postgresGauges.TableSizes())
	reg.MustRegister(postgresGauges.DatabaseReadingUsage())
	reg.MustRegister(postgresGauges.DatabaseWritingUsage())
	reg.MustRegister(postgresGauges.HOTUpdates())
	reg.MustRegister(postgresGauges.TableLiveRows())
	reg.MustRegister(postgresGauges.TableDeadRows())
	reg.MustRegister(postgresGauges.DatabaseLiveRows())
	reg.MustRegister(postgresGauges.DatabaseDeadRows())
	reg.MustRegister(postgresGauges.UnvacuumedTransactions())
	reg.MustRegister(postgresGauges.LastTimeVacuumRan())
	reg.MustRegister(postgresGauges.LastTimeAutoVacuumRan())
	reg.MustRegister(postgresGauges.VacuumRunningTotal())
	reg.MustRegister(postgresGauges.ReplicationSlotStatus())
	reg.MustRegister(postgresGauges.ReplicationSlotLagInBytes())

	reg.MustRegister(postgresGauges.SlowestQueries())
	reg.MustRegister(postgresGauges.DeadTuples())

	reg.MustRegister(postgresGauges.Latency())
	reg.MustRegister(postgresGauges.TransactionsCommitSum())
	reg.MustRegister(postgresGauges.TransactionsRollbackSum())
	reg.MustRegister(postgresGauges.QueriesSum())

	reg.MustRegister(postgresGauges.MtmStatus())
	reg.MustRegister(postgresGauges.MtmGenNum())
}
