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

import "os"

// ******** Public functions ********

// SensibleGlobFilesOnly globs a pattern.
// This call returns only files, no directories.
func SensibleGlobFilesOnly(pattern string) ([]string, error) {
	return platformGlobWithSwitch(pattern, false, true)
}

// SensibleGlobDirsOnly globs a pattern.
// This call returns only directories, no files.
func SensibleGlobDirsOnly(pattern string) ([]string, error) {
	return platformGlobWithSwitch(pattern, true, false)
}

// SensibleGlob globs a pattern.
func SensibleGlob(pattern string) ([]string, error) {
	return platformGlobWithSwitch(pattern, true, true)
}

// ******** Private functions ********

// ensureNoTrailingSeparator ensures that the pattern does not end with a trailing separator.
func ensureNoTrailingSeparator(pattern string) string {
	var pos int
	for pos = len(pattern) - 1; pos >= 0; {
		if os.IsPathSeparator(pattern[pos]) {
			pos--
		} else {
			break
		}
	}

	return pattern[:pos+1]
}
