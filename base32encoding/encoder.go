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

package base32encoding

// ******** Public functions ********

// EncodeToString encodes a byte slice as a string.
func EncodeToString(b []byte) string {
	emx.Lock()
	defer emx.Unlock()

	return enc.EncodeToString(b)
}

// EncodeToBytes encodes a byte slice as a byte slice.
func EncodeToBytes(b []byte) []byte {
	emx.Lock()
	defer emx.Unlock()

	result := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(result, b)

	return result
}

// EncodeKey encodes a key with groups of letters and numbers.
func EncodeKey(k []byte) string {
	encodedKey := encKey.EncodeToString(k)
	l := len(encodedKey)
	if l == 0 {
		return ""
	}

	di := 0

	result := make([]byte, l+(l-1)/keyGroupSize+1)
	for l > keyGroupSize {
		result[di] = keySeparator
		di++

		copy(result[di:], encodedKey[:keyGroupSize])

		encodedKey = encodedKey[keyGroupSize:]
		di += keyGroupSize
		l -= keyGroupSize
	}

	result[di] = keySeparator
	di++
	copy(result[di:], encodedKey[:])

	return string(result[1:])
}
