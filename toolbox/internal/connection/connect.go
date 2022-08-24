package connection

import (
	"backup/internal/config"
	"database/sql"
	"fmt"
	"github.com/k0sproject/rig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
	nurl "net/url"
)

type Conf struct {
	ConnName    string
	ConnectDb   bool
	DbRequired  bool
	ConnectSsh  bool
	SshRequired bool
}

func Connect(cfg config.Config, connConf []Conf) Connections {
	logger := log.New()
	logger.SetLevel(log.GetLevel())
	rig.SetLogger(logger)

	connNames := cfg.GetAllConnNames()
	for _, conf := range connConf {
		if !arrayContains(connNames, conf.ConnName) {
			log.Fatalf("config doesn't contain connection for '%s'", conf.ConnName)
		}
	}

	connections := make(Connections)

	for _, conf := range connConf {
		conn := cfg.GetConn(conf.ConnName)

		connLog := log.WithField("conn", conn.Name)

		ssh, err := connectToHost(conn.Ssh, conf.ConnectSsh)
		if err != nil {
			if conf.SshRequired {
				connLog.WithError(err).Fatal("cannot connect to host, but it is required")
			}
			connLog.WithError(err).Warn("cannot connect to host, skipping")
		} else if conf.ConnectSsh {
			ssh.name = conn.Name
			ssh.pgdata = cfg.GetPGDATA()
			ssh.pgbin = cfg.GetPGBIN()
		}

		db, err := connectToPostgres(conn.URL, conf.ConnectDb)
		if err != nil {
			if conf.DbRequired {
				connLog.WithError(err).Fatal("cannot connect to database, but it is required")
			}
			connLog.WithError(err).Warn("cannot connect to database, skipping")
		}

		connections[conn.Name] = Connection{
			SSH: ssh,
			DB:  NewDB(db, conn.Name),
		}

		connLog.Info("connected")
	}

	return connections
}

func arrayContains(arr []string, name string) bool {
	for _, e := range arr {
		if e == name {
			return true
		}
	}
	return false
}

func connectToHost(sshConf config.Ssh, required bool) (*SSH, error) {
	if !required {
		return nil, nil
	}

	ssh := rig.SSH{
		User:    sshConf.User,
		Address: sshConf.Host,
		Port:    sshConf.Port,
		KeyPath: sshConf.Key,
		Bastion: getBastion(sshConf.Bastion),
	}
	ssh.SetDefaults()

	h := rig.Connection{SSH: &ssh}
	if err := h.Connect(); err != nil {
		return nil, err
	}
	s := &SSH{Connection: h}
	return s, nil
}

func getBastion(bastion config.Bastion) *rig.SSH {
	if bastion.Host == "" {
		return nil
	}

	b := &rig.SSH{
		User:    bastion.User,
		Address: bastion.Host,
		Port:    bastion.Port,
		KeyPath: bastion.Key,
	}
	b.SetDefaults()

	return b
}

func connectToPostgres(url string, required bool) (*sql.DB, error) {
	if !required {
		return nil, nil
	}

	url, err := pgpass.UpdateURL(url)
	if err != nil {
		return nil, err
	}
	if !checkPass(url) {
		return nil, fmt.Errorf("no password provided for user")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)

	return db, nil
}

func checkPass(url string) bool {
	u, err := nurl.Parse(url)
	if err != nil {
		return false
	}
	if user := u.User; user != nil {
		if _, ok := user.Password(); !ok {
			return false
		}
	}
	return true
}
