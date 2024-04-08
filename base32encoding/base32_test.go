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
// Version: 1.1.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-03-23: V1.1.0: Added tests for all functions.
//

package base32encoding

import (
	"bytes"
	cryptorand "crypto/rand"
	"math/rand"
	"strings"
	"testing"
)

// ******** Private constants ********

// testLoopCount contains the no. of times a test is repeated.
const testLoopCount = 1_000

func TestEncodeKey(t *testing.T) {
	for i := 0; i < testLoopCount; i++ {
		sl := rand.Intn(30)
		s := make([]byte, sl)
		_, _ = cryptorand.Read(s)
		es := EncodeKey(s)
		ds, err := DecodeKey(es)
		if err != nil {
			t.Fatalf(`error decoding key '%s': %v`, es, err)
		}

		if !bytes.Equal(s, ds) {
			t.Fatalf(`decoding '%s' did not result in '%x', but '%x'`, es, s, ds)
		}
	}
}

func TestDecodeKeyInvalidCharacter(t *testing.T) {
	k := `ABCD-ABCD-AB`
	_, err := DecodeKey(k)
	if err != nil {
		if !strings.Contains(err.Error(), `illegal base32`) {
			t.Fatalf(`Wrong error: %v`, err)
		}
	} else {
		t.Fatal(`No error with invalid character`)
	}
}

func TestDecodeKeyInvalidGroupSize(t *testing.T) {
	keys := []string{`BCDF-B-AB`, `BCDF-BC-AB`, `BCDF-BCD-AB`, `BCDF-BCDFG-AB`}
	for _, k := range keys {
		_, err := DecodeKey(k)
		if err != nil {
			if !strings.Contains(err.Error(), `group size`) {
				t.Fatalf(`Wrong error with key '%s': %v`, k, err)
			}
		} else {
			t.Fatal(`No error with wrong group size`)
		}
	}
}

func TestEncodeToString(t *testing.T) {
	for i := 0; i < testLoopCount; i++ {
		sl := rand.Intn(30)
		s := make([]byte, sl)
		_, _ = cryptorand.Read(s)
		es := EncodeToString(s)
		ds, err := DecodeFromString(es)
		if err != nil {
			t.Fatalf(`error decoding key '%s': %v`, es, err)
		}

		if !bytes.Equal(s, ds) {
			t.Fatalf(`decoding '%s' did not result in '%02x', but '%02x'`, es, s, ds)
		}
	}
}

func TestDecodeInvalidCharacter(t *testing.T) {
	k := `123456`
	_, err := DecodeFromString(k)
	if err != nil {
		if !strings.Contains(err.Error(), `illegal base32`) {
			t.Fatalf(`Wrong error: %v`, err)
		}
	} else {
		t.Fatal(`No error with invalid character`)
	}
}

func TestEncodeToBytes(t *testing.T) {
	for i := 0; i < testLoopCount; i++ {
		sl := rand.Intn(30)
		s := make([]byte, sl)
		_, _ = cryptorand.Read(s)
		es := EncodeToBytes(s)
		ds, err := DecodeFromBytes(es)
		if err != nil {
			t.Fatalf(`error decoding bytes '%02x': %v`, s, err)
		}

		if !bytes.Equal(s, ds) {
			t.Fatalf(`decoding '%02x' did not result in '%02x', but '%02x'`, es, s, ds)
		}
	}
}

func BenchmarkEncodeKey(b *testing.B) {
	source := make([]byte, 16)
	_, _ = cryptorand.Read(source)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = EncodeKey(source)
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	source := make([]byte, 16)
	_, _ = cryptorand.Read(source)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = EncodeToString(source)
	}
}

func BenchmarkEncodeToBytes(b *testing.B) {
	source := make([]byte, 16)
	_, _ = cryptorand.Read(source)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = EncodeToBytes(source)
	}
}
