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

package hashsignature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
)

// ******** Public types ********

// EcDsa521HashSigner contains the objects necessary for file signing with curve secp521r1 and hash.
type EcDsa521HashSigner struct {
	privateKey *ecdsa.PrivateKey
	isValid    bool
}

// ******** Type creation ********

// NewEcDsa521HashSigner creates a new EcDsa521HashSigner.
func NewEcDsa521HashSigner() (HashSigner, error) {
	var err error

	result := &EcDsa521HashSigner{
		isValid: true,
	}

	result.privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Public functions ********

// GetPublicKey returns a copy of the public key
func (hs *EcDsa521HashSigner) GetPublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(hs.privateKey.Public())
}

// -------- Sign functions --------

// SignHash signs the supplied has value.
func (hs *EcDsa521HashSigner) SignHash(hashValue []byte) ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return hs.doSignHash(hashValue)
}

func (hs *EcDsa521HashSigner) Destroy() {
	if hs.isValid {
		hs.privateKey.D.SetInt64(-1)
		hs.isValid = false
	}
}

// ******** Private functions ********

// checkValidity checks if this Ed25519HashSigner is usable.
func (hs *EcDsa521HashSigner) checkValidity() error {
	if hs.isValid {
		return nil
	} else {
		return IsDestroyedErr
	}
}

// doSignHash signs a supplied hash value.
func (hs *EcDsa521HashSigner) doSignHash(hashValue []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, hs.privateKey, hashValue)
}
