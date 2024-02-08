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

package hashsignature

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
)

// ******** Public types ********

// Ec521HashVerifier contains the objects necessary for file curve sec521r1 and hash verification.
type Ec521HashVerifier struct {
	publicKey     *ecdsa.PublicKey
	publicKeyHash []byte
}

// ******** Private constants ********

// p521KeyLength is the length of a curve P-521 key.
const p521KeyLength = 158

// ******** Type creation ********

// NewEc521HashVerifier creates a new Ec521HashVerifier.
func NewEc521HashVerifier(publicKey []byte) (HashVerifier, error) {
	lenKey := len(publicKey)
	if lenKey != p521KeyLength {
		return nil, fmt.Errorf(`bad ec dsa p521 public key length: %d`, lenKey)
	}

	var err error
	result := &Ec521HashVerifier{}

	var pk any
	pk, err = x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf(`invalid public key: %v`, err)
	}

	var ok bool
	result.publicKey, ok = pk.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New(`public key is not an ECDSA key`)
	}

	return result, nil
}

// ******** Public functions ********

// VerifyHash verifies the supplied hash with the supplied signature.
func (hv *Ec521HashVerifier) VerifyHash(hashValue []byte, signature []byte) (bool, error) {
	return hv.doVerifyHash(hashValue, signature)
}

// ******** Private functions ********

// doVerifyHash verifies a supplied hash value with a supplied signature.
func (hv *Ec521HashVerifier) doVerifyHash(hashValue []byte, signature []byte) (bool, error) {
	return ecdsa.VerifyASN1(hv.publicKey, hashValue, signature), nil
}
