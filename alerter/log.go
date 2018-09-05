package alerter

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// LogAlerter - Handles alerts by logging them.
type LogAlerter struct {
	Path string `yaml:"path"`
	log *logrus.Logger
}

func (a *LogAlerter) SetConfig(config interface{}) error {
	var err error
	var f io.Writer

	m, ok := config.(map[string]interface{})

	if !ok {
		return errors.New("invalid config")
	}

	path, ok := m["path"].(string)

	if !ok || path == "" {
		return errors.New("missing required parameter: path")
	}

	if path == "-"  || path == "STDOUT" {
		f = os.Stdout
	} else {

		f, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			return err
		}
	}

	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{}
	l.Out = f

	a.log = l
	a.Path = path

	return err
}

func (a *LogAlerter) Alert(e Event) error {
	a.log.WithFields(logrus.Fields{
		"watcher": e.WatcherName,
		"operation": e.Operation,
		"ts": e.T,
		"path": e.Path,
	}).Warnln(e.Path, "was", e.Operation, "at", e.T)

	return nil
}