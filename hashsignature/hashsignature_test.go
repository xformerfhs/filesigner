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
//    2024-04-06: V1.0.0: Created.
//

package hashsignature

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	mrand "math/rand"
	"strings"
	"testing"
)

// ******** Private constants ********

// testLoopCount is the loop count for verification tests.
const testLoopCount = 100

// ******** Private types ********

// testEnvironment contains the data needed for tests.
type testEnvironment struct {
	algorithmName string
	signer        HashSigner
	verifier      HashVerifier
	publicKey     []byte
	publicKeyLen  int
}

// ******** Private variables ********

// testEnvironments contains the test environments.
var testEnvironments []testEnvironment

// ******** Setup functions ********

func ensureEnvironment(t *testing.T) {
	if len(testEnvironments) == 0 {
		setupEnvironment(t)
	}
}

func setupEnvironment(t *testing.T) {
	testEnvironments = make([]testEnvironment, 2)

	// 1. EcDsaP521 environment.
	testEnvironments[0] = makeEcDsaP521Environment(t)

	// 2. Ed25519 environment
	testEnvironments[1] = makeEd25519Environment(t)
}

func makeEcDsaP521Environment(t *testing.T) testEnvironment {
	algorithmName := `EcDsaP521`

	signer, err := NewEcDsaP521HashSigner()

	// Check if signer creation had no error.
	if err != nil {
		fatalActionExitf(t, `creat`, algorithmName, `data signer`)
	}

	// Check if signer is not nil.
	if signer == nil {
		fatalNilExitf(t, algorithmName, `signer`)
	}

	var publicKey []byte
	publicKey, err = signer.PublicKey()

	// Check if getting public key had an error.
	if err != nil {
		fatalActionExitf(t, `gett`, algorithmName, `public key`)
	}

	var verifier HashVerifier
	verifier, err = NewEcDsaP521HashVerifier(publicKey)

	// Check if verifier creation had no error.
	if err != nil {
		fatalActionExitf(t, `creat`, algorithmName, `data verifier`)
	}

	// Check if verifier is not nil.
	if verifier == nil {
		fatalNilExitf(t, algorithmName, `verifier`)
	}

	return testEnvironment{
		algorithmName: algorithmName,
		signer:        signer,
		verifier:      verifier,
		publicKey:     publicKey,
		publicKeyLen:  p521PublicKeyLength,
	}
}

func makeEd25519Environment(t *testing.T) testEnvironment {
	algorithmName := `Ed25519`

	signer, err := NewEd25519HashSigner()
	if err != nil {
		fatalActionExitf(t, `creat`, algorithmName, `data signer`)
	}

	// Check if signer is not nil.
	if signer == nil {
		fatalNilExitf(t, algorithmName, `signer`)
	}

	var publicKey []byte
	publicKey, err = signer.PublicKey()

	// Check if getting public key had an error.
	if err != nil {
		fatalActionExitf(t, `gett`, algorithmName, `public key`)
	}

	var verifier HashVerifier
	verifier, err = NewEd25519HashVerifier(publicKey)

	// Check if verifier creation had no error.
	if err != nil {
		fatalActionExitf(t, `creat`, algorithmName, `data verifier`)
	}

	// Check if verifier is not nil.
	if verifier == nil {
		fatalNilExitf(t, algorithmName, `verifier`)
	}

	return testEnvironment{
		algorithmName: algorithmName,
		signer:        signer,
		verifier:      verifier,
		publicKey:     publicKey,
		publicKeyLen:  ed25519.PublicKeySize,
	}
}

func fatalActionExitf(t *testing.T, actionName string, algorithmName string, objectName string) {
	t.Fatalf(`Error %sing %s %s`, actionName, algorithmName, objectName)
}

func fatalNilExitf(t *testing.T, algorithmName string, objectName string) {
	t.Fatalf(`New %s returned nil %s`, algorithmName, objectName)
}

// ******** Tests ********

func TestPublicKeyLen(t *testing.T) {
	ensureEnvironment(t)

	for _, env := range testEnvironments {
		doTestPublicKeyLen(t, env.algorithmName, env.publicKey, env.publicKeyLen)
	}
}

