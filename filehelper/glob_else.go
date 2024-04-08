//go:build !windows

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
	"os"
	"path/filepath"
)

// ******** Private functions ********

// platformGlobWithSwitch is the globbing function for all OSes except Windows.
func platformGlobWithSwitch(pattern string, withDirs bool, withFiles bool) ([]string, error) {
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
