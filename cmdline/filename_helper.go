package cmdline

import (
	"filesigner/filehelper"
	"path/filepath"
)

// ******** Public constants ********

// NegatePrefix is the prefix that begins a negation specification.
const NegatePrefix = '-'

// ******** Public functions ********

// GetAllFilePaths gets the command line file names and adds them to a file path array with all files matching a pattern.
func GetAllFilePaths(rawFilePaths []string) ([]string, error) {
	numFilePaths := len(rawFilePaths)
	result := make([]string, 0, numFilePaths)          // Result
	resultNames := make(map[string]bool, numFilePaths) // Check list to avoid duplicates
	negateList := make([]string, 0, numFilePaths)      // List of patterns which should *not* be contained in the result

	// Glob is always called, even when there is no pattern.
	for _, filePath := range rawFilePaths {
		if filePath[0] != NegatePrefix {
			matchingFiles, err := filehelper.SensibleGlob(filePath)
			if err != nil {
				return nil, err
			}

			for _, matchingFilePath := range matchingFiles {
				_, found := resultNames[matchingFilePath]
				if !found {
					result = append(result, filepath.Clean(matchingFilePath))
					resultNames[matchingFilePath] = true
				}
			}
		} else {
			negatePath := filePath[1:]
			if len(negatePath) != 0 {
				negateList = append(negateList, filePath[1:])
			}
		}
	}

	if len(negateList) > 0 {
		result = removeNegates(result, negateList)
	}

	return result, nil
}

// removeNegates removes all files matching the negate list.
func removeNegates(filePathList []string, negateList []string) []string {
	result := make([]string, 0, len(filePathList))

	for _, filePath := range filePathList {
		if !matchesPattern(filePath, negateList) {
			result = append(result, filePath)
		}
	}

	return result
}

// matchesPattern checks if a supplied file path matches any of the supplied file name patterns.
func matchesPattern(filePath string, patternList []string) bool {
	for _, pattern := range patternList {
		matched, _ := filepath.Match(pattern, filePath)
		if matched {
			return true
		}
	}

	return false
}
