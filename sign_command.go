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
	"os"
	"time"
)

// ******** Private constants ********

// timeStampFormat Format for signature file time stamp
const timeStampFormat = "2006-01-02 15:04:05 Z07:00"

// ******** Private functions ********

// doSigning signs all files with the given context id.
func doSigning(
	signaturesFileName string,
	signatureType signaturehandler.SignatureType,
	contextId string,
	filePaths []string,
) int {
	var err error

	signatureData := &signaturehandler.SignatureData{
		Format:        signaturehandler.SignatureFormatV1,
		Timestamp:     time.Now().Format(timeStampFormat),
		SignatureType: signatureType,
	}

	signatureData.Hostname, err = os.Hostname()
	if err != nil {
		logger.PrintErrorf(31, `could not get host name: %v`, err)
		return rcProcessError
	}

	resultList := filehasher.FileHashes(filePaths, contextId)

	if existHashErrors(resultList) {
		return rcProcessError
	}

	var hashSigner hashsignature.HashSigner
	if signatureType == signaturehandler.SignatureTypeEd25519 {
		hashSigner, err = hashsignature.NewEd25519HashSigner()
	} else {
		hashSigner, err = hashsignature.NewEcDsa521HashSigner()
	}
	if err != nil {
		logger.PrintErrorf(32, `could not create hash-signer: %v`, err)
		return rcProcessError
	}
	defer hashSigner.Destroy()

	var publicKeyBytes []byte
	publicKeyBytes, err = hashSigner.GetPublicKey()
	if err != nil {
		logger.PrintErrorf(33, `could not get public key bytes: %v`, err)
		return rcProcessError
	}
	signatureData.PublicKey = base32encoding.EncodeToString(publicKeyBytes)

	var successList []string
	signatureData.FileSignatures, successList, err = filesignature.SignFileHashes(hashSigner, resultList)
	if err != nil {
		logger.PrintErrorf(34, `could not sign file hashes: %v`, err)
		return rcProcessError
	}

	err = signatureData.Sign(hashSigner, contextId)
	if err != nil {
		logger.PrintErrorf(35, `could not sign signature file data: %v`, err)
		return rcProcessError
	}

	err = signaturefile.WriteSignatureFile(signaturesFileName, signatureData)
	if err != nil {
		logger.PrintError(36, err.Error())
		return rcProcessError
	}

	printMetaData(contextId, keyid.KeyId(publicKeyBytes), signatureData.Timestamp, signatureData.Hostname)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList(`signing`, successList)
	}

	successEnding := texthelper.GetCountEnding(successCount)

	logger.PrintInfof(37, `signature%s for %d file%s successfully created`, successEnding, len(successList), successEnding)
	return rcOK
}
