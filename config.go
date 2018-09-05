package main

import (
	"errors"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"os"
)

// Configuration - Top level configuration data
type Configuration struct {
	Watchers []WatcherConfig `yaml:"watchers"`
	Alerters []AlerterConfig `yaml:"alerters"`
}

// WatcherConfig - Configuration for a single watcher
type WatcherConfig struct {
	Name string `yaml:"name"`
	Paths []string `yaml:"paths"`
	Excludes []string `yaml:"excludes"`
	Alerters []string `yaml:"alerters"`
}

// AlerterConfig - Configuration for a single alerter.
// Alerter type-specific params are kept in an opaque map for processing by the alerter initializer.
type AlerterConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Configuration map[string]interface{} `yaml:"config"`
}

func readConfig() (*Configuration, error) {
		path := pflag.String("config-path", "./panicroom.yaml", "Configuration file path")
		pflag.Parse()

		if *path == "" {
			return nil, errors.New("no configuration file specified")
		}

		f, err := os.OpenFile(*path, os.O_RDONLY, 0644)

		if err != nil {
			return nil, err
		}

		dec := yaml.NewDecoder(f)

		config := &Configuration{}

		err = dec.Decode(config)

		return config, err
}