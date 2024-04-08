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
//    2024-04-05: V1.0.1: Make type private, add validity check for PublicKey.
//    2024-04-05: V2.0.0: Correct name of type and creation function.
//

package hashsignature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
)

// ******** Private types ********

// ecDsaP521HashSigner contains the objects necessary for signing a hash with curve secp521r1.
type ecDsaP521HashSigner struct {
	privateKey *ecdsa.PrivateKey
	isValid    bool
}

// ******** Type creation ********

// NewEcDsaP521HashSigner creates a new ecDsaP521HashSigner.
func NewEcDsaP521HashSigner() (HashSigner, error) {
	var err error

	result := &ecDsaP521HashSigner{
		isValid: true,
	}

	result.privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Public functions ********

// PublicKey returns a copy of the public key.
func (hs *ecDsaP521HashSigner) PublicKey() ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return x509.MarshalPKIXPublicKey(hs.privateKey.Public())
}

// SignHash signs the supplied hash value.
func (hs *ecDsaP521HashSigner) SignHash(hashValue []byte) ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return ecdsa.SignASN1(rand.Reader, hs.privateKey, hashValue)
}

// Destroy overwrites the private key, so the signer is no longer usable.
func (hs *ecDsaP521HashSigner) Destroy() {
	if hs.isValid {
		hs.privateKey.D.SetInt64(-1)
		hs.isValid = false
	}
}

// ******** Private functions ********

// checkValidity checks if this ecDsaP521HashSigner is usable.
func (hs *ecDsaP521HashSigner) checkValidity() error {
	if hs.isValid {
		return nil
	} else {
		return IsDestroyedErr
	}
}
