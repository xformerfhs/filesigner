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

	return filterMatches(matches, withDirs, withFiles)
}

// filterMatches filters all elements in the globbing list according to being directories or files.
func filterMatches(matchList []string, withDirs bool, withFiles bool) ([]string, error) {
	result := make([]string, 0, len(matchList))
	for _, fileName := range matchList {
		fi, err := os.Stat(fileName)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			if withFiles {
				result = append(result, fileName)
			}
		} else {
			if withDirs {
				if fileName != `.` && fileName != `..` {
					result = append(result, fileName)
				}
			}
		}
	}

	return result, nil
}
