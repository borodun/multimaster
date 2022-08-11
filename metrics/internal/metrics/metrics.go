package metrics

import (
	"database/sql"
	"github.com/lib/pq"
	"metrics/internal/config"
	"metrics/internal/gauges"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Run(cfg config.Config) {
	var server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Spec.ListenPort),
	}

	http.Handle("/metrics", promhttp.Handler())

	for _, con := range cfg.Spec.Databases {
		var dblog = log.WithField("db", con.Name)

		db, err := sql.Open("postgres", con.URL)
		if err != nil {
			dblog.WithError(err).Fatal("failed to open url")
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			dblog.WithError(err).Warn("failed to ping database")
		}
		db.SetMaxOpenConns(cfg.Spec.ConnectionPoolMaxSize)

		dbName := getDbName(con.URL, con.Name)
		watch(db, prometheus.DefaultRegisterer, con.Name, dbName, cfg.Spec.Interval, cfg.Spec.QueryTimeout)
		dblog.Info("started monitoring")
	}

	log.WithField("port", strconv.Itoa(cfg.Spec.ListenPort)).Info("server started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
}

func getDbName(url, name string) string {
	pgUrl, _ := pq.ParseURL(url)
	params := strings.Split(pgUrl, " ")
	for _, param := range params {
		if strings.Contains(param, "dbname") {
			dbName := strings.Split(param, "=")[1]
			dbName = strings.Trim(dbName, "'")
			return dbName
		}
	}
	log.WithField("db", name).Warn("couldn't parse db name from url")
	return "null"
}

func watch(db *sql.DB, reg prometheus.Registerer, conName, dbName string, interval, timeout int) {
	var postgresGauges = gauges.New(db, conName, dbName, time.Duration(interval)*time.Second, time.Duration(timeout)*time.Second)
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
