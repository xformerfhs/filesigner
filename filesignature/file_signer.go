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

// SignFileHashes creates signatures for file hashes.
func SignFileHashes(hashSigner hashsignature.HashSigner,
	hashResultList map[string]*filehasher.HashResult) (map[string]string, []string, error) {
	filePathList := maphelper.SortedKeys(hashResultList)

	return makeHashSignatures(hashSigner, filePathList, hashResultList)
}

// ******** Private functions ********

func makeHashSignatures(hashSigner hashsignature.HashSigner,
	filePathList []string,
	hashResultList map[string]*filehasher.HashResult) (map[string]string, []string, error) {
	var err error
	signatures := make(map[string]string, len(filePathList))

	var signature []byte
	for _, filePath := range filePathList {
		signature, err = hashSigner.SignHash(hashResultList[filePath].HashValue)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not sign hash of file '%s': %w", filePath, err)
		}

		signatures[filepath.ToSlash(filePath)] = base32encoding.EncodeToString(signature)
	}

	return signatures, filePathList, nil
}
