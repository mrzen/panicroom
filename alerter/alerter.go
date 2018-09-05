package alerter

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"time"
)

type Event struct {
	WatcherName string `json:"watcher"`
	Path string `json:"path"`
	T time.Time `json:"ts"`
	Operation fsnotify.Op `json:"op"`
}

// An alerter is a channel which can handle alerts/events from a set of watchers.
type Alerter interface {
	// Alert will handle the alert for a single event
	// TODO: Refactor as a chan.
	Alert(Event) error

	// Set the configuration
	SetConfig(interface{}) error
}


var errUnrecognizedType = errors.New("unrecognized type")

// New - Create a new alerter with the given type name and configuration
func New(typeName string, config interface{}) (Alerter, error) {

	l, err := getAlerter(typeName)

	if err != nil {
		return nil, err
	}

	err = l.SetConfig(config)

	return l, err
}


func getAlerter(typeName string) (Alerter, error) {
	switch typeName {
	case "log":
		return new(LogAlerter), nil
	case "sns":
		return new(SNSAlerter), nil
	default:
		return nil, errUnrecognizedType
	}
}