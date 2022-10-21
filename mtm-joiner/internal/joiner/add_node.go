package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (j *Joiner) addNode() string {
	addNodeURL := fmt.Sprintf("%s/api/v1/add-node?host=%s&port=%s", j.URL, j.Addr, j.Port)

	log.Debugf("add url: %s", addNodeURL)

	resp, err := http.Get(addNodeURL)
	if err != nil {
		log.WithError(err).Fatalf("getting connection string: http get")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatalf("getting connection string: reading body")
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("response", string(body)).Fatal("getting connection string: status code not 200")
	}

	log.Info("node added")

	return string(body)
}
