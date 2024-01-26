package filehelper

import (
	"os"
	"strings"
)

// ******** Private variables ********

var isCaseSensitivityKnown bool
var fileSystemIsCaseSensitive bool

// ******** Public functions ********

// IsFileSystemCaseSensitive determines whether the file system is case-sensitive or not.
//
// Returns true if the file system is case-sensitive, false otherwise.
// Returns an error if any error occurred during the process.
func IsFileSystemCaseSensitive() (bool, error) {
	if isCaseSensitivityKnown {
		return fileSystemIsCaseSensitive, nil
	} else {
		var err error
		fileSystemIsCaseSensitive, err = isFileSystemCaseSensitive()
		if err != nil {
			return false, err
		}

		isCaseSensitivityKnown = true
		return fileSystemIsCaseSensitive, nil
	}
}

// ******** Private functions ********

// isFileSystemCaseSensitive determines whether the file system is case-sensitive or not.
//
// Returns true if the file system is case-sensitive, false otherwise.
// Returns an error if any error occurred during the process.
func isFileSystemCaseSensitive() (bool, error) {
	const testPattern = "fScSc"

	testFile, err := os.CreateTemp(".", testPattern)
	if err != nil {
		return false, err
	}
	testFilePath := testFile.Name()
	CloseFile(testFile)
	defer DeleteFile(testFilePath)

	_, err = os.Stat(strings.ToLower(testFilePath))
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		} else {
			return false, err
		}
	}

	return false, nil
}
