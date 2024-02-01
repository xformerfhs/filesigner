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

package signaturefile

import (
	"encoding/json"
	"errors"
	"filesigner/filehelper"
	"filesigner/signaturehandler"
	"fmt"
	"os"
)

// ******** Private constants ********

const maxFileSize = 10_000_000

// ******** Public functions ********

func WriteSignatureFile(filePath string, signatureData *signaturehandler.SignatureData) error {
	jsonOutput, err := json.MarshalIndent(signatureData, "", "   ")
	if err != nil {
		return fmt.Errorf("Could not create json format for signature file: %w", err)
	}

	err = os.WriteFile(filePath, jsonOutput, 0600)
	if err != nil {
		return fmt.Errorf("Could not write signature file: %w", err)
	}

	return nil
}

// ReadSignatureFile reads a signature file and returns the signature data.
func ReadSignatureFile(filePath string) (*signaturehandler.SignatureData, error) {
	err := checkFileSize(filePath)
	if err != nil {
		return nil, err
	}

	var result *signaturehandler.SignatureData
	result, err = getSignatureData(filePath)
	if err != nil {
		return nil, err
	}

	err = checkSignatureForm(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Private functions ********

// checkFileSize checks if the size of the signature file is within the allowed boundaries.
func checkFileSize(filePath string) error {
	fileSize, err := filehelper.FileSize(filePath)
	if err != nil {
		return err
	}

	if fileSize > maxFileSize {
		return errors.New("Signature file is too large.")
	}

	return nil
}

// getSignatureData reads the signature data from a file and returns the data in a SignatureData structure.
func getSignatureData(filePath string) (*signaturehandler.SignatureData, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	result := &signaturehandler.SignatureData{
		Format:        signaturehandler.SignatureFormatInvalid,
		SignatureType: signaturehandler.SignatureTypeInvalid,
	}
	err = json.Unmarshal(fileContent, result)

	return result, err
}

// checkSignatureForm checks if the signature data are complete and satisfy some formal requirements.
func checkSignatureForm(signatureData *signaturehandler.SignatureData) error {
	err := checkMissingInformation(signatureData)
	if err != nil {
		return err
	}

	if signatureData.Format > signaturehandler.SignatureFormatMax {
		return fmt.Errorf("Invalid signature format id: %d", signatureData.Format)
	}

	if signatureData.SignatureType > signaturehandler.SignatureTypeMax {
		return fmt.Errorf("Invalid signature type: %d", signatureData.SignatureType)
	}

	return nil
}

// checkMissingInformation checks if any required signature result data is missing.
func checkMissingInformation(signatureData *signaturehandler.SignatureData) error {
	if len(signatureData.DataSignature) == 0 {
		return makeMissingFieldError("DataSignature")
	}
	if signatureData.FileSignatures == nil {
		return makeMissingFieldError("FileSignatures")
	}
	if signatureData.Format == signaturehandler.SignatureFormatInvalid {
		return makeMissingFieldError("Format")
	}
	if len(signatureData.Hostname) == 0 {
		return makeMissingFieldError("Hostname")
	}
	if len(signatureData.PublicKey) == 0 {
		return makeMissingFieldError("PublicKey")
	}
	if signatureData.SignatureType == signaturehandler.SignatureTypeInvalid {
		return makeMissingFieldError("SignatureType")
	}
	if len(signatureData.Timestamp) == 0 {
		return makeMissingFieldError("Timestamp")
	}
	return nil
}

// makeMissingFieldError build the error for a "missing field" error type.
func makeMissingFieldError(fieldName string) error {
	return fmt.Errorf("Field '%s' is missing from signature file", fieldName)
}
