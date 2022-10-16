package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (j *Joiner) joinNode(lsn string) {
	joinNodeURL := fmt.Sprintf("%s/api/v1/join-node?host=%s&port=%s&lsn=%s", j.URL, j.Addr, j.Port, lsn)

	log.Infof("join url: %s", joinNodeURL)

	resp, err := http.Get(joinNodeURL)
	if err != nil {
		log.WithError(err).Fatalf("joining node: http get")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatalf("joining node: reading body")
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("response", string(body)).Fatal("joining node: status code not 200")
	}

	log.Infof("response: %s", string(body))
}
