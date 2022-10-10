package joiner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (j *Joiner) addNode(ip string) string {
	addNodeURL := fmt.Sprintf("%s/api/v1/add-node?host=%s", j.URL, ip)

	return getConnStr(addNodeURL)
}

func getConnStr(url string) string {
	resp, err := http.Get(url)
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

	return string(body)
}
