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
		log.WithError(err).Error("dropping node: http get")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("dropping node: reading body")
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("response", string(body)).Error("dropping node: status code not 200")
	}

	log.Infof("response: %s", string(body))
}