func TestWrongPublicKeyWithVerifyEcDsaP521(t *testing.T) {
	wrongPublicKey := make([]byte, p521PublicKeyLength)
	_, err := NewEcDsaP521HashVerifier(wrongPublicKey)
	algorithmName := `EcDsaP521`
	if err == nil || !strings.Contains(err.Error(), `public key`) {
		t.Fatalf(`%s verifier has no error with invalid public key`, algorithmName)
	}

	var privKey *rsa.PrivateKey
	privKey, err = rsa.GenerateKey(rand.Reader, 1000)
	if err != nil {
		t.Fatalf(`Error generating RSA key: %v`, err)
	}

	pubKey := &privKey.PublicKey
	var pubKeyBytes []byte
	pubKeyBytes, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		t.Fatalf(`Error converting RSA key: %v`, err)
	}

	_, err = NewEcDsaP521HashVerifier(pubKeyBytes)
	if err == nil || !strings.Contains(err.Error(), `not an ECDSA key`) {
		t.Fatalf(`Error generating RSA key: %v`, err)
	}
}

func TestShortPublicKeyWithVerify(t *testing.T) {
	wrongPublicKey := make([]byte, 3)
	_, err := NewEcDsaP521HashVerifier(wrongPublicKey)
	algorithmName := `EcDsaP521`
	if err == nil || !strings.Contains(err.Error(), `key length`) {
		t.Fatalf(`%s verifier has no error with wrong public key`, algorithmName)
	}

	_, err = NewEd25519HashVerifier(wrongPublicKey)
	algorithmName = `Ed25519`
	if err == nil || !strings.Contains(err.Error(), `key length`) {
		t.Fatalf(`%s verifier has no error with wrong public key`, algorithmName)
	}
}

func TestSignAndVerify(t *testing.T) {
	ensureEnvironment(t)

	for _, env := range testEnvironments {
		doTestSignAndVerify(t, env.algorithmName, env.signer, env.verifier, false)
	}
}

func TestSignAndVerifyInvalid(t *testing.T) {
	ensureEnvironment(t)

	for _, env := range testEnvironments {
		doTestSignAndVerify(t, env.algorithmName, env.signer, env.verifier, true)
	}
}

func TestSignerDestroy(t *testing.T) {
	// This must be the last test. All other tests will fail, after this one.
	ensureEnvironment(t)

	for _, env := range testEnvironments {
		doTestSignerDestroy(t, env.algorithmName, env.signer)
	}
}

// ******** Tests for one type ********

// These tests were written with the help of ChatGPT.

func doTestPublicKeyLen(t *testing.T, algorithmName string, publicKey []byte, publicKeyLength int) {
	if len(publicKey) != publicKeyLength {
		t.Fatalf(`%s: Invalid public key length: expected %d, got %d`, algorithmName, publicKeyLength, len(publicKey))
	}
}

func doTestSignAndVerify(
	t *testing.T,
	algorithmName string,
	signer HashSigner,
	verifier HashVerifier,
	testInvalid bool,
) {
	var validWord string
	if testInvalid {
		validWord = `Inv`
	} else {
		validWord = `V`
	}

	for i := 0; i < testLoopCount; i++ {
		// Generate some data to sign.
		dataLen := mrand.Intn(100) + 4

		data := make([]byte, dataLen)
		_, _ = rand.Read(data)

		// Sign the data.
		var err error
		var signature []byte
		signature, err = signer.SignHash(data)

		if err != nil {
			t.Fatalf(`Error signing hash with %s: %v`, algorithmName, err)
		}

		// Check if signature is not nil.
		if signature == nil {
			t.Fatalf(`Signing hash with %s returned nil signature`, algorithmName)
		}

		// Invalidate signature if invalid signature is to be tested.
		if testInvalid {
			signature[len(signature)>>1] ^= 0xff
		}

		// Verify signature.
		var verifyResult bool
		verifyResult, err = verifier.VerifyHash(data, signature)

		if err != nil {
			t.Fatalf(`Error verifying hash signature with %s: %v`, algorithmName, err)
		}

		// Test verification result.
		if verifyResult == testInvalid {
			t.Fatalf(`%salid hash signature verification with %s returned %t`, validWord, algorithmName, verifyResult)
		}
	}
}

func doTestSignerDestroy(t *testing.T, algorithmName string, signer HashSigner) {
	// Destroy the signer.
	signer.Destroy()

	// Attempt to retrieve public key after destruction.
	_, err := signer.PublicKey()
	if err == nil {
		t.Fatalf(`No error when retrieving %s public key after destruction`, algorithmName)
	}

	// Generate some random has value.
	data := make([]byte, 16)
	_, _ = rand.Read(data)

	// Attempt to sign data after destruction.
	_, err = signer.SignHash(data)
	if err == nil {
		t.Fatalf(`No error when signing hash with %s after destruction`, algorithmName)
	}
}
