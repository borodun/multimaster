package connection

import (
	"database/sql"
	"fmt"
	"mtmctl/internal/config"
	nurl "net/url"
	"strings"
	"sync"

	"github.com/k0sproject/rig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
)

type Conf struct {
	ConnName    string
	ConnectDb   bool
	DbRequired  bool
	ConnectSsh  bool
	SshRequired bool
}

var confMapMutex = sync.Mutex{}

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

	fmt.Printf("Connecting to nodes: %s\n", strings.Join(connNames, ", "))

	connections := make(Connections)

	var wg sync.WaitGroup

	for _, conf := range connConf {
		wg.Add(1)
		go connectNode(conf, cfg, connections, &wg)
	}

	wg.Wait()

	return connections
}

func connectNode(conf Conf, cfg config.Config, connections Connections, wg *sync.WaitGroup) {
	defer wg.Done()

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

	confMapMutex.Lock()
	connections[conn.Name] = Connection{
		SSH: ssh,
		DB:  NewDB(db, conn.Name),
	}
	confMapMutex.Unlock()

	connLog.Info("connected")
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
