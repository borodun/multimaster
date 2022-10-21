package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (j *Joiner) backupNodeAndGetLSN(connStr string) string {
	cmdStr := fmt.Sprintf("pg_basebackup -D %s -d '%s' -c fast -v 2>&1", j.PGDATA, connStr)

	log.Info("starting backup process")
	out, err := execCmd(cmdStr)
	if err != nil {
		log.WithError(err).Fatal("cannot backup node")
	}
	log.Debug(out)
	log.Info("backup completed")

	lsn := getLSN(out)
	if lsn == "" {
		log.WithField("out", out).Fatal("couldn't get lsn")
	}
	log.Debug("lsn: %s", lsn)

	return lsn
}

func getLSN(out string) string {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "write-ahead log end point:") {
			lsn := strings.Split(line, ":")[2]
			lsn = strings.Trim(lsn, " ")
			return lsn
		}
	}
	return ""
}
