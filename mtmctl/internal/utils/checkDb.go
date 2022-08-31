package utils

import (
	"backup/internal/connection"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func CheckDatabase(db *connection.DB) error {
	if err := db.TruePing(); err != nil {
		return fmt.Errorf("cannot ping '%s' database: %s", db.GetName(), err.Error())
	}
	if stat := db.MtmStatus(); stat != "online" {
		return fmt.Errorf("'%s' node is not online, current status: %s", db.GetName(), stat)
	}
	return nil
}

func CheckDatabases(dbs ...*connection.DB) {
	for _, db := range dbs {
		if err := CheckDatabase(db); err != nil {
			println(err.Error())
			log.WithError(err).Fatal("bad connection to db")
		}
	}
}
