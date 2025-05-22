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
// Version: 1.4.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-17: V1.1.0: Create contextBytes as early, as possible.
//    2024-03-04: V1.2.0: Get public key bytes, not id.
//    2025-03-01: V1.2.1: Correct message levels of verification success messages.
//    2025-03-01: V1.3.0: Add message base.
//    2025-03-01: V1.4.0: Correct handling of os.Stat errors.
//

package main

import (
	"errors"
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

// ******** Private constants ********

// verifyCmdMsgBase is the base number for all messages in verify_command.
// This file reserves numbers 50-69.
const verifyCmdMsgBase = 50

const errMsgCouldNotConvert = `Could not convert %s to bytes: %v`

// ******** Private functions ********

// doVerification verifies a signatures file.
func doVerification(signaturesFileName string) int {
	logger.PrintInfof(verifyCmdMsgBase+1, `Reading signatures file '%s'`, signaturesFileName)

	signatureData, err := signaturefile.ReadJson(signaturesFileName)
	if err != nil {
		logger.PrintErrorf(verifyCmdMsgBase+2, `Error reading signatures file: %v`, err)
		return rcProcessError
	}

	var publicKeyBytes []byte
	publicKeyBytes, err = base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		logger.PrintErrorf(verifyCmdMsgBase+3, errMsgCouldNotConvert, `public key`, err)
		return rcProcessError
	}

	var dataSignature []byte
	dataSignature, err = base32encoding.DecodeFromString(signatureData.DataSignature)
	if err != nil {
		logger.PrintErrorf(verifyCmdMsgBase+3, errMsgCouldNotConvert, `data signature`, err)
		return rcProcessError
	}

	var hashVerifier hashsignature.HashVerifier
	hashVerifier, err = getHashVerifier(signatureData, publicKeyBytes)
	if err != nil {
		logger.PrintErrorf(verifyCmdMsgBase+3, `Error getting hash verifier: %v`, err)
		return rcProcessError
	}

	contextKey := stretcher.KeyFromBytes(stringhelper.UnsafeStringBytes(signatureData.ContextId))
	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextKey, dataSignature)
	if err == nil {
		if !ok {
			logger.PrintError(verifyCmdMsgBase+4, `Signatures file has been modified`)
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(verifyCmdMsgBase+5, `Error verifying signatures file data signature: %v`, err)
		return rcProcessError
	}

	printMetaData(signatureData, publicKeyBytes)

	successCount, errorCount, rc := verifyFiles(contextKey, signatureData, hashVerifier)

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(verifyCmdMsgBase+6, `Verification of %d file%s successful`, successCount, successEnding)

	case rcProcessWarning:
		logger.PrintInfof(verifyCmdMsgBase+7, `Verification of %d file%s successful and warnings present`, successCount, successEnding)

	case rcProcessError:
		logger.PrintInfof(verifyCmdMsgBase+8, `Verification of %d file%s successful and %d file%s unsuccessful`, successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

// verifyFiles verifies the signatures of the files in the signature data.
func verifyFiles(contextBytes []byte,
	signatureData *signaturehandler.SignatureData,
	hashVerifier hashsignature.HashVerifier) (int, int, int) {
	filePaths, rc := getExistingFiles(maphelper.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(verifyCmdMsgBase+9, `No files from signatures file present`)
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
func getHashVerifier(signatureData *signaturehandler.SignatureData, publicKeyBytes []byte) (hashsignature.HashVerifier, error) {
	var err error
	var hashVerifier hashsignature.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsignature.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsignature.NewEcDsaP521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		return nil, fmt.Errorf(`Could not create hash verifier: %w`, err)
	}

	return hashVerifier, nil
}

// getExistingFiles gets the files from a signature list that exist in the directory that is to be verified.
func getExistingFiles(filePaths []string) ([]string, int) {
	rc := rcOK

	result := make([]string, 0, len(filePaths))
	for _, fp := range filePaths {
		nfp := filepath.FromSlash(fp)
		fi, err := os.Stat(nfp)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				logger.PrintWarningf(verifyCmdMsgBase+10, `File '%s' in signatures file does not exist`, nfp)
				rc = max(rc, rcProcessWarning)
			} else {
				logger.PrintErrorf(verifyCmdMsgBase+11, `Error checking if file '%s' in signatures file exists: %v`, nfp, err)
				rc = rcProcessError
			}
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(verifyCmdMsgBase+12, `'%s' in signatures file is a directory`, nfp)
				rc = max(rc, rcProcessWarning)
			} else {
				result = append(result, nfp)
			}
		}
	}

	return result, rc
}
