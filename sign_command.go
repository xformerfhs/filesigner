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
// Version: 1.3.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-17: V1.1.0: Create contextBytes as early, as possible.
//    2024-03-04: V1.2.0: Use public key bytes, not id.
//    2025-03-01: V1.3.0: Add message base.
//

package main

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/filesignature"
	"filesigner/hashsignature"
	"filesigner/logger"
	"filesigner/signaturefile"
	"filesigner/signaturehandler"
	"filesigner/stretcher"
	"filesigner/stringhelper"
	"filesigner/texthelper"
	"os"
	"time"
)

// ******** Private constants ********

// signCmdMsgBase is the base number for all messages in sign_command.
// This file reserves numbers 30-49.
const signCmdMsgBase = 30

// timeStampFormat RFC3339 format for signatures file time stamp
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
		ContextId:     contextId,
	}

	signatureData.Hostname, err = os.Hostname()
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+1, `Could not get host name: %v`, err)
		return rcProcessError
	}

	contextKey := stretcher.KeyFromBytes(stringhelper.UnsafeStringBytes(contextId))
	resultList := filehasher.FileHashes(filePaths, contextKey)

	if existHashErrors(resultList) {
		return rcProcessError
	}

	var hashSigner hashsignature.HashSigner
	if signatureType == signaturehandler.SignatureTypeEd25519 {
		hashSigner, err = hashsignature.NewEd25519HashSigner()
	} else {
		hashSigner, err = hashsignature.NewEcDsaP521HashSigner()
	}
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+2, `Could not create hash-signer: %v`, err)
		return rcProcessError
	}
	defer hashSigner.Destroy()

	var publicKeyBytes []byte
	publicKeyBytes, err = hashSigner.PublicKey()
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+3, `Could not get public key bytes: %v`, err)
		return rcProcessError
	}
	signatureData.PublicKey = base32encoding.EncodeToString(publicKeyBytes)

	var successList []string
	signatureData.FileSignatures, successList, err = filesignature.SignFileHashes(hashSigner, resultList)
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+4, `Could not sign file hashes: %v`, err)
		return rcProcessError
	}

	err = signatureData.Sign(hashSigner, contextKey)
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+5, `Could not sign signatures file data: %v`, err)
		return rcProcessError
	}

	err = signaturefile.WriteJson(signaturesFileName, signatureData)
	if err != nil {
		logger.PrintErrorf(signCmdMsgBase+6, `Error writing signatures file '%s': %v`, signaturesFileName, err)
		return rcProcessError
	}

	printMetaData(signatureData, publicKeyBytes)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList(`Signing`, successList)
	}

	successEnding := texthelper.GetCountEnding(successCount)

	logger.PrintInfof(signCmdMsgBase+7,
		`Signature%s for %d file%s successfully created and written to '%s'`,
		successEnding,
		len(successList),
		successEnding,
		signaturesFileName)
	return rcOK
}
