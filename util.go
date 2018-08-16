package gitignore

import (
	"fmt"
	"runtime"
	"os"
	"strings"
)

func cutN(path string, n int) (string, bool) {
	isLast := true

	var i, count int
	for i < len(path)-1 {
		if os.IsPathSeparator(path[i]) {
			count++
			if count >= n {
				isLast = false
				break
			}
		}
		i++
	}
	return path[:i+1], isLast
}

func cutLastN(path string, n int) (string, bool) {
	isLast := true
	i := len(path) - 1

	var count int
	for i >= 0 {
		if os.IsPathSeparator(path[i]) {
			count++
			if count >= n {
				isLast = false
				break
			}
		}
		i--
	}
	return path[i+1:], isLast
}

func hasMeta(path string) bool {
	return strings.IndexAny(path, "*?[") >= 0
}

const isWindows bool = runtime.GOOS == "windows"

func fixRootPrefix(path string) string {
	if (isWindows) {
		if (len(path) >=3) {
			prefix := path[1:3]
			hasWindowsRootPrefix := prefix == ":/" || prefix == ":\\"
			if hasWindowsRootPrefix {
				driveLetter := path[0:1]
				pathTail := path[3:]
				return fmt.Sprintf("/%v/%v", driveLetter, pathTail)
			}
		}

		return path
	}

	return path
}

func fixPath(path string) string {
	if (isWindows) {
		return strings.Replace(path, "\\", "/", -1)
	}

	return path
}
