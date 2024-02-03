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
//    2024-02-03: V1.0.0: Created.
//

package filehelper

import (
	"path/filepath"
	"runtime"
	"strings"
)

// ******** Private variables ********

// matcherMatchFunc is the pointer to the platform-dependent matcher function.
var matcherMatchFunc func(string, string) (bool, error)

// isMatcherInitialized is true, if the matcher has been initialized.
var isMatcherInitialized bool

// ******** Public functions ********

// Matches returns true if the pattern matches the given name.
func Matches(pattern string, name string) (bool, error) {
	ensureMatcherIsInitialized()

	return matcherMatchFunc(pattern, name)
}

// MatchesAny returns true if any pattern matches the given name.
func MatchesAny(patterns []string, name string) (bool, error) {
	ensureMatcherIsInitialized()

	for _, entry := range patterns {
		isMatch, err := matcherMatchFunc(entry, name)
		if err != nil {
			return false, err
		}

		if isMatch {
			return true, nil
		}
	}

	return false, nil
}

// ******** Private functions ********

// ensureMatcherIsInitialized ensures, that the matcher is initialized.
func ensureMatcherIsInitialized() {
	if !isMatcherInitialized {
		initMatcher()
		isMatcherInitialized = true
	}
}

// initMatcher initializes the platform-dependent match function of the matcher.
func initMatcher() {
	if runtime.GOOS == "windows" {
		modMatchFunc = caseInsensitiveMatchFunction
	} else {
		modMatchFunc = filepath.Match
	}
}

// caseInsensitiveMatchFunction is a filepath.Match-equivalent function for case-insensitive file systems.
func caseInsensitiveMatchFunction(pattern string, name string) (bool, error) {
	return filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
}
