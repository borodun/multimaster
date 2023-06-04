package connection

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/k0sproject/rig"
	log "github.com/sirupsen/logrus"
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

func (s *SSH) Initdb() {
	out, err := s.ExecOutputf("%s/initdb -D %s", s.pgbin, s.pgdata)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithError(err).Error("initdb error")
	}
}

func (s *SSH) AddPostgresqlConfLine(line string) {
	out, err := s.ExecOutputf("echo %s >> %s/postgresql.conf", line, s.pgdata)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithField("line", line).
			WithError(err).Error("error while editing postgresql.conf")
	}
}

func (s *SSH) AddPgHBAConfLine(line string) {
	out, err := s.ExecOutputf("echo %s >> %s/pg_hba.conf", line, s.pgdata)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithField("line", line).
			WithError(err).Error("error while editing pg_hba.conf")
	}
}

func (s *SSH) RunPSQL(query string) string {
	out, err := s.ExecOutputf("%s/psql -d postgres -p 5432 -Atc \"%s\"", s.pgbin, query)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithField("query", query).
			WithError(err).Error("psql error")
	}

	return out
}

func (s *SSH) RunPSQL_MTM(query string) string {
	out, err := s.ExecOutputf("PGPASSWORD=1234 %s/psql -d mydb -U mtmuser -p 5432 -Atc \"%s\"", s.pgbin, query)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithField("query", query).
			WithError(err).Error("mtm psql error")
	}

	return out
}

func (s *SSH) WaitForCluster(nodeCount int) bool {
	query := "SELECT count(*) FROM mtm.nodes() WHERE enabled = 't' AND connected = 't'"

	tries, limit := 0, 60
	good := false

	for {
		fmt.Printf(".")

		out, _ := s.ExecOutputf("PGPASSWORD=1234 %s/psql -d mydb -U mtmuser -p 5432 -Atc \"%s\"", s.pgbin, query)
		count, _ := strconv.Atoi(out)
		if count == nodeCount {
			good = true
			break
		}

		if tries > limit {
			break
		}
		tries++

		time.Sleep(500 * time.Millisecond)
	}

	if good {
		fmt.Println("ok")
	} else {
		fmt.Println("fail")
	}

	return good
}

func (s *SSH) ExecInShell(toExec string) string {
	out, err := s.ExecOutputf("sh -c \"%s\"", toExec)
	if err != nil {
		log.WithField("conn", s.name).
			WithField("out", out).
			WithField("cmd", toExec).
			WithError(err).Error("executing cmd in shell failed")
	}

	return out
}
