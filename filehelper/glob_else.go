//go:build !windows

package filehelper

import (
	"os"
	"path/filepath"
)

// ******** Public functions ********

// SensibleGlob is the globbing function for all OSes except Windows.
func SensibleGlob(pattern string) ([]string, error) {
	// Remove trailing separators, if any
	pattern = ensureNoTrailingSeparator(pattern)

	// Find files matching pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return matches, err
	}

	return removeDirs(matches)
}

// ******** Private functions ********

// removeDirs removes all directories in the globbing list.
func removeDirs(matchList []string) ([]string, error) {
	result := make([]string, 0, len(matchList))
	for _, filePath := range matchList {
		fi, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			result = append(result, filePath)
		}
	}

	return result, nil
}
