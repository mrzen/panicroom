package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"panicroom/alerter"
	"panicroom/watcher"
	"time"
)

func main() {

	config, err := readConfig()

	if err != nil {
		logrus.Fatal("Unable to load config:", err)
	}

	alerters := make(map[string]alerter.Alerter)

	for _, a := range config.Alerters {

		al, err := alerter.New(a.Type, a.Configuration)

		if err != nil {
			logrus.Fatal("Unable to initialize alerter", a.Name, ":", err)
		}

		alerters[a.Name] = al
	}

	nFiles := uint64(0)

	for _, wc := range config.Watchers {
		w, err := watcher.New(wc.Paths, wc.Excludes)

		if err != nil {
			logrus.Fatal("Unable to set up watcher", wc.Name, ":", err)
		}

		nFiles += uint64(len(w.FilePaths))

		// TODO: Refactor this to fan-out watcher events chan into Alerter events chans.
		go func() {
			for {
				e := <-w.Watcher.Events

				if e.Op == fsnotify.Create {
					nFiles++
				}
				event := alerter.Event{
					WatcherName: wc.Name,
					Path:        e.Name,
					Operation:   e.Op,
					T:           time.Now(),
				}

				for _, an := range wc.Alerters {
					al := alerters[an]

					if al == nil {
						logrus.Errorln("Alerter ", an, " does not exist. Cannot alert")
						continue
					}

					if err = al.Alert(event); err != nil {
						logrus.WithError(err).Errorln("Failed to push alert to", an, ":", err)
					}
				}
			}
		}()
	}

	logrus.Infoln("Started watching ", nFiles, "files for changes.")

	for {
		time.Sleep(time.Second)
	}
}