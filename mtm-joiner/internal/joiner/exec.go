package joiner

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func execCmd(cmdStr string, args ...string) (string, error) {
	cmd := exec.Command(cmdStr, args...)

	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	err := cmd.Run()

	if err != nil {
		log.WithField("stdout", o.String()).
			WithField("stderr", e.String()).
			WithField("cmd", cmdStr).
			WithField("args", args).
			WithError(err).Error("exec error")

		return "", err
	}

	return o.String(), nil
}
