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

package main

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/filesignature"
	"filesigner/hashsignature"
	"filesigner/keyid"
	"filesigner/logger"
	"filesigner/signaturefile"
	"filesigner/signaturehandler"
	"filesigner/texthelper"
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
)

// ******** Private functions ********

// doVerification verifies a signature file.
func doVerification(contextId string, signaturesFileName string) int {
	signatureData, err := signaturefile.ReadSignatureFile(signaturesFileName)
	if err != nil {
		logger.PrintError(51, err.Error())
		return rcProcessError
	}

	var hashVerifier hashsignature.HashVerifier
	var publicKeyId string
	hashVerifier, publicKeyId, err = getHashVerifier(signatureData)
	if err != nil {
		logger.PrintError(52, err.Error())
		return rcProcessError
	}

	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextId)
	if err == nil {
		if !ok {
			logger.PrintError(53, `signature file has been tampered with or wrong context id`)
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(54, `error verifying signature file data signature: %v`, err)
		return rcProcessError
	}

	printMetaData(contextId, publicKeyId, signatureData.Timestamp, signatureData.Hostname)

	successCount, errorCount, rc := verifyFiles(contextId, signatureData, hashVerifier)

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(59, "Verification of %d file%s successful", successCount, successEnding)

	case rcProcessWarning:
		logger.PrintWarningf(60, "Verification of %d file%s successful and warnings present", successCount, successEnding)

	case rcProcessError:
		logger.PrintErrorf(61, "Verification of %d file%s successful and %d file%s unsuccessful", successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

// verifyFiles verifies the signatures of the files in the signature data.
func verifyFiles(contextId string,
	signatureData *signaturehandler.SignatureData,
	hashVerifier hashsignature.HashVerifier) (int, int, int) {
	filePaths, rc := getExistingFiles(maps.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(62, "No files from signature file present")
		return 0, 0, rcProcessWarning
	}

	hashList := filehasher.FileHashes(filePaths, contextId)
	if existHashErrors(hashList) {
		return 0, 0, rcProcessError
	}

	successList, errorList := filesignature.VerifyFileHashes(hashVerifier, signatureData.FileSignatures, hashList)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList("Verification", successList)
	}

	errorCount := len(errorList)
	if errorCount > 0 {
		printErrorList(errorList)
		rc = rcProcessError
	}

	return successCount, errorCount, rc
}

// getHashVerifier constructs the hash verifier and the key id from the signature data.
func getHashVerifier(signatureData *signaturehandler.SignatureData) (hashsignature.HashVerifier, string, error) {
	publicKeyBytes, err := base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		return nil, "", fmt.Errorf("Could not convert public key to bytes: %w", err)
	}

	var hashVerifier hashsignature.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsignature.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsignature.NewEc521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		return nil, "", fmt.Errorf("Could not create hash verifier: %w", err)
	}

	return hashVerifier, keyid.KeyId(publicKeyBytes), nil
}

// getExistingFiles gets the files from a signature list that exist in the directory that is to be verified.
func getExistingFiles(filePaths []string) ([]string, int) {
	rc := rcOK

	result := make([]string, 0, len(filePaths))
	for _, fp := range filePaths {
		nfp := filepath.FromSlash(fp)
		fi, err := os.Stat(nfp)
		if err != nil {
			logger.PrintWarningf(63, "'%s' from signature file does not exist", nfp)
			rc = rcProcessWarning
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(64, "'%s' from signature file is a directory", nfp)
				rc = rcProcessWarning
			} else {
				result = append(result, nfp)
			}
		}
	}

	return result, rc
}
