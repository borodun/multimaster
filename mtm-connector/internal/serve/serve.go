package serve

import (
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
	"mtm-connector/internal/connector"
	"mtm-connector/internal/database"
	nurl "net/url"
)

func Run(port, url string) {
	url, err := pgpass.UpdateURL(url)
	if err != nil {
		log.WithError(err).Error("failed to add password from ~/.pgpass")
	}
	if !checkPass(url) {
		log.Error("no password provided for user")
	}
	connInfo, _ := pq.ParseURL(url)

	mtmConnector := connector.MtmConnector{
		Db:       database.NewDatabase(connInfo),
		ConnInfo: connInfo,
		Hosts:    map[string]string{},
		Joined:   map[string]bool{},
	}

	mtmConnector.Serve(port)
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
