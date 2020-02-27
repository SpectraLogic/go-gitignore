package gitignore

import (
	"strings"
)

// Separator fixed value "/"
var Separator = string("/")

type pattern struct {
	hasRootPrefix     bool
	hasDirSuffix      bool
	pathDepth         int
	matcher           pathMatcher
	onlyEqualizedPath bool
}

func newPattern(path string) (*pattern, error) {
	hasRootPrefix := path[0] == '/'
	hasDirSuffix := path[len(path)-1] == '/'

	var pathDepth int
	if !hasRootPrefix {
		pathDepth = strings.Count(path, "/")
	}

	var matcher pathMatcher
	matchingPath := strings.Trim(path, "/")
	if hasMeta(path) {
		globMatcher, err := newGlobMatcher(matchingPath)
		if err != nil {
			return nil, err
		}
		matcher = globMatcher
	} else {
		matcher = simpleMatcher{path: matchingPath}
	}

	return &pattern{
		hasRootPrefix: hasRootPrefix,
		hasDirSuffix:  hasDirSuffix,
		pathDepth:     pathDepth,
		matcher:       matcher,
	}, nil
}

func newPatternForEqualizedPath(path string) (*pattern, error) {
	pattern, err := newPattern(path)
	if err != nil {
		return nil, err
	}
	pattern.onlyEqualizedPath = true
	return pattern, nil
}

func (p pattern) match(path string, isDir bool) bool {

	if p.hasDirSuffix && !isDir {
		return false
	}

	var targetPath string
	if p.hasRootPrefix || p.onlyEqualizedPath {
		// absolute pattern or only equalized path mode
		targetPath = path
	} else {
		// relative pattern
		targetPath = p.equalizeDepth(path)
	}
	return p.matcher.match(targetPath)
}

func (p pattern) equalizeDepth(path string) string {
	equalizedPath, _ := cutLastN(path, p.pathDepth+1)
	return equalizedPath
}
