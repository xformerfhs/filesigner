//go:build !windows

package filehelper

import (
	"os"
	"path/filepath"
)

// ******** Private functions ********

// sensibleGlobWithSwitch is the globbing function for all OSes except Windows.
func sensibleGlobWithSwitch(pattern string, withDirs bool, withFiles bool) ([]string, error) {
	// Remove trailing separators, if any
	pattern = ensureNoTrailingSeparator(pattern)

	// Find files and directories matching pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return matches, err
	}

	if !withDirs {
		matches = removeElements(matches, true)
	}
	if !withFiles {
		matches = removeElements(matches, false)
	}

	return matches
}

// removeElements removes all elements in the globbing list that are either no directories or no no files.
func removeElements(matchList []string, noDirs bool) ([]string, error) {
	result := make([]string, 0, len(matchList))
	for _, filePath := range matchList {
		fi, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if fi.IsDir() {
			if !noDirs {
				result = append(result, filePath)
			}
		} else {
			if noDirs {
				result = append(result, filePath)
			}
		}
	}

	return result, nil
}
