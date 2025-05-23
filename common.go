//
// SPDX-FileCopyrightText: Copyright 2024-2025 Frank Schwab
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
// Version: 1.1.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2025-03-01: V1.1.0: Add message base.
//

package main

import (
	"filesigner/filehasher"
	"filesigner/keyid"
	"filesigner/logger"
	"filesigner/maphelper"
	"filesigner/signaturehandler"
	"filesigner/stringhelper"
	"sort"
)

// ******** Private constants ********

// commonMsgBase is the base number for all messages in common.
// This file reserves numbers 20-29.
const commonMsgBase = 20

// ******** Private functions ********

// printSuccessList prints the successful executions of an operation.
func printSuccessList(operation string, successList []string) {
	sort.Strings(successList)

	for _, filePath := range successList {
		logger.PrintInfof(commonMsgBase+1, `%s succeeded for file '%s'`, operation, filePath)
	}
}

// printErrorList prints the errors that occurred during an operation.
func printErrorList(errorList []error) {
	for _, err := range errorList {
		logger.PrintError(commonMsgBase+2, err.Error())
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
			logger.PrintErrorf(commonMsgBase+3, `Could not get hash of file '%s': %v`, hr.FilePath, hr.Err)
			result = true
		}
	}

	return result
}

// printMetaData prints the meta data of the signatures.
func printMetaData(
	signatureData *signaturehandler.SignatureData,
	publicKeyBytes []byte) {
	logger.PrintInfof(commonMsgBase+4, `Context id         : %s`, signatureData.ContextId)
	logger.PrintInfof(commonMsgBase+5, `Public key id      : %s`, keyid.KeyId(publicKeyBytes))
	logger.PrintInfof(commonMsgBase+6, `Signature timestamp: %s`, signatureData.Timestamp)
	logger.PrintInfof(commonMsgBase+7, `Signature host name: %s`, signatureData.Hostname)
	logger.PrintInfof(commonMsgBase+8, `Verifier           : %s`, verifyString(signatureData, publicKeyBytes))
}

// verifyString returns the verify string for the given data.
func verifyString(
	signatureData *signaturehandler.SignatureData,
	publicKeyBytes []byte) string {
	return keyid.KeyId(
		stringhelper.UnsafeStringBytes(signatureData.ContextId),
		publicKeyBytes,
		stringhelper.UnsafeStringBytes(signatureData.Timestamp),
		stringhelper.UnsafeStringBytes(signatureData.Hostname),
	)
}
