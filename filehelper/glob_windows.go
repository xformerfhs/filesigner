//go:build windows

package filehelper

import (
	"errors"
	"golang.org/x/sys/windows"
	"path/filepath"
)

// ******** Private functions ********

// sensibleGlobWithSwitch globs a pattern with Windows API calls as this is the only correct
// way to handle globbing on the case-insensitive Windows file system.
func sensibleGlobWithSwitch(pattern string, withDirs bool, withFiles bool) ([]string, error) {
	// Initialize result
	var result []string

	// Remove trailing separators, if any.
	pattern = ensureNoTrailingSeparator(pattern)

	// Return if pattern is empty
	if len(pattern) == 0 {
		return result, nil
	}

	// Split path into directory and file.
	patternDir, patternFile := filepath.Split(pattern)

	// "." and ".." must not be handled by the Windows API calls as they will replace them with
	// the directory names, which is wrong when constructing paths.
	if patternFile == `.` || patternFile == `..` {
		result = append(result, pattern)
		return result, nil
	}

	// Convert pattern into a UTF-16 string.
	patternUTF16Ptr, _ := windows.UTF16PtrFromString(pattern)
	var findData windows.Win32finddata

	// See if there is a match.
	findHandle, err := windows.FindFirstFile(patternUTF16Ptr, &findData)
	if err != nil {
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
			return result, nil
		} else {
			return result, err
		}
	}

	// Ensure that the find handle is closed on exit.
	defer findCloseHelper(findHandle)

	// Append first file to result.
	result = appendNameIfEligible(patternDir, result, findData, withDirs, withFiles)

	// Now loop through more matching files.
	for {
		err = windows.FindNextFile(findHandle, &findData)

		switch {
		case err == nil:
			result = appendNameIfEligible(patternDir, result, findData, withDirs, withFiles)
		case errors.Is(err, windows.ERROR_NO_MORE_FILES):
			return result, nil
		default:
			return result, err
		}
	}
}

// appendNameIfEligible appends a file name to the result list, if it is eligible.
func appendNameIfEligible(patternDir string, result []string, findData windows.Win32finddata, withDirs bool, withFiles bool) []string {
	fileName := windows.UTF16ToString(findData.FileName[:])

	if (findData.FileAttributes & windows.FILE_ATTRIBUTE_DIRECTORY) == 0 {
		if withFiles {
			result = append(result, filepath.Join(patternDir, fileName))
		}
	} else {
		if withDirs {
			// Never return "." or "..".
			if fileName != `.` && fileName != `..` {
				result = append(result, filepath.Join(patternDir, fileName))
			}
		}
	}

	return result
}

// findCloseHelper wraps windows.FindClose which will never return an error with this application
func findCloseHelper(handle windows.Handle) {
	_ = windows.FindClose(handle)
}
