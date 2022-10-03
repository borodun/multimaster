package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (j *Joiner) configureBackup() {
	_, err := execCmd(fmt.Sprintf("touch %s/recovery.signal", j.PGDATA))
	if err != nil {
		log.WithError(err).Fatalf("cannot create %s/recovery.signal", j.PGDATA)
	}

	settings := []string{
		`"restore_command = 'false'"`,
		`"recovery_target = 'immediate'"`,
		`"recovery_target_action = 'promote'"`,
	}

	for _, setting := range settings {
		_, err = execCmd(fmt.Sprintf("echo %s >> %s/postgresql.conf", setting, j.PGDATA))
		if err != nil {
			log.WithError(err).Fatalf("insert %s setting into %s/postgresql.conf", setting, j.PGDATA)
		}
	}

	log.Info("postgres configured")
}
