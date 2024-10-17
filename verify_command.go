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
// Version: 1.1.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-17: V1.1.0: Create contextBytes as early, as possible.
//    2024-03-04: V1.2.0: Get public key bytes, not id.
//

package main

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/filesignature"
	"filesigner/hashsignature"
	"filesigner/logger"
	"filesigner/maphelper"
	"filesigner/signaturefile"
	"filesigner/signaturehandler"
	"filesigner/stretcher"
	"filesigner/stringhelper"
	"filesigner/texthelper"
	"fmt"
	"os"
	"path/filepath"
)

// ******** Private functions ********

// doVerification verifies a signatures file.
func doVerification(signaturesFileName string) int {
	logger.PrintInfof(51, `Reading signatures file '%s'`, signaturesFileName)

	signatureData, err := signaturefile.ReadJson(signaturesFileName)
	if err != nil {
		logger.PrintErrorf(52, `Error reading signatures file: %v`, err)
		return rcProcessError
	}

	var hashVerifier hashsignature.HashVerifier
	var publicKeyId []byte
	hashVerifier, publicKeyId, err = getHashVerifier(signatureData)
	if err != nil {
		logger.PrintErrorf(53, `Error getting hash verifier: %v`, err)
		return rcProcessError
	}

	contextKey := stretcher.KeyFromBytes(stringhelper.UnsafeStringBytes(signatureData.ContextId))
	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextKey)
	if err == nil {
		if !ok {
			logger.PrintError(54, `Signatures file has been modified`)
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(55, `Error verifying signatures file data signature: %v`, err)
		return rcProcessError
	}

	printMetaData(signatureData, publicKeyId)

	successCount, errorCount, rc := verifyFiles(contextKey, signatureData, hashVerifier)

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(56, `Verification of %d file%s successful`, successCount, successEnding)

	case rcProcessWarning:
		logger.PrintWarningf(57, `Verification of %d file%s successful and warnings present`, successCount, successEnding)

	case rcProcessError:
		logger.PrintErrorf(58, `Verification of %d file%s successful and %d file%s unsuccessful`, successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

// verifyFiles verifies the signatures of the files in the signature data.
func verifyFiles(contextBytes []byte,
	signatureData *signaturehandler.SignatureData,
	hashVerifier hashsignature.HashVerifier) (int, int, int) {
	filePaths, rc := getExistingFiles(maphelper.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(59, `No files from signatures file present`)
		return 0, 0, rcProcessWarning
	}

	hashList := filehasher.FileHashes(filePaths, contextBytes)
	if existHashErrors(hashList) {
		return 0, 0, rcProcessError
	}

	successList, errorList := filesignature.VerifyFileHashes(hashVerifier, signatureData.FileSignatures, hashList)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList(`Verification`, successList)
	}

	errorCount := len(errorList)
	if errorCount > 0 {
		printErrorList(errorList)
		rc = rcProcessError
	}

	return successCount, errorCount, rc
}

// getHashVerifier constructs the hash verifier and the key id from the signature data.
func getHashVerifier(signatureData *signaturehandler.SignatureData) (hashsignature.HashVerifier, []byte, error) {
	publicKeyBytes, err := base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf(`Could not convert public key to bytes: %w`, err)
	}

	var hashVerifier hashsignature.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsignature.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsignature.NewEcDsaP521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		return nil, nil, fmt.Errorf(`Could not create hash verifier: %w`, err)
	}

	return hashVerifier, publicKeyBytes, nil
}

// getExistingFiles gets the files from a signature list that exist in the directory that is to be verified.
func getExistingFiles(filePaths []string) ([]string, int) {
	rc := rcOK

	result := make([]string, 0, len(filePaths))
	for _, fp := range filePaths {
		nfp := filepath.FromSlash(fp)
		fi, err := os.Stat(nfp)
		if err != nil {
			logger.PrintWarningf(60, `File '%s' in signatures file does not exist`, nfp)
			rc = rcProcessWarning
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(61, `'%s' in signatures file is a directory`, nfp)
				rc = rcProcessWarning
			} else {
				result = append(result, nfp)
			}
		}
	}

	return result, rc
}
