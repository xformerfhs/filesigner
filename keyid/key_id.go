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
//    2025-05-23: V2.0.0: Functions can be called with multiple byte slices.
//

package keyid

import (
	"crypto/subtle"
	"filesigner/base32encoding"
	"golang.org/x/crypto/sha3"
)

// ******** Private constants ********

// beginFence is the first bytes of the fence.
var beginFence = []byte{'k', 'e', 'y', 0x5a}

// endFence is the last bytes of the fence.
var endFence = []byte{0xa5, 'h', 's', 'h'}

// ******** Private variables ********

// hasher is the Shake-128 hasher.
var hasher = sha3.NewShake128()

// singleByte is a byte slice with a single element used for hashing integers.
var singleByte = make([]byte, 1)

// ******** Public functions ********

// KeyHash calculates the Shake-128 hash of a slice of byte slices.
func KeyHash(s ...[]byte) []byte {
	hasher.Reset()

	_, _ = hasher.Write(beginFence)
	for i, b := range s {
		// Hash length of byte slice.
		singleByte[0] = byte(len(b))
		_, _ = hasher.Write(singleByte)

		// Hash content of byte slice.
		_, _ = hasher.Write(b)

		// Hash position.
		singleByte[0] = byte(i)
		_, _ = hasher.Write(singleByte)
	}
	_, _ = hasher.Write(endFence)

	// Get hash result.
	rawResult := make([]byte, 32)
	_, _ = hasher.Read(rawResult)

	// Xor upper and lower half of hash result as the final result.
	result := make([]byte, 16)
	subtle.XORBytes(result, rawResult[:16], rawResult[16:])

	return result
}

// KeyId returns the key id of some key bytes.
func KeyId(key ...[]byte) string {
	return base32encoding.EncodeKey(KeyHash(key...))
}
