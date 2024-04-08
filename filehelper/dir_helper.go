//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//

package filehelper

import (
	"filesigner/slicehelper"
	"filesigner/stringhelper"
	"os"
	"path/filepath"
	"runtime"
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

		parts = slicehelper.Prepend(parts, file)
		lastDirIndex := dirLen - 1
		if dirLen > 1 && os.IsPathSeparator(dir[lastDirIndex]) {
			dir = dir[:lastDirIndex]
		}
	}

	if dirLen != 0 {
		parts = slicehelper.Prepend(parts, dir)
	}

	if volLen != 0 {
		parts[0] = vol + parts[0]
	}

	return
}

// PathGlob returns all globbed files from a path that may contain wildcards.
func PathGlob(path string, excludeDirList []string, excludeFileList []string) ([]string, error) {
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

	fullPaths := emptyFullPaths()

	var isExcluded bool
	for i, part := range parts {
		if os.IsPathSeparator(part[len(part)-1]) {
			if runtime.GOOS == `windows` {
				EnsureDriveLetterIsUpperCase(part)
			}

			fullPaths[0] = part
		} else {
			findDirs := i != lastPartIndex

			isExcluded, err = isNameExcluded(part, findDirs, excludeDirList, excludeFileList)
			if err != nil {
				return nil, err
			}

			if isExcluded {
				return nil, nil
			}

			fullPaths, err = walkThroughPart(fullPaths, part, findDirs)
			if err != nil {
				return nil, err
			}
		}
	}

	return fullPaths, nil
}

// EnsureDriveLetterIsUpperCase ensures that the drive letter is an upper-case letter.
func EnsureDriveLetterIsUpperCase(path string) {
	if path[1] == ':' {
		driveLetter := path[0]
		if driveLetter >= 'a' && driveLetter <= 'z' {
			aPartBytes := stringhelper.UnsafeStringBytes(path)
			aPartBytes[0] ^= 0x20
		}
	}
}

// ******** Private functions ********

// emptyFullPaths returns a string slice with one element that is a string of length 0.
// This is the starting point for the fullPaths string slice.
func emptyFullPaths() []string {
	result := make([]string, 1)
	result[0] = ``
	return result
}

// isNameExcluded returns "true", if the name is member of an exclude list and "false", otherwise.
func isNameExcluded(name string, findDirs bool, excludeDirList []string, excludeFileList []string) (bool, error) {
	var err error
	var isExcluded bool

	if findDirs {
		isExcluded, err = MatchesAny(excludeDirList, name)
	} else {
		isExcluded, err = MatchesAny(excludeFileList, name)
	}

	return isExcluded, err
}

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
