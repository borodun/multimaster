package connection

import (
	"fmt"
	"github.com/k0sproject/rig"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SSH struct {
	rig.Connection
	name   string
	pgdata string
	pgbin  string
}

func (s *SSH) PgCtl(cmd string) (string, error) {
	c := fmt.Sprintf("%s/pg_ctl -D %s -l %s/logfile %s 2>&1", s.pgbin, s.pgdata, s.pgdata, cmd)
	out, err := s.ExecOutput(c)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("cmd", c).
			WithField("out", out).
			WithError(err).Error("pg_ctl error")
	}
	return out, err
}

func (s *SSH) PgRunning() bool {
	status, _ := s.PgCtl("status")
	if strings.Contains(status, "server is running") {
		return true
	} else if strings.Contains(status, "no server running") {
		return false
	}
	return false
}

func (s *SSH) PgCtlStop() {
	running := s.PgRunning()
	if running {
		_, err := s.PgCtl("stop")
		if err != nil {
			log.WithError(err).Warn("error occurred while trying to stop Postgres")
		}
		log.Info("stopped Postgres")
	}
}

func (s *SSH) PgCtlStart() {
	running := s.PgRunning()
	if !running {
		_, err := s.PgCtl("start")
		if err != nil {
			log.WithError(err).Warn("error occurred while trying to start Postgres")
		}
		log.Info("started Postgres")
	}
}

func (s *SSH) RemovePGDATA() {
	out, err := s.ExecOutputf("rm -rf %s 2>&1", s.pgdata)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithError(err).Error("remove pgdata error")
	}
}
