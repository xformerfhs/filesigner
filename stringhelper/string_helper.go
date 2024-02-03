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

package stringhelper

import (
	"strings"
	"unsafe"
)

// UnsafeStringBytes returns a byte slice that points to the bytes of the supplied string.
// No bytes are copied. Attention: This is *unsafe*! Do not change those bytes!
func UnsafeStringBytes(s string) []byte {
	// This is a streamlined version of
	// https://josestg.medium.com/140x-faster-string-to-byte-and-byte-to-string-conversions-with-zero-allocation-in-go-200b4d7105fc .
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeStringFrom returns a string that starts at the specified offset in another string.
// This avoids copying a string.
func UnsafeStringFrom(s string, offset int) string {
	return unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(unsafe.StringData(s)), offset)), len(s)-offset)
}

// HasCaseInsensitivePrefix tests whether the string s begins with prefix in a case-insensitive way.
func HasCaseInsensitivePrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}
