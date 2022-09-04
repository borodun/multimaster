package metrics

import (
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
	"metrics/internal/cluster"
	"metrics/internal/config"
	"net/http"
	nurl "net/url"
	"strconv"
	"time"
)

func Run(cfg config.Config) {
	var server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Metrics.ListenPort),
	}

	http.Handle("/metrics", promhttp.Handler())

	cl := cluster.Cluster{
		DbMaxConnections: cfg.Metrics.ConnectionPoolMaxSize,
		Interval:         time.Duration(cfg.Metrics.Interval) * time.Second,
		Timeout:          time.Duration(cfg.Metrics.QueryTimeout) * time.Second,
	}

	for _, con := range cfg.Metrics.Databases {
		url, err := pgpass.UpdateURL(con.URL)
		if err != nil {
			log.WithField("conn", con.Name).
				WithError(err).Error("failed to add password from ~/.pgpass")
		}
		if !checkPass(url) {
			log.WithField("conn", con.Name).
				Error("no password provided for user")
		}

		connInfo, _ := pq.ParseURL(url)
		cl.AddNode(con.Name, connInfo)
	}
	cl.StartNodeDiscovery()

	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
	log.WithField("port", strconv.Itoa(cfg.Metrics.ListenPort)).Info("server started")
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
