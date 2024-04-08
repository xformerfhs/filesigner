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

package signaturefile

import (
	"bytes"
	"encoding/json"
	"errors"
	"filesigner/filehelper"
	"filesigner/signaturehandler"
	"fmt"
	"os"
)

// ******** Private constants ********

// maxFileSize is the maximum allowed size of a signatures file.
const maxFileSize = 50_000_000

// ******** Public functions ********

// WriteJson writes the signature data to the specified file in JSON format.
func WriteJson(filePath string, signatureData *signaturehandler.SignatureData) error {
	jsonOutput, err := json.MarshalIndent(signatureData, "", "   ")
	if err != nil {
		return fmt.Errorf(`Could not convert data to JSON format: %w`, err)
	}

	err = os.WriteFile(filePath, jsonOutput, 0600)
	if err != nil {
		return fmt.Errorf(`Could not write signatures file: %w`, err)
	}

	return nil
}

// ReadJson reads a signatures file in JSON format and returns the signature data.
func ReadJson(filePath string) (*signaturehandler.SignatureData, error) {
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

// checkFileSize checks if the size of the signatures file is within the allowed boundaries.
func checkFileSize(filePath string) error {
	fileSize, err := filehelper.FileSize(filePath)
	if err != nil {
		return err
	}

	if fileSize > maxFileSize {
		return errors.New(`Signatures file is too large`)
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
	err = strictJsonUnmarshal(fileContent, result)

	return result, err
}

// strictJsonUnmarshal implements a json.Unmarshal which throws an error when encountering unknown fields.
// json.Unmarshal ignores unknown fields.
func strictJsonUnmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// checkSignatureForm checks if the signature data are complete and satisfy some formal requirements.
func checkSignatureForm(signatureData *signaturehandler.SignatureData) error {
	err := checkMissingInformation(signatureData)
	if err != nil {
		return err
	}

	if signatureData.Format > signaturehandler.SignatureFormatMax {
		return fmt.Errorf(`Invalid signature format id: %d`, signatureData.Format)
	}

	if signatureData.SignatureType > signaturehandler.SignatureTypeMax {
		return fmt.Errorf(`Invalid signature type: %d`, signatureData.SignatureType)
	}

	return nil
}

// checkMissingInformation checks if any required signature result data is missing.
func checkMissingInformation(signatureData *signaturehandler.SignatureData) error {
	if len(signatureData.DataSignature) == 0 {
		return makeMissingFieldError(`dataSignature`)
	}
	if len(signatureData.ContextId) == 0 {
		return makeMissingFieldError(`contextId`)
	}
	if signatureData.FileSignatures == nil {
		return makeMissingFieldError(`fileSignatures`)
	}
	if signatureData.Format == signaturehandler.SignatureFormatInvalid {
		return makeMissingFieldError(`format`)
	}
	if len(signatureData.Hostname) == 0 {
		return makeMissingFieldError(`hostname`)
	}
	if len(signatureData.PublicKey) == 0 {
		return makeMissingFieldError(`publicKey`)
	}
	if signatureData.SignatureType == signaturehandler.SignatureTypeInvalid {
		return makeMissingFieldError(`signatureType`)
	}
	if len(signatureData.Timestamp) == 0 {
		return makeMissingFieldError(`timestamp`)
	}
	return nil
}

// makeMissingFieldError build the error for a "missing field" error type.
func makeMissingFieldError(fieldName string) error {
	return fmt.Errorf(`Field '%s' is missing from signatures file`, fieldName)
}
