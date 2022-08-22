package add

import (
	"backup/internal/config"
	"backup/internal/connection"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type MtmAddNode struct {
	Cfg         config.Config
	Connections connection.Connections
	SrcNode     string
	DstNode     string

	srcDB  *connection.DB
	dstSSH *connection.SSH

	ConnStr string
}

func (m *MtmAddNode) Run() {
	m.srcDB = m.Connections[m.SrcNode].DB
	m.dstSSH = m.Connections[m.DstNode].SSH

	fmt.Println("Checking connection")
	if err := m.srcDB.TruePing(); err != nil {
		fmt.Printf("Cannot ping '%s' database", m.SrcNode)
		log.WithError(err).Fatalf("cannot ping '%s' database", m.SrcNode)
	}
	if stat := m.srcDB.MtmStatus(); stat != "online" {
		fmt.Printf("Multimaster node is not online, current status: %s\n", stat)
		log.Fatalf("multimaster node is not online, current status: %s", stat)
	}

	fmt.Printf("Creating replication slot for new node \n")
	id := m.mtmAddNodeAndGetID()
	fmt.Printf("Created replication slot, node id: %s\n", id)

	log.RegisterExitHandler(func() {
		m.mtmDropNodeByID(id)
		fmt.Printf("Something went wrong, add --verbose flag to see logs")
	})

	m.dstSSH.PgCtlStop()
	fmt.Printf("Postgres stopped\n")
	fmt.Printf("Getting backup from '%s'\n", m.SrcNode)
	lsn := m.backupNodeAndGetLSN()
	fmt.Printf("Backup finished successfully, lsn: %s\n", lsn)

	m.configureBackup()
	fmt.Printf("Db '%s' configured for recovery\n", m.DstNode)
	m.dstSSH.PgCtlStart()
	fmt.Printf("Postgres started\n")
	fmt.Printf("Wait untill it will be joined by multimaster cluster\n")

	time.Sleep(time.Duration(10) * time.Second)

	m.mtmJoinNode(id, lsn)
	fmt.Printf("Successfully joined '%s' to cluster\n", m.DstNode)
}

func (m *MtmAddNode) mtmJoinNode(id, lsn string) {
	mtmJoinNodeQuery := fmt.Sprintf(`SELECT mtm.join_node(%s, '%s')`, id, lsn)

	var info []string
	err := m.srcDB.Query(mtmJoinNodeQuery, &info)

	if err != nil {
		log.WithError(err).Fatal("cannot join node")
	}
	log.Infof("node added successfully, id: %s", id)
}

func (m *MtmAddNode) setNodeConnStr() {
	if m.ConnStr != "" {
		return
	}

	conn := m.Cfg.GetConn(m.DstNode)
	if conn.URL == "" {
		log.WithField("conn", conn.Name).
			Fatal("you need to specify URL in config or connection string in cmd")
	}

	host := conn.Ssh.Host
	user, err := conn.GetFieldFromConnStr("user")
	if err != nil {
		log.WithField("conn", conn.Name).WithError(err).Fatal()
	}
	pass, err := conn.GetFieldFromConnStr("password")
	if err != nil {
		log.WithField("conn", conn.Name).WithError(err).Fatal()
	}
	dbname, err := conn.GetFieldFromConnStr("dbname")
	if err != nil {
		log.WithField("conn", conn.Name).WithError(err).Fatal()
	}

	m.ConnStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, pass, dbname)
}

func (m *MtmAddNode) mtmDropNodeByID(id string) {
	mtmAddNodeQuery := fmt.Sprintf(`SELECT mtm.drop_node(%s)`, id)

	db := m.Connections[m.SrcNode].DB

	var info []string
	err := db.Query(mtmAddNodeQuery, &info)

	if err != nil {
		log.WithError(err).Warn("cannot drop node")
	}
}

func (m *MtmAddNode) mtmAddNodeAndGetID() string {
	m.setNodeConnStr()
	mtmAddNodeQuery := fmt.Sprintf(`SELECT mtm.add_node('%s')`, m.ConnStr)

	db := m.Connections[m.SrcNode].DB

	var id []string
	err := db.Query(mtmAddNodeQuery, &id)

	if err != nil {
		log.WithError(err).Fatal("cannot add node")
	}

	retId := id[0]
	log.Infof("id: %s", retId)
	return retId
}

func (m *MtmAddNode) backupNodeAndGetLSN() string {
	cmd, err := m.getPgBaseBackupCmd()
	if err != nil {
		log.WithError(err).Fatal("cannot create pg_basebackup cmd")
	}

	m.dstSSH.RemovePGDATA()

	log.Info("starting backup process")
	out, err := m.dstSSH.ExecOutput(cmd)
	if err != nil {
		log.WithError(err).
			WithField("out", out).
			WithField("cmd", cmd).
			Fatal("cannot exec pg_basebackup cmd")
	}
	log.Info(out)
	log.Info("backup completed")

	lsn := m.getLSN(out)
	if lsn == "" {
		log.Fatal("couldn't get lsn")
	}
	log.Infof("lsn: %s", lsn)

	return lsn
}

func (m *MtmAddNode) configureBackup() {
	m.dstSSH.Execf("touch %s/recovery.signal", m.Cfg.GetPGDATA())

	settings := []string{
		`"restore_command = 'false'"`,
		`"recovery_target = 'immediate'"`,
		`"recovery_target_action = 'promote'"`,
	}

	for _, setting := range settings {
		m.dstSSH.Execf("echo %s >> %s/postgresql.conf", setting, m.Cfg.GetPGDATA())
	}

	log.Info("postgres configured")
}

func (m *MtmAddNode) getLSN(out string) string {
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

func (m *MtmAddNode) getPgBaseBackupCmd() (string, error) {
	conn := m.Cfg.GetConn(m.SrcNode)
	if conn == nil {
		log.Fatalf("no connection specified for %s", m.SrcNode)
	}

	connStr := conn.ParseURL()

	cmd := fmt.Sprintf("%s/pg_basebackup -D %s -d '%s' -c fast -v 2>&1", m.Cfg.GetPGBIN(), m.Cfg.GetPGDATA(), connStr)
	return cmd, nil
}
