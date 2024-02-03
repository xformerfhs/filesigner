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

package filesignature

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/hashsignature"
	"filesigner/maphelper"
	"fmt"
	"path/filepath"
)

// ******** Public functions ********

// VerifyFileHashes verifies file hashes
func VerifyFileHashes(hashVerifier hashsignature.HashVerifier,
	fileSignatures map[string]string,
	fileHashList map[string]*filehasher.HashResult) ([]string, []error) {
	var err error

	successCollection := make([]string, 0, len(fileHashList))
	errCollection := make([]error, 0, len(fileHashList))

	filePathList := maphelper.SortedKeys(fileSignatures)

	var signatureString string
	var signatureValue []byte
	for _, filePath := range filePathList {
		normalizedFilePath := filepath.FromSlash(filePath)
		fileHashResult, haveHashForFilePath := fileHashList[normalizedFilePath]
		if haveHashForFilePath {
			signatureString = fileSignatures[filePath]
			signatureValue, err = base32encoding.DecodeFromString(signatureString)
			if err != nil {
				errCollection = append(errCollection, fmt.Errorf("Signature of file '%s' has invalid encoding: %w", normalizedFilePath, err))
			} else {
				var ok bool
				ok, err = hashVerifier.VerifyHash(fileHashResult.HashValue, signatureValue)
				if ok {
					successCollection = append(successCollection, normalizedFilePath)
				} else {
					errCollection = append(errCollection, fmt.Errorf("File '%s' has been tampered with", normalizedFilePath))
				}
			}
		}
	}

	return successCollection, errCollection
}
