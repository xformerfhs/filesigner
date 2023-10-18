package main

import (
	"filesigner/filehashing"
	"filesigner/logger"
	"filesigner/maphelper"
	"sort"
)

// printSuccessList prints the successful executions of an operation.
func printSuccessList(operation string, successList []string) {
	sort.Strings(successList)

	for _, filePath := range successList {
		logger.PrintInfof(21, "%s succeeded for file '%s'", operation, filePath)
	}
}

// printErrorList prints the errors that occurred during an operation.
func printErrorList(errorList []error) {
	for _, err := range errorList {
		logger.PrintError(22, err.Error())
	}
}

// existHashErrors checks if hash errors exist and prints them.
func existHashErrors(hashResults map[string]*filehashing.HashResult) bool {
	result := false

	keyList := maphelper.GetSortedKeys(hashResults)

	var hr *filehashing.HashResult
	for _, filePath := range keyList {
		hr = hashResults[filePath]
		if hr.Err != nil {
			logger.PrintErrorf(23, "Could not get hash of file '%s': %v", hr.FilePath, hr.Err)
			result = true
		}
	}

	return result
}
