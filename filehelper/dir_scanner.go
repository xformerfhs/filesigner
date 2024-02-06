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
	"filesigner/flaglist"
	"filesigner/set"
	"io/fs"
	"path/filepath"
)

var modExcludeFileNameList []string
var modIncludeFileNameList []string
var modExcludeDirNameList []string
var modIncludeDirNameList []string
var modDoRecursion bool
var modResultList *set.FileSystemStringSet

func ScanDir(includeFileList *flaglist.FileSystemFlagList,
	excludeFileList *flaglist.FileSystemFlagList,
	includeDirList *flaglist.FileSystemFlagList,
	excludeDirList *flaglist.FileSystemFlagList,
	doRecursion bool) (*set.FileSystemStringSet, error) {
	modIncludeFileNameList = includeFileList.Elements()
	modExcludeFileNameList = excludeFileList.Elements()
	modIncludeDirNameList = includeDirList.Elements()
	modExcludeDirNameList = excludeDirList.Elements()
	modDoRecursion = doRecursion
	modResultList = set.NewFileSystemStringSet()

	// Always walk the current directory
	return modResultList, filepath.WalkDir(".", WalkEntryFunction)
}

func WalkEntryFunction(path string, dirEntry fs.DirEntry, dirErr error) error {
	// Return immediately if walking the directory tree returned an error
	if dirErr != nil {
		return dirErr
	}

	// Never process current or parent directory
	if path == `.` || path == `..` {
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
		isEntryInList, err = MatchesAny(excludeNames, entryName)
		if isEntryInList || err != nil {
			return false, err
		}
	}

	if len(includeNames) != 0 {
		isEntryInList, err = MatchesAny(includeNames, entryName)
		if !isEntryInList || err != nil {
			return false, err
		}
	}

	return true, nil
}
