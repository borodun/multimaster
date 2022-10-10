package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (j *Joiner) backupNodeAndGetLSN(connStr string) string {
	cmdStr := fmt.Sprintf("pg_basebackup -D %s -d '%s' -c fast -v 2>&1", j.PGDATA, connStr)

	err := j.removePGDATA()
	if err != nil {
		log.WithError(err).Fatal("while removing PGDATA")
	}

	log.Info("starting backup process")
	out, err := execCmd(cmdStr)
	if err != nil {
		log.WithError(err).Fatal("cannot backup node")
	}
	log.Info(out)
	log.Info("backup completed")

	lsn := getLSN(out)
	if lsn == "" {
		log.WithField("out", out).Fatal("couldn't get lsn")
	}
	log.Infof("lsn: %s", lsn)

	return lsn
}

func (j *Joiner) removePGDATA() error {
	_, err := os.Stat(j.PGDATA)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	var ans string

	fmt.Printf("%s is not empty. Do you want to delete it? (y/n): ", j.PGDATA)
	_, err = fmt.Scan(&ans)
	if err != nil {
		return err
	}

	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		log.Infof("removing %s", j.PGDATA)
		_, err = execCmd(fmt.Sprintf("rm -rf %s", j.PGDATA))
		return err
	}
	return fmt.Errorf("user chose to keep %s", j.PGDATA)
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
