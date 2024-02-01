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
	"filesigner/flaglist"
	"filesigner/set"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
)

var modExcludeFileNameList []string
var modIncludeFileNameList []string
var modExcludeDirNameList []string
var modIncludeDirNameList []string
var modDoRecursion bool
var modResultList *set.Set[string]
var modMatchFunc func(string, string) (bool, error)

func ScanDir(includeFileList *flaglist.FlagList,
	excludeFileList *flaglist.FlagList,
	includeDirList *flaglist.FlagList,
	excludeDirList *flaglist.FlagList,
	doRecursion bool) (*set.Set[string], error) {
	modIncludeFileNameList = includeFileList.GetNames()
	modExcludeFileNameList = excludeFileList.GetNames()
	modIncludeDirNameList = includeDirList.GetNames()
	modExcludeDirNameList = excludeDirList.GetNames()
	modDoRecursion = doRecursion
	modResultList = set.New[string]()
	if runtime.GOOS == "windows" {
		modMatchFunc = CaseInvariantMatchFunction
	} else {
		modMatchFunc = filepath.Match
	}

	// Always walk the current directory
	return modResultList, filepath.WalkDir(".", WalkEntryFunction)
}

func WalkEntryFunction(path string, dirEntry fs.DirEntry, dirErr error) error {
	// Return immediately if walking the directory tree returned an error
	if dirErr != nil {
		return dirErr
	}

	// Never process current or parent directory
	if path == "." || path == ".." {
		return nil
	}

	isDir := dirEntry.IsDir()

	// If the entry is a directory and subdirectories are not allowed return SkipDir
	if isDir && !modDoRecursion {
		return filepath.SkipDir
	}

	entryName := dirEntry.Name()

	var shouldProcess bool
	var err error

	if !isDir {
		shouldProcess, err = shouldProcessEntry(entryName,
			modIncludeFileNameList,
			modExcludeFileNameList)
	} else {
		shouldProcess, err = shouldProcessEntry(entryName,
			modIncludeDirNameList,
			modExcludeDirNameList)
	}

	if err != nil {
		return err
	}

	if !shouldProcess {
		if !isDir {
			return nil
		} else {
			return filepath.SkipDir
		}
	}

	// Add the entry to the result list. We come here if the entry is not excluded.
	// If there are includes the entry also has to match an include specification.
	// Only add files, not directories.
	if !isDir {
		modResultList.Add(path)
	}

	return nil
}

func shouldProcessEntry(entryName string, includeNames []string, excludeNames []string) (bool, error) {
	var isEntryInList bool
	var err error

	if len(excludeNames) != 0 {
		isEntryInList, err = matchesAnyInList(entryName, excludeNames)
		if isEntryInList || err != nil {
			return false, err
		}
	}

	if len(includeNames) != 0 {
		isEntryInList, err = matchesAnyInList(entryName, includeNames)
		if !isEntryInList || err != nil {
			return false, err
		}
	}

	return true, nil
}

func matchesAnyInList(entryName string, nameList []string) (bool, error) {
	isMatch, err := matchesEntry(entryName, nameList)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}

func CaseInvariantMatchFunction(pattern string, name string) (bool, error) {
	return filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
}

func matchesEntry(name string, entries []string) (bool, error) {
	for _, entry := range entries {
		isMatch, err := modMatchFunc(entry, name)
		if err != nil {
			return false, err
		}

		if isMatch {
			return true, nil
		}
	}

	return false, nil
}
