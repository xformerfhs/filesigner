//
// SPDX-FileCopyrightText: Copyright 2023 Frank Schwab
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
			if runtime.GOOS == `windows` {
				ensureDriveLetterIsUpperCase(aPart)
			}

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

// ensureDriveLetterIsUpperCase ensures that the drive letter is an upper-case letter.
func ensureDriveLetterIsUpperCase(aPart string) {
	if aPart[1] == ':' {
		driveLetter := aPart[0]
		if driveLetter >= 'a' && driveLetter <= 'z' {
			aPartBytes := stringhelper.UnsafeStringBytes(aPart)
			aPartBytes[0] ^= 0x20
		}
	}
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
