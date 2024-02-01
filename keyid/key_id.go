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

package keyid

import (
	"crypto/subtle"
	"filesigner/base32encoding"
	"golang.org/x/crypto/sha3"
)

// ******** Private constants ********

var beginFence = []byte{'k', 'e', 'y', 0x5a}
var endFence = []byte{0xa5, 'h', 's', 'h'}

// ******** Public functions ********

// KeyHash calculates the Shake-128 hash of some key bytes.
func KeyHash(key []byte) []byte {
	hasher := sha3.NewShake128()

	_, _ = hasher.Write(beginFence)
	_, _ = hasher.Write(key)
	_, _ = hasher.Write(endFence)

	rawResult := make([]byte, 32)
	_, _ = hasher.Read(rawResult)

	result := make([]byte, 16)
	subtle.XORBytes(result, rawResult[:16], rawResult[16:])

	return result
}

// KeyId returns the key id of some key bytes.
func KeyId(key []byte) string {
	return base32encoding.EncodeKey(KeyHash(key))
}
