package filehelper

import (
	"filesigner/slicehelper"
	"os"
	"path/filepath"
)

// ******** Public functions ********

// SplitPath takes a path string and splits it into its individual parts.
func SplitPath(path string) (parts []string) {
	vol := filepath.VolumeName(path)
	volLen := len(vol)
	if volLen > 0 {
		path = path[volLen:]
	}

	var dir string
	var file string
	var dirLen int

	dir = path
	for {
		dir, file = filepath.Split(dir)
		if len(file) == 0 {
			break
		}
		dirLen = len(dir)

		parts = slicehelper.Prepend(file, parts)
		lastDirIndex := dirLen - 1
		if dirLen > 1 && os.IsPathSeparator(dir[lastDirIndex]) {
			dir = dir[:lastDirIndex]
		}
	}

	if dirLen != 0 {
		parts = slicehelper.Prepend(dir, parts)
	}

	if volLen != 0 {
		parts[0] = vol + parts[0]
	}

	return
}

// PathGlob returns all globbed files from a path that may contain wildcards.
func PathGlob(path string) ([]string, error) {
	// An empty path returns nil
	if len(path) == 0 {
		return nil, nil
	}

	// Remove unnecessary "." and "..".
	// Clean converts "" to "." which is wrong for this function. This is why this case is handled
	// before the call to Clean.
	path = filepath.Clean(path)

	parts := SplitPath(path)
	lastPartIndex := len(parts) - 1
	var err error

	fullPaths := make([]string, 1)
	fullPaths[0] = ``

	for i, aPart := range parts {
		if os.IsPathSeparator(aPart[len(aPart)-1]) {
			fullPaths[0] = aPart
		} else {
			fullPaths, err = walkThroughPart(fullPaths, aPart, i != lastPartIndex)
			if err != nil {
				return nil, err
			}
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
		var testPath string

		testPath = filepath.Join(aFullPath, part)

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
			result = append(result, aMatch)
		}
	}

	return result, nil
}
