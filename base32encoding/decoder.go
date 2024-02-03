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

import "strings"

// ******** Public functions ********

// DecodeFromString decodes a string into a byte slice.
func DecodeFromString(s string) ([]byte, error) {
	emx.Lock()
	defer emx.Unlock()

	return enc.DecodeString(s)
}

// DecodeKey decodes a key
func DecodeKey(s string) ([]byte, error) {
	sl := len(s)
	result := make([]byte, sl)

	di := 0
	t := 0

	for {
		t = strings.IndexByte(s, keySeparator)
		if t < 0 {
			copy(result[di:], s)
			di += sl
			break
		}

		copy(result[di:], s[:t])
		di += t
		t++
		s = s[t:]
		sl -= t
	}

	return encKey.DecodeString(string(result[:di]))
}
