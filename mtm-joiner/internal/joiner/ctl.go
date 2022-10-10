package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (j *Joiner) pgIsReady() bool {
	out, err := execCmd(fmt.Sprintf("pg_isready -p %s", j.Port))
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		log.WithError(err).Fatal("pg_isready error")
	}

	return strings.Contains(out, "accepting connections")
}

func (j *Joiner) pgCtl(cmd string) {
	c := fmt.Sprintf("pg_ctl -D %s -l %s/logfile -o '-p %s' %s", j.PGDATA, j.PGDATA, j.Port, cmd)
	_, err := execCmd(c)
	if err != nil {
		log.WithError(err).
			WithField("cmd", cmd).
			Fatal("pg_ctl error")
	}
}

func (j *Joiner) startPg() {
	if j.pgIsReady() {
		return
	}
	j.pgCtl("start")
}

func (j *Joiner) stopPg() {
	if !j.pgIsReady() {
		return
	}
	j.pgCtl("stop")
}
