package filehelper

import (
	"filesigner/typelist"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
)

var modExcludeFileNameList []string
var modIncludeFileNameList []string
var modExcludeDirNameList []string
var modIncludeDirNameList []string
var modNoSubDirs bool
var modResultList []string
var modMatchFunc func(string, string) (bool, error)

func ScanDir(includeFileList *typelist.FlagTypeList,
	excludeFileList *typelist.FlagTypeList,
	includeDirList *typelist.FlagTypeList,
	excludeDirList *typelist.FlagTypeList,
	noSubDirs bool) ([]string, error) {
	modIncludeFileNameList = includeFileList.GetNames()
	modExcludeFileNameList = excludeFileList.GetNames()
	modIncludeDirNameList = includeDirList.GetNames()
	modExcludeDirNameList = excludeDirList.GetNames()
	modNoSubDirs = noSubDirs
	modResultList = make([]string, 0, 100)
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
	if isDir && modNoSubDirs {
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
		return nil
	}

	// Add the entry to the result list. We come here if the entry is not excluded.
	// If there are includes the entry also has to match an include specification.
	// Only add files, not directories.
	if !isDir {
		modResultList = append(modResultList, path)
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
