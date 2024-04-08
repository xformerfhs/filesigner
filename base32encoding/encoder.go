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
	return enc.EncodeToString(b)
}

// EncodeToBytes encodes a byte slice as a byte slice.
func EncodeToBytes(b []byte) []byte {
	result := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(result, b)

	return result
}

// EncodeKey encodes a key with groups of letters and numbers.
func EncodeKey(key []byte) string {
	encodedKey := encKey.EncodeToString(key)
	encodedKeyLength := len(encodedKey)
	if encodedKeyLength == 0 {
		return ``
	}

	destinationIndex := 0

	result := make([]byte, encodedKeyLength+(encodedKeyLength-1)/keyGroupSize+1)
	for encodedKeyLength > keyGroupSize {
		result[destinationIndex] = keySeparator
		destinationIndex++

		copy(result[destinationIndex:], encodedKey[:keyGroupSize])

		encodedKey = encodedKey[keyGroupSize:]
		destinationIndex += keyGroupSize
		encodedKeyLength -= keyGroupSize
	}

	result[destinationIndex] = keySeparator
	destinationIndex++
	copy(result[destinationIndex:], encodedKey[:])

	return string(result[1:])
}
