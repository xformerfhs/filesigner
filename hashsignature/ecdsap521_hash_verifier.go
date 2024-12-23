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
// Version: 3.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-04-05: V1.0.1: Make type private.
//    2024-04-05: V2.0.0: Correct name of type and creation function.
//    2024-12-23: V3.0.0: Do not return an error.
//

package hashsignature

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
)

// ******** Public types ********

// ecDsaP521HashVerifier contains the objects necessary for signature verification with curve sec521r1.
type ecDsaP521HashVerifier struct {
	publicKey     *ecdsa.PublicKey
	publicKeyHash []byte
}

// ******** Private constants ********

// p521PublicKeyLength is the length of a curve P-521 public key.
const p521PublicKeyLength = 158

// ******** Type creation ********

// NewEcDsaP521HashVerifier creates a new ecDsaP521HashVerifier.
func NewEcDsaP521HashVerifier(publicKey []byte) (HashVerifier, error) {
	lenKey := len(publicKey)
	if lenKey != p521PublicKeyLength {
		return nil, fmt.Errorf(`Bad ec dsa p521 public key length: %d`, lenKey)
	}

	var err error
	result := &ecDsaP521HashVerifier{}

	var pk any
	pk, err = x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf(`Invalid public key: %v`, err)
	}

	var ok bool
	result.publicKey, ok = pk.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New(`Public key is not an ECDSA key`)
	}

	return result, nil
}

// ******** Public functions ********

// VerifyHash verifies the supplied hash with the supplied signature.
func (hv *ecDsaP521HashVerifier) VerifyHash(hashValue []byte, signature []byte) bool {
	return ecdsa.VerifyASN1(hv.publicKey, hashValue, signature)
}
