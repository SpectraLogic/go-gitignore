package gitignore

import "strings"

const (
	asc = iota
	desc
)

type depthPatternHolder struct {
	patterns depthPatterns
	order    int
}

func newDepthPatternHolder(order int) depthPatternHolder {
	return depthPatternHolder{
		patterns: depthPatterns{m: map[int]initialPatternHolder{}},
		order:    order,
	}
}

func (h *depthPatternHolder) add(pattern string) error {
	count := strings.Count(strings.Trim(pattern, "/"), "/")
	return h.patterns.set(count+1, pattern)
}

func (h depthPatternHolder) match(path string, isDir bool) bool {
	if h.patterns.size() == 0 {
		return false
	}

	for depth := 1; ; depth++ {
		var part string
		var isLast, isDirCurrent bool
		if h.order == asc {
			part, isLast = cutN(path, depth)
			if isLast {
				isDirCurrent = isDir
			} else {
				isDirCurrent = false
			}
		} else {
			part, isLast = cutLastN(path, depth)
			isDirCurrent = isDir
		}
		if patterns, ok := h.patterns.get(depth); ok {
			if patterns.match(part, isDirCurrent) {
				return true
			}
		}
		if isLast {
			break
		}
	}
	return false
}

type depthPatterns struct {
	m map[int]initialPatternHolder
}

func (p *depthPatterns) set(depth int, pattern string) error {
	if ps, ok := p.m[depth]; ok {
		return ps.add(pattern)
	} else {
		holder := newInitialPatternHolder()
		err := holder.add(pattern)
		if err != nil {
			return err
		}
		p.m[depth] = holder
	}
	return nil
}

func (p depthPatterns) get(depth int) (initialPatternHolder, bool) {
	patterns, ok := p.m[depth]
	return patterns, ok
}

func (p depthPatterns) size() int {
	return len(p.m)
}
