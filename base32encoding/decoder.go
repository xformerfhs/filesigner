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

import (
	"errors"
	"strings"
)

// ******** Public functions ********

// DecodeFromString decodes a string into a byte slice.
func DecodeFromString(s string) ([]byte, error) {
	return enc.DecodeString(s)
}

// DecodeFromBytes decodes a byte slice into a byte slice.
func DecodeFromBytes(b []byte) ([]byte, error) {
	result := make([]byte, enc.DecodedLen(len(b)))
	n, err := enc.Decode(result, b)

	return result[:n], err
}

// DecodeKey decodes a key.
func DecodeKey(keyId string) ([]byte, error) {
	keyIdLength := len(keyId)
	result := make([]byte, keyIdLength)

	destinationIndex := 0
	separatorPosition := 0
	separatorCount := 0
	for {
		separatorPosition = strings.IndexByte(keyId, keySeparator)
		if separatorPosition < 0 {
			copy(result[destinationIndex:], keyId)
			destinationIndex += keyIdLength
			break
		}

		if separatorPosition != keyGroupSize {
			return nil, errors.New(`Invalid group size in key id`)
		}
		separatorCount++
		copy(result[destinationIndex:], keyId[:separatorPosition])
		destinationIndex += separatorPosition
		separatorPosition++
		keyId = keyId[separatorPosition:]
		keyIdLength -= separatorPosition
	}

	return encKey.DecodeString(string(result[:destinationIndex]))
}
