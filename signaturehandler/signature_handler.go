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

package signaturehandler

import (
	"filesigner/base32encoding"
	"filesigner/contexthasher"
	"filesigner/hashsignature"
	"filesigner/maphelper"
	"filesigner/numberhelper"
	"filesigner/stringhelper"
	"golang.org/x/crypto/sha3"
	"hash"
)

// ******** Public types ********

// SignatureFormat contains the the format id of the signature file.
type SignatureFormat byte

// SignatureType contains the the code for the signature algorithm.
type SignatureType byte

// SignatureData contains all the data that comprise a filesigner signature.
type SignatureData struct {
	Format         SignatureFormat   `json:"format"`
	PublicKey      string            `json:"publicKey"`
	Timestamp      string            `json:"timeStamp"`
	Hostname       string            `json:"hostname"`
	SignatureType  SignatureType     `json:"signatureType"`
	FileSignatures map[string]string `json:"fileSignatures"`
	DataSignature  string            `json:"dataSignature"`
}

// ******** Public constants ********

// These are the possible values for SignatureFormat.
const (
	SignatureFormatInvalid SignatureFormat = iota
	SignatureFormatV1
	SignatureFormatMax = iota - 1
)

// These are the possible values for SignatureType.
const (
	SignatureTypeInvalid SignatureType = iota
	SignatureTypeEd25519
	SignatureTypeEcDsaP521
	SignatureTypeMax = iota - 1
)

// ******** Public type functions ********

// Sign adds the data signature to a SignatureData
func (sd *SignatureData) Sign(hashSigner hashsignature.HashSigner, contextId string) error {
	hashValue := getHashValueOfSignatureData(sd, stringhelper.UnsafeStringBytes(contextId))
	signatureValue, err := hashSigner.SignHash(hashValue)
	if err != nil {
		return err
	}

	sd.DataSignature = base32encoding.EncodeToString(signatureValue)

	return nil
}

// Verify verifies the data signature of a SignatureData
func (sd *SignatureData) Verify(hashVerifier hashsignature.HashVerifier, contextId string) (bool, error) {
	dataSignature, err := base32encoding.DecodeFromString(sd.DataSignature)
	if err != nil {
		return false, err
	}

	hashValue := getHashValueOfSignatureData(sd, stringhelper.UnsafeStringBytes(contextId))
	return hashVerifier.VerifyHash(hashValue, dataSignature)
}

// ******** Private functions ********

// getHashValueOfSignatureData calculates the hash value of a SignatureData
func getHashValueOfSignatureData(signatureData *SignatureData, contextBytes []byte) []byte {
	hasher := contexthasher.NewContextHasher(sha3.New512(), contextBytes)

	position, _ := numberhelper.NewByteCounterForCount(uint64((len(signatureData.FileSignatures) << 1) + 5))
	oneByteSlice := make([]byte, 1)

	hashPosition(hasher, position)
	oneByteSlice[0] = byte(signatureData.Format)
	hasher.Write(oneByteSlice)

	hashPosition(hasher, position)
	hasher.Write(stringhelper.UnsafeStringBytes(signatureData.PublicKey))

	hashPosition(hasher, position)
	hasher.Write(stringhelper.UnsafeStringBytes(signatureData.Timestamp))

	hashPosition(hasher, position)
	hasher.Write(stringhelper.UnsafeStringBytes(signatureData.Hostname))

	hashPosition(hasher, position)
	oneByteSlice[0] = byte(signatureData.SignatureType)
	hasher.Write(oneByteSlice)

	sortedFileNames := maphelper.SortedKeys(signatureData.FileSignatures)
	for _, fileName := range sortedFileNames {
		hashPosition(hasher, position)
		hasher.Write(stringhelper.UnsafeStringBytes(fileName))

		hashPosition(hasher, position)
		hasher.Write(stringhelper.UnsafeStringBytes(signatureData.FileSignatures[fileName]))
	}

	return hasher.Sum(nil)
}

// hashPosition writes the position into the hasher
func hashPosition(hasher hash.Hash, position *numberhelper.ByteCounter) {
	position.Inc()
	hasher.Write(position.Slice())
}
