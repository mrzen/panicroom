package watcher

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

// Watcher is a device to watch a set of named files and pass their events to Alerters.
type Watcher struct {
	Globs     []string          `yaml:"paths"`
	FilePaths []string          `yaml:"-"`
	Watcher   *fsnotify.Watcher `yaml:"-"`
}

// New - Create a new watcher to monitor files matching the given include/exclude rules.
func New(globs, excludes []string) (*Watcher, error) {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	inc, err := expand(globs)

	if err != nil {
		return nil, err
	}

	exp, err := expand(excludes)

	if err != nil {
		return nil, err
	}

	paths, err := filter(inc, exp)

	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		err := watcher.Add(p)

		if err != nil {
			return nil, err
		}
	}

	return &Watcher{
		Globs:     globs,
		FilePaths: paths,
		Watcher:   watcher,
	}, nil

}

func expand(globs []string) ([]string, error) {
	var r []string

	for _, g := range globs {
		p, err := filepath.Glob(g)

		if err != nil {
			return nil, err
		}

		r = append(r, p...)
	}

	return r, nil
}


func filter(inc, exp []string) ([]string, error) {


	var paths []string

	for _, p := range inc {

		// Don't watch anything that's excluded explicitly
		for _, ex := range exp {
			if p == ex {
				continue
			}
		}

		// Don't watch directories directly
		// inofity will watch them recursively
		stat, err := os.Lstat(p)
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			paths = append(paths, p)
		}
	}

	return paths, nil
}