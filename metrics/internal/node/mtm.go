package node

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type statusTup struct {
	Id     int    `db:"my_node_id"`
	Status string `db:"status"`
}

// MtmStatus returns status of a node in metrics cluster
func (n *Node) MtmStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mtm_status",
		Help:        "Node status in metrics cluster",
		ConstLabels: n.Labels,
	}, []string{"status"})

	if !n.mmts || n.removed {
		return gauge
	}

	const mtmNodeStatusQuery = `
		SELECT status FROM mtm.status()
	`

	var statusTypes = []string{
		"online",
		"recovery",
		"catchup",
		"disabled",
		"isolated",
	}

	go func() {
		log.Info("status: starting goroutine")
		for {
			if n.removed {
				log.WithField("name", n.Name).Info("status: returning from goroutine")
				return
			}

			var status []statusTup
			err := n.Db.Query(mtmNodeStatusQuery, &status)

			var statusFromErr string
			if err != nil {
				if strings.Contains(err.Error(), "pq: [MTM] multimaster node is not online: current status") {
					statusFromErr = strings.Split(err.Error(), ":")[2]
				}
			}

			for _, statusType := range statusTypes {
				gauge.With(prometheus.Labels{
					"status": statusType,
				}).Set(0)

				if (len(status) != 0 && status[0].Status == statusType) || strings.Contains(statusFromErr, statusType) {
					gauge.With(prometheus.Labels{
						"status": statusType,
					}).Set(1)
				}
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}

// MtmGenNum returns generation of a node in metrics cluster
func (n *Node) MtmGenNum() prometheus.Gauge {
	var gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "mtm_gen_num",
		Help:        "Node generation in metrics cluster",
		ConstLabels: n.Labels,
	})

	if !n.mmts || n.removed {
		return gauge
	}

	const genNumQuery = `SELECT gen_num FROM mtm.status()`

	go func() {
		log.Info("gen: starting goroutine")
		for {
			if n.removed {
				log.WithField("name", n.Name).Info("gen: returning from goroutine")
				return
			}

			var genNum []float64
			if err := n.Db.Query(genNumQuery, &genNum); err == nil {
				gauge.Set(genNum[0])
			}
			time.Sleep(n.Interval)
		}
	}()
	return gauge
}
