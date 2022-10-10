package joiner

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func execCmd(args ...string) (string, error) {
	ar := append([]string{"-c"}, args...)
	cmd := exec.Command("sh", ar...)

	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	err := cmd.Run()

	if err != nil {
		log.WithField("stdout", o.String()).
			WithField("stderr", e.String()).
			WithField("cmd", cmd.Args).
			WithError(err).Warn("exec error")

		return "", err
	}

	return o.String(), nil
}
