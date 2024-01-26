package filehelper

import (
	"os"
	"strings"
)

// ******** Private constants ********

var pathSeparatorString = string(os.PathSeparator)

// ******** Public functions ********

// PathGlobs returns all globbed files from an arbitrary absolute path that may contain wildcards.
func PathGlobs(arbitraryAbsPath string) ([]string, error) {
	parts := strings.Split(arbitraryAbsPath, pathSeparatorString)
	lastPartIndex := len(parts) - 1
	var err error

	fullPaths := make([]string, 1)
	fullPaths[0] = ""

	for i, aPart := range parts {
		fullPaths, err = walkThroughPart(fullPaths, aPart, i != lastPartIndex)
		if err != nil {
			return nil, err
		}
	}

	return fullPaths, nil
}

// ******** Private functions ********

// walkThroughPart loops through each element of full path, adds part to it and returns all file names
// that match the resulting specification. If findDirs is true, it will search for directories,
// otherwise it will search for files.
func walkThroughPart(fullPath []string, part string, findDirs bool) ([]string, error) {
	result := make([]string, 0, len(fullPath)<<1)

	for _, aFullPath := range fullPath {
		fullPathLen := len(aFullPath)

		var testPath string

		if fullPathLen != 0 {
			testPath = aFullPath + pathSeparatorString + part
		} else {
			testPath = part
		}

		var matches []string
		var err error
		if findDirs {
			matches, err = SensibleGlobDirsOnly(testPath)
		} else {
			matches, err = SensibleGlobFilesOnly(testPath)
		}

		if err != nil {
			return nil, err
		}

		for _, aMatch := range matches {
			if fullPathLen != 0 {
				result = append(result, aFullPath+pathSeparatorString+aMatch)
			} else {
				result = append(result, aMatch)
			}
		}
	}

	return result, nil
}
