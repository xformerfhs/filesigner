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
//    2024-03-04: V1.0.0: Created.
//

package paddedhasher

import (
	"fmt"
	"hash"
)

// ******** Public types ********

// PaddedHasher is a hash.Hash with a padding.
type PaddedHasher struct {
	hasher            hash.Hash
	padding           []byte
	paddingSplitIndex int
	count             uint64
}

// ******** Type creation functions ********

// NewPaddedHasherAsHash creates a new padded hasher as a hash.Hash.
func NewPaddedHasherAsHash(hashFunc hash.Hash, padding []byte) hash.Hash {
	return hash.Hash(NewPaddedHasher(hashFunc, padding))
}

// NewPaddedHasher creates a new padded hasher.
func NewPaddedHasher(hashFunc hash.Hash, padding []byte) *PaddedHasher {
	result := &PaddedHasher{
		hasher:            hashFunc,
		padding:           padding,
		paddingSplitIndex: len(padding) >> 1,
		count:             0,
	}

	result.Reset()

	return result
}

// ******** Public functions ********

// Write writes data to the context hasher.
func (ph *PaddedHasher) Write(p []byte) (int, error) {
	n, err := ph.hasher.Write(p)
	ph.count += uint64(n)
	return n, err
}

// Sum returns the hash either into b or creates a new byte slice.
func (ph *PaddedHasher) Sum(b []byte) []byte {
	hasher := ph.hasher
	// This write does not increase the counter as this is not written data.
	_, err := hasher.Write(ph.padding[ph.paddingSplitIndex:])
	if err != nil {
		panic(fmt.Sprintf(`Write to hash in Sum had error: %v`, err))
	}
	return hasher.Sum(b)
}

// Reset resets the padded hasher.
func (ph *PaddedHasher) Reset() {
	hasher := ph.hasher
	hasher.Reset()
	// This write does not increase the counter as this is not written data.
	hasher.Write(ph.padding[:ph.paddingSplitIndex])
}

// Size returns the size of the hash value.
func (ph *PaddedHasher) Size() int {
	return ph.hasher.Size()
}

// BlockSize returns the block size of the hash algorithm.
func (ph *PaddedHasher) BlockSize() int {
	return ph.hasher.BlockSize()
}

// Count returns the number of bytes that have been written to this hasher.
func (ph *PaddedHasher) Count() uint64 {
	return ph.count
}
