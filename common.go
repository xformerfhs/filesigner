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

package main

import (
	"filesigner/filehasher"
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
func existHashErrors(hashResults map[string]*filehasher.HashResult) bool {
	result := false

	keyList := maphelper.SortedKeys(hashResults)

	var hr *filehasher.HashResult
	for _, filePath := range keyList {
		hr = hashResults[filePath]
		if hr.Err != nil {
			logger.PrintErrorf(23, "Could not get hash of file '%s': %v", hr.FilePath, hr.Err)
			result = true
		}
	}

	return result
}
