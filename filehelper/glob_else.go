//go:build !windows

//
// SPDX-FileCopyrightText: Copyright 2024-2026 Frank Schwab
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
// Version: 2.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2026-09-02: V1.1.0: Check, if file is regular.
//    2026-09-02: V2.0.0: Only follow symbolic links, if followSymlinks is true.
//

package filehelper

import (
	"os"
	"path/filepath"
)

// ******** Private functions ********

// platformGlobWithSwitch is the globbing function for all OSes except Windows.
// Symbolic links are only resolved, i.e. classified and returned by the type of their
// target, if followSymlinks is true. Otherwise, they are skipped.
func platformGlobWithSwitch(pattern string, withDirs bool, withFiles bool, followSymlinks bool) ([]string, error) {
	// Remove trailing separators, if any
	pattern = ensureNoTrailingSeparator(pattern)

	// Find files and directories matching pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return matches, err
	}

	return filterMatches(matches, withDirs, withFiles, followSymlinks)
}

// filterMatches filters all elements in the globbing list according to being directories or files.
func filterMatches(matchList []string, withDirs bool, withFiles bool, followSymlinks bool) ([]string, error) {
	// Stat follows symbolic links, Lstat does not. With Lstat a symbolic link is neither a
	// regular file nor a directory, so it is filtered out below.
	statFunc := os.Lstat
	if followSymlinks {
		statFunc = os.Stat
	}

	result := make([]string, 0, len(matchList))
	for _, fileName := range matchList {
		fi, err := statFunc(fileName)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			if withFiles && fi.Mode().IsRegular() {
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
