package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Database struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewDatabase(connInfo string) Database {
	log.Infof(connInfo)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.WithError(err).Warn("failed to open url")
	}

	if err = db.Ping(); err != nil {
		log.WithError(err).Warn("failed to ping database")
	}

	var dbx = sqlx.NewDb(db, "postgres")
	return Database{
		db:      dbx,
		timeout: 600 * time.Second,
	}
}

var EmptyParams []interface{}

func (d *Database) Query(
	query string,
	result interface{},
) error {
	return d.queryWithTimeout(query, result, EmptyParams, d.timeout)
}

func (d *Database) QueryWithParams(
	query string,
	result interface{},
	params []interface{},
) error {
	return d.queryWithTimeout(query, result, params, d.timeout)
}

func (d *Database) queryWithTimeout(
	query string,
	result interface{},
	params []interface{},
	timeout time.Duration,
) error {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(timeout),
	)
	defer func() {
		<-ctx.Done()
	}()
	var err = d.db.SelectContext(ctx, result, query, params...)
	if err != nil {
		var q = strings.Join(strings.Fields(query), " ")
		if len(q) > 50 {
			q = q[:50] + "..."
		}
		log.WithError(err).
			WithField("query", q).
			WithField("params", params).
			Error("query failed")
	}
	cancel()
	return err
}
