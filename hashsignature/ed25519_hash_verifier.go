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
//    2024-02-25: V1.1.0: Use "Ed25519", instead of "Ed25519ph".
//    2024-02-26: V1.2.0: Use a strengthened version of "Ed25519ph".
//    2024-02-26: V1.3.0: Use a strengthened version of "Ed25519".
//    2024-04-05: V1.3.1: Make type private.
//    2024-12-23: V2.0.0: Do not return an error.
//

package hashsignature

import (
	"crypto/ed25519"
	"fmt"
)

// ******** Private types ********

// ed25519HashVerifier contains the objects necessary for ed25519 signature verification.
type ed25519HashVerifier struct {
	publicKey []byte
	options   *ed25519.Options
}

// ******** Private constants ********

// ed25519MinKeyLength is the minimum valid key length for an Ed25519 key.
const ed25519MinKeyLength = 32

// ed25519MaxKeyLength is the maximum valid key length for an Ed25519 key.
const ed25519MaxKeyLength = 34

// ******** Type creation ********

// NewEd25519HashVerifier creates a new ed25519HashVerifier.
func NewEd25519HashVerifier(publicKey []byte) (HashVerifier, error) {
	lenKey := len(publicKey)
	if lenKey < ed25519MinKeyLength || lenKey > ed25519MaxKeyLength {
		return nil, fmt.Errorf(`Bad ed25519 public key length: %d`, lenKey)
	}

	result := &ed25519HashVerifier{
		publicKey: publicKey,
	}

	return result, nil
}

// ******** Public functions ********

// VerifyHash verifies the supplied hash with the supplied signature.
func (hv *ed25519HashVerifier) VerifyHash(hashValue []byte, signature []byte) bool {
	// Ed25519 does its own hashing and expects the full source as a parameter.
	// This is not possible here, as files can be arbitrarily large and Ed25519 can not have
	// an io.Reader interface. So we can only supply the already computed hash.
	// However, this hash only has 64 bytes, which is too short to provide enough security.
	// So we pad the hash value with a constant padding. This is similar to the HashEdDSA
	// variant ed25519ph of RFC8032. This RFC uses a constant prefix. Here we use a constant
	// prefix and suffix.
	return ed25519.Verify(hv.publicKey, paddedHash(hashValue), signature)
}
