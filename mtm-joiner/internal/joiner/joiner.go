package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Joiner struct {
	URL    string
	PGDATA string
}

func (j *Joiner) Start() {
	localIP := getLocalIP()

	connStr := j.addNode(localIP)

	log.RegisterExitHandler(func() {
		fmt.Println("Something went wrong: dropping node")
		j.dropNode(localIP)
	})

	j.stopPg()
	lsn := j.backupNodeAndGetLSN(connStr)
	j.configureBackup()
	j.startPg()

	for !j.pgIsReady() {
		log.Infof("postgres isn't ready yet: waiting")
		time.Sleep(100 * time.Millisecond)
	}

	j.joinNode(localIP, lsn)
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.WithError(err).Fatal("while getting local IP")
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String()
}
