//go:build windows

package filehelper

import (
	"errors"
	"golang.org/x/sys/windows"
)

// SensibleGlob globs a pattern with Windows API calls as this is the only correct
// way to handle globbing on the case-insensitive Windows file system.
func SensibleGlob(pattern string) ([]string, error) {
	// Initialize result
	result := make([]string, 0, 64)

	// Remove trailing separators, if any
	pattern = ensureNoTrailingSeparator(pattern)

	// Return if pattern is empty
	if len(pattern) == 0 {
		return result, nil
	}

	// Convert pattern into a UTF-16 string
	patternUTF16Ptr, _ := windows.UTF16PtrFromString(pattern)
	var findData windows.Win32finddata

	// See if there is a match
	findHandle, err := windows.FindFirstFile(patternUTF16Ptr, &findData)
	if err != nil {
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
			return result, nil
		} else {
			return result, err
		}
	}

	// Ensure that the find handle is closed on exit
	defer findCloseHelper(findHandle)

	// Append first file to result
	result = appendFileNameIfNoDir(result, findData)

	// Now loop through more matching files
	for {
		err = windows.FindNextFile(findHandle, &findData)

		switch err {
		case nil:
			result = appendFileNameIfNoDir(result, findData)

		case windows.ERROR_NO_MORE_FILES:
			return result, nil

		default:
			return result, err
		}
	}
}

// appendFileNameIfNoDir appends a file name to the result list, if it does not represent a directory.
func appendFileNameIfNoDir(result []string, findData windows.Win32finddata) []string {
	if (findData.FileAttributes & windows.FILE_ATTRIBUTE_DIRECTORY) == 0 {
		result = append(result, windows.UTF16ToString(findData.FileName[:]))
	}

	return result
}

// findCloseHelper wraps windows.FindClose which will never return an error with this application
func findCloseHelper(handle windows.Handle) {
	_ = windows.FindClose(handle)
}
