package gitignore

import (
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
)

type pathMatcher interface {
	match(path string) bool
}

type simpleMatcher struct {
	path string
}

func (m simpleMatcher) match(path string) bool {
	return m.path == path
}

type filepathMatcher struct {
	path string
}

func (m filepathMatcher) match(path string) bool {
	match, _ := filepath.Match(m.path, path)
	return match
}

type globMatcher struct {
	path string
	glob glob.Glob
}

func newGlobMatcher(path string) (*globMatcher, error) {
	matcher, err := glob.Compile(path, os.PathSeparator)
	if err != nil {
		return nil, err
	}
	g := globMatcher{
		path: path,
		glob: matcher,
	}
	return &g, nil
}

func (m globMatcher) match(path string) bool {
	match := m.glob.Match(path)
	return match
}
