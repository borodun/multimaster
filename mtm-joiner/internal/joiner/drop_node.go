package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (j *Joiner) dropNode(ip string) {
	dropNodeURL := fmt.Sprintf("%s/api/v1/drop-node?host=%s", j.URL, ip)

	resp, err := http.Get(dropNodeURL)
	if err != nil {
		log.WithError(err).Fatalf("dropping node: http get")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatalf("dropping node: reading body")
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("response", body).Fatal("dropping node: status code not 200")
	}

	log.Infof("response: %s", body)
}
