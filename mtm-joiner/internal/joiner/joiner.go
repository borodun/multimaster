package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
	"time"
)

type Joiner struct {
	URL    string
	PGDATA string
	Port   string
	Addr   string
}

func (j *Joiner) Start(drop bool) {
	if j.Addr == "" {
		j.Addr = getLocalIP()
	}

	if drop {
		log.Info("dropping node")
		j.stopPg()
		j.dropNode()
		return
	}

	err := j.askToStopPg()
	if err != nil {
		log.WithError(err).Error("while stopping Postgres")
		return
	}

	err = j.removePGDATA()
	if err != nil {
		log.WithError(err).Error("while removing PGDATA")
		return
	}

	connStr := j.addNode()

	log.RegisterExitHandler(func() {
		log.Warn("something went wrong: dropping node")
		j.stopPg()
		j.dropNode()
	})

	j.stopPg()
	lsn := j.backupNodeAndGetLSN(connStr)
	j.configureBackup()
	j.startPg()

	log.Info("waiting for node to become ready for joining")
	for {
		fmt.Printf(".")
		_, err := execCmd(fmt.Sprintf("psql -U mtmuser -d mydb -p %s -c 'SELECT 1'", j.Port))
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n")
	log.Info("joining node")

	go func() {
		for {
			fmt.Printf(".")
			time.Sleep(1 * time.Second)
		}
	}()
	j.joinNode(lsn)
}

func (j *Joiner) askToStopPg() error {
	if j.pgIsReady() {
		var ans string

		fmt.Printf("Postgres is already running on port '%s'. Do you want to stop it? (y/n): ", j.Port)
		_, err := fmt.Scan(&ans)
		if err != nil {
			return err
		}

		ans = strings.TrimSpace(ans)
		ans = strings.ToLower(ans)

		if ans == "y" || ans == "yes" {
			j.stopPg()
		}

		return nil
	}
	return fmt.Errorf("user chose not to stop Postgres")
}

func (j *Joiner) removePGDATA() error {
	_, err := os.Stat(j.PGDATA)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	var ans string

	fmt.Printf("'%s' is not empty. Do you want to delete it? (y/n): ", j.PGDATA)
	_, err = fmt.Scan(&ans)
	if err != nil {
		return err
	}

	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		log.Infof("removing %s", j.PGDATA)
		_, err = execCmd(fmt.Sprintf("rm -rf %s", j.PGDATA))
		return err
	}
	return fmt.Errorf("user chose to keep %s", j.PGDATA)
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
