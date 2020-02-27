package gitignore

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type IgnoreMatcher interface {
	Match(path string, isDir bool) bool
}

type gitIgnore struct {
	ignorePatterns scanStrategy
	acceptPatterns scanStrategy
	path           string
}

func NewGitIgnore(gitignore string, base ...string) (IgnoreMatcher, error) {
	var path string
	if len(base) > 0 {
		path = base[0]
	} else {
		path = filepath.Dir(gitignore)
	}

	file, err := os.Open(gitignore)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return NewGitIgnoreFromReader(path, file)
}

func NewGitIgnoreFromReader(path string, r io.Reader) (*gitIgnore, error) {
	g := gitIgnore{
		ignorePatterns: newIndexScanPatterns(),
		acceptPatterns: newIndexScanPatterns(),
		path:           path,
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		line = fixPath(line)
		if strings.HasPrefix(line, "!") {
			err := g.acceptPatterns.add(fixRootPrefix(strings.TrimPrefix(line, "!")))
			if err != nil {
				return nil, err
			}
		} else {
			err := g.ignorePatterns.add(fixRootPrefix(line))
			if err != nil {
				return nil, err
			}
		}
	}
	return &g, nil
}

func NewGitIgnoreFromStrings(path string, patterns []string) (*gitIgnore, error) {
	g := gitIgnore{
		ignorePatterns: newIndexScanPatterns(),
		acceptPatterns: newIndexScanPatterns(),
		path:           path,
	}
	for _, line := range patterns {
		line := strings.Trim(line, " ")
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		line = fixPath(line)
		if strings.HasPrefix(line, "!") {
			err := g.acceptPatterns.add(fixRootPrefix(strings.TrimPrefix(line, "!")))
			if err != nil {
				return nil, err
			}
		} else {
			err := g.ignorePatterns.add(fixRootPrefix(line))
			if err != nil {
				return nil, err
			}
		}
	}
	return &g, nil
}

func (g gitIgnore) Match(path string, isDir bool) bool {
	relativePath, err := filepath.Rel(g.path, path)
	_ = "breakpoint"
	if err != nil {
		return false
	}

	relativePath = fixPath(relativePath)
	pbuild := ""
	pdir := true
	for _, p := range strings.Split(relativePath, "/") {
		// Get path up to a specific depth
		if pbuild == "" {
			pbuild += p
		} else {
			pbuild += "/" + p
		}

		// Everything except for the last element is a directory
		if pbuild == relativePath {
			pdir = isDir
		}

		// Try accepted patterns (which never match)
		if g.acceptPatterns.match(pbuild, pdir) {
			return false
		}
		// Try ignored patterns (which match if there's no matching accepted pattern)
		if g.ignorePatterns.match(pbuild, pdir) {
			return true
		}
	}
	return false
}
