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
// Version: 1.3.1
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-25: V1.1.0: Use "Ed25519", instead of "Ed25519ph".
//    2024-02-26: V1.2.0: Use a strengthened version of "Ed25519ph".
//    2024-02-26: V1.3.0: Use a strengthened version of "Ed25519".
//    2024-04-05: V1.3.1: Make type private, add validity check for PublicKey.
//

package hashsignature

import (
	"crypto/ed25519"
	"filesigner/slicehelper"
)

// ******** Private types ********

// ed25519HashSigner contains the objects necessary for ed25519 hash signing.
type ed25519HashSigner struct {
	signer    ed25519.PrivateKey
	publicKey []byte
	isValid   bool
}

// ******** Type creation ********

// NewEd25519HashSigner creates a new ed25519HashSigner.
func NewEd25519HashSigner() (HashSigner, error) {
	var err error

	result := &ed25519HashSigner{
		isValid: true,
	}

	result.publicKey, result.signer, err = ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Public functions ********

// PublicKey returns a copy of the public key.
func (hs *ed25519HashSigner) PublicKey() ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return slicehelper.Copy(hs.publicKey), nil
}

// SignHash signs the supplied hash value.
func (hs *ed25519HashSigner) SignHash(hashValue []byte) ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	// Ed25519 does its own hashing and expects the full source as a parameter.
	// This is not possible here, as files can be arbitrarily large and Ed25519 can not have
	// an io.Reader interface. So we can only supply the already computed hash.
	// However, this hash only has 64 bytes, which is too short to provide enough security.
	// So we pad the hash value with a constant padding. This is similar to the HashEdDSA
	// variant ed25519ph of RFC8032. This RFC uses a constant prefix. Here we use a constant
	// prefix and suffix.
	return ed25519.Sign(hs.signer, paddedHash(hashValue)), nil
}

// Destroy removes the private key from this ed25519HashSigner, so it can no longer be used.
func (hs *ed25519HashSigner) Destroy() {
	if hs.isValid {
		slicehelper.ClearNumber(hs.signer)
		hs.signer = nil
		hs.isValid = false
	}
}

// ******** Private functions ********

// checkValidity checks if this ed25519HashSigner is usable.
func (hs *ed25519HashSigner) checkValidity() error {
	if hs.isValid {
		return nil
	} else {
		return IsDestroyedErr
	}
}
