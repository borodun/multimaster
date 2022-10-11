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
	Port   string
}

func (j *Joiner) Start(drop bool) {
	localIP := getLocalIP()

	if drop {
		log.Info("dropping node")
		j.stopPg()
		j.dropNode(localIP)
		return
	}

	connStr := j.addNode(localIP)

	log.RegisterExitHandler(func() {
		log.Info("something went wrong: dropping node")
		j.stopPg()
		j.dropNode(localIP)
	})

	j.stopPg()
	lsn := j.backupNodeAndGetLSN(connStr)
	j.configureBackup()
	j.startPg()

	log.Info("waiting for node to become ready")
	for i := 0; i < 10; i++ {
		_, err := execCmd(fmt.Sprintf("psql -U mtmuser -d mydb -p %s -c 'SELECT 1'", j.Port))
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
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
