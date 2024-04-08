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
// Version: 2.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-17: V1.1.0: Use contextBytes.
//    2024-02-25: V2.0.0: Rename "Ed25519" to "Ed25519Ph".
//

package signaturehandler

import (
	"filesigner/base32encoding"
	"filesigner/hashsignature"
	"filesigner/maphelper"
	"filesigner/numberhelper"
	"filesigner/paddedhasher"
	"filesigner/stringhelper"
	"golang.org/x/crypto/sha3"
	"hash"
)

// ******** Public types ********

// SignatureFormat contains the format id of the signatures file.
type SignatureFormat byte

// SignatureType contains the code for the signature algorithm.
type SignatureType byte

// SignatureData contains all the data that comprise a filesigner signature.
type SignatureData struct {
	Format         SignatureFormat   `json:"format"`
	ContextId      string            `json:"contextId"`
	PublicKey      string            `json:"publicKey"`
	Timestamp      string            `json:"timestamp"`
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

// Sign adds the data signature to a SignatureData.
func (sd *SignatureData) Sign(hashSigner hashsignature.HashSigner, contextKey []byte) error {
	hashValue := hashValueOfSignatureData(sd, contextKey)
	signatureValue, err := hashSigner.SignHash(hashValue)
	if err != nil {
		return err
	}

	sd.DataSignature = base32encoding.EncodeToString(signatureValue)

	return nil
}

// Verify verifies the data signature of a SignatureData.
func (sd *SignatureData) Verify(hashVerifier hashsignature.HashVerifier, contextKey []byte) (bool, error) {
	dataSignature, err := base32encoding.DecodeFromString(sd.DataSignature)
	if err != nil {
		return false, err
	}

	hashValue := hashValueOfSignatureData(sd, contextKey)
	return hashVerifier.VerifyHash(hashValue, dataSignature)
}

// ******** Private functions ********

// hashValueOfSignatureData calculates the hash value of a SignatureData.
func hashValueOfSignatureData(signatureData *SignatureData, contextKey []byte) []byte {
	hasher := paddedhasher.NewPaddedHasher(sha3.New512(), contextKey)

	oneByteSlice := make([]byte, 1)

	position := uint32(0)

	oneByteSlice[0] = byte(signatureData.Format)
	position = hashBytesWithPosition(hasher, position, oneByteSlice)

	position = hashStringWithPosition(hasher, position, signatureData.ContextId)
	position = hashStringWithPosition(hasher, position, signatureData.PublicKey)
	position = hashStringWithPosition(hasher, position, signatureData.Timestamp)
	position = hashStringWithPosition(hasher, position, signatureData.Hostname)

	oneByteSlice[0] = byte(signatureData.SignatureType)
	position = hashBytesWithPosition(hasher, position, oneByteSlice)

	sortedFileNames := maphelper.SortedKeys(signatureData.FileSignatures)
	for _, fileName := range sortedFileNames {
		position = hashStringWithPosition(hasher, position, fileName)
		position = hashStringWithPosition(hasher, position, signatureData.FileSignatures[fileName])
	}

	return hasher.Sum(nil)
}

// hashStringWithPosition hashes a string with the given position.
func hashStringWithPosition(hasher hash.Hash, position uint32, text string) uint32 {
	return hashBytesWithPosition(hasher, position, stringhelper.UnsafeStringBytes(text))
}

// hashBytesWithPosition hashes a byte slice with the given position.
func hashBytesWithPosition(hasher hash.Hash, position uint32, b []byte) uint32 {
	position = hashPosition(hasher, position)
	hasher.Write(b)
	hasher.Write(numberhelper.StaticIntAsShortestBigEndianBytes(len(b)))
	return position
}

// hashPosition hashes the position.
func hashPosition(hasher hash.Hash, position uint32) uint32 {
	position++
	// It is safe to use StaticUint32AsShortestBigEndianBytes here, as the returned byte slice
	// is immediately used and there is no concurrency.
	hasher.Write(numberhelper.StaticUint32AsShortestBigEndianBytes(position))
	return position
}
