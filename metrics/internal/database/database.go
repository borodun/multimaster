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
	name    string
	db      *sqlx.DB
	timeout time.Duration
}

func NewDatabase(name string, db *sql.DB, timeout time.Duration) Database {
	var dbx = sqlx.NewDb(db, "postgres")
	return Database{
		name:    name,
		db:      dbx,
		timeout: timeout,
	}
}

func (d *Database) HasSharedPreloadLibrary(lib string) bool {
	var libs []string
	if err := d.Query("SHOW shared_preload_libraries", &libs); err != nil {
		return false
	}
	return strings.Contains(libs[0], lib)
}

func (d *Database) HasExtension(ext string) bool {
	var count int64
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(d.timeout),
	)
	defer func() {
		<-ctx.Done()
	}()
	if err := d.db.GetContext(
		ctx,
		&count,
		`
			SELECT count(*)
			FROM pg_available_extensions
			WHERE name = $1
			AND installed_version is not null
		`,
		ext,
	); err != nil {
		log.WithError(err).Errorf("failed to determine if %s is installed", ext)
	}
	cancel()
	return count > 0
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
			WithField("conn", d.name).
			WithField("query", q).
			WithField("params", params).
			Error("query failed")
	}
	cancel()
	return err
}

func (d *Database) Version() int {
	var version int
	if err := d.db.QueryRow("show server_version_num").Scan(&version); err != nil {
		log.WithField("db", d.name).WithError(err).Error("failed to get postgresql version, assuming 9.6.0")
		return 90600
	}
	return version
}

func (d *Database) Ping() error {
	return d.db.Ping()
}
