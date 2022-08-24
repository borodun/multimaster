package connection

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DB struct {
	db      *sqlx.DB
	name    string
	timeout time.Duration
}

func NewDB(db *sql.DB, name string) *DB {
	if db == nil {
		return nil
	}

	var dbx = sqlx.NewDb(db, "postgres")
	d := &DB{
		db:      dbx,
		name:    name,
		timeout: time.Duration(600) * time.Second,
	}
	return d
}

func (d *DB) GetName() string {
	return d.name
}

func (d *DB) HasSharedPreloadLibrary(lib string) bool {
	var libs []string
	if err := d.Query("SHOW shared_preload_libraries", &libs); err != nil {
		return false
	}
	return strings.Contains(libs[0], lib)
}

func (d *DB) HasExtension(ext string) bool {
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

var emptyParams []interface{}

func (d *DB) Query(
	query string,
	result interface{},
) error {
	return d.queryWithTimeout(query, result, emptyParams, d.timeout)
}

func (d *DB) QueryP(
	query string,
	result interface{},
	params []interface{},
) error {
	return d.queryWithTimeout(query, result, params, d.timeout)
}

func (d *DB) queryWithTimeout(
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

func (d *DB) Ping() bool {
	var ret []string

	err := d.Query(`SELECT 1`, &ret)
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		return false
	}
	return true
}

func (d *DB) TruePing() error {
	return d.db.Ping()
}

type statusTup struct {
	Id     string `db:"my_node_id"`
	Status string `db:"status"`
}

const mtmNodeStatusQuery = `
		SELECT my_node_id, status FROM mtm.status()
	`

func (d *DB) MtmStatus() string {

	var statusTypes = []string{
		"online",
		"recovery",
		"catchup",
		"disabled",
		"isolated",
	}

	var status []statusTup
	err := d.Query(mtmNodeStatusQuery, &status)
	if err != nil {
		if strings.Contains(err.Error(), "multimaster node is not online: current status") {
			return strings.Split(err.Error(), ":")[2]
		} else if strings.Contains(err.Error(), "multimaster is not configured") {
			return "multimaster is not configured"
		}
		return fmt.Sprintf("unknown err: %s", err.Error())
	}

	for _, statusType := range statusTypes {
		if status[0].Status == statusType {
			return statusType
		}
	}

	return ""
}

func (d *DB) GetMtmNodeID() string {
	var status []statusTup
	err := d.Query(mtmNodeStatusQuery, &status)
	if err != nil {
		log.WithError(err).
			WithField("conn", d.name).Fatal("couldn't get node id in mtm")
	}
	return status[0].Id
}
