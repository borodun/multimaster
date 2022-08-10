package gauges

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (g *Gauges) CheckMtm() {
	if !g.hasExtension("multimaster") {
		log.WithField("db", g.name).
			Warn("mtm monitoring is disabled because multimaster extension is not installed")
	}
}

type status struct {
	Status string `db:"status"`
}

// MtmStatus returns status of a node in mtm cluster
func (g *Gauges) MtmStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mtm_status",
		Help:        "Node status in mtm cluster",
		ConstLabels: g.labels,
	}, []string{"status"})

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
		for {
			var status []status
			err := g.query(mtmNodeStatusQuery, &status, emptyParams)

			for _, statusType := range statusTypes {
				gauge.With(prometheus.Labels{
					"status": statusType,
				}).Set(0)

				if err != nil && strings.Contains(strings.Split(err.Error(), ":")[2], statusType) ||
					err == nil && status[0].Status == statusType {
					gauge.With(prometheus.Labels{
						"status": statusType,
					}).Set(1)
				}
			}

			time.Sleep(g.interval)
		}
	}()

	return gauge
}

// MtmGenNum returns generation of a node in mtm cluster
func (g *Gauges) MtmGenNum() prometheus.Gauge {
	return g.new(prometheus.GaugeOpts{
		Name:        "mtm_gen_num",
		Help:        "Node generation in mtm cluster",
		ConstLabels: g.labels,
	}, "SELECT gen_num FROM mtm.status()")
}
