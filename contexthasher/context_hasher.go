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

package contexthasher

import (
	"filesigner/numberhelper"
	"hash"
)

// ******** Public types ********

// ContextHasher is a hash.Hash with a context.
type ContextHasher struct {
	hasher     hash.Hash
	context    []byte
	contextLen []byte
}

// ******** Private constants ********

// separator separates the context from the data.
var separator = []byte{0x33, 0x17, 0xd1, 0xdb, 0xc2, 0xf1}

// ******** Type creation functions ********

// NewContextHasher creates a new context hasher.
func NewContextHasher(hashFunc hash.Hash, contextBytes []byte) hash.Hash {
	contextLen := uint64(len(contextBytes))
	bc, _ := numberhelper.NewByteCounterForCount(contextLen)
	bc.SetCount(contextLen)
	result := &ContextHasher{hasher: hashFunc, context: contextBytes, contextLen: bc.Slice()}

	// context points to the bytes of the string. It is not a copy.
	result.context = contextBytes
	result.Reset()

	return result
}

// ******** Public functions ********

// Write writes data to the context hasher.
func (ch *ContextHasher) Write(p []byte) (int, error) {
	return ch.hasher.Write(p)
}

// Sum returns the hash either into b or creates a new byte slice.
func (ch *ContextHasher) Sum(b []byte) []byte {
	return ch.hasher.Sum(b)
}

// Reset resets the context hasher.
func (ch *ContextHasher) Reset() {
	hasher := ch.hasher
	hasher.Reset()
	hasher.Write(ch.context)
	hasher.Write(ch.contextLen)
	hasher.Write(separator)
}

// Size returns the size of the hash value.
func (ch *ContextHasher) Size() int {
	return ch.hasher.Size()
}

// BlockSize returns the block size of the hash algorithm.
func (ch *ContextHasher) BlockSize() int {
	return ch.hasher.BlockSize()
}
