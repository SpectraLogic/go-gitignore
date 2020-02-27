package gitignore

import "strings"

type indexScanPatterns struct {
	absolute depthPatternHolder
	relative depthPatternHolder
}

func newIndexScanPatterns() *indexScanPatterns {
	return &indexScanPatterns{
		absolute: newDepthPatternHolder(asc),
		relative: newDepthPatternHolder(desc),
	}
}

func (ps *indexScanPatterns) add(pattern string) error {
	if strings.HasPrefix(pattern, "/") {
		return ps.absolute.add(pattern)
	} else {
		return ps.relative.add(pattern)
	}
}

func (ps indexScanPatterns) match(path string, isDir bool) bool {
	if ps.absolute.match(path, isDir) {
		return true
	}
	return ps.relative.match(path, isDir)
}

type scanStrategy interface {
	add(pattern string) error
	match(path string, isDir bool) bool
}
