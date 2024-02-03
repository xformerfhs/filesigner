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

package numberhelper

import (
	"errors"
	"filesigner/slicehelper"
)

// ******** Private constants ********

const maxByteCounterLen = 8

// ******** Public types ********

type ByteCounter struct {
	counter []byte
	first   byte
	last    byte
}

// ******** Public creation functions ********

func NewByteCounter(length byte) (*ByteCounter, error) {
	if length == 0 {
		return nil, errors.New("Byte counter length must not be 0")
	}

	if length > maxByteCounterLen {
		return nil, errors.New("Byte counter length is too large")
	}

	maxIndex := length - 1

	return &ByteCounter{counter: make([]byte, length), first: maxIndex, last: maxIndex}, nil
}

func NewByteCounterForCount(count uint64) (*ByteCounter, error) {
	return NewByteCounter(ByteCountForNumber(count))
}

// ******** Public functions ********

// -------- Counting functions --------

func (bc *ByteCounter) Inc() {
	counter := bc.counter
	i := bc.last
	for {
		a := counter[i]
		a++
		counter[i] = a
		if a != 0 || i == 0 {
			break
		}

		i--
	}

	bc.first = i
}

func (bc *ByteCounter) Dec() {
	counter := bc.counter
	i := byte(len(counter) - 1)
	for {
		a := counter[i]
		a--
		counter[i] = a
		if a != 0xff || i == 0 {
			break
		}

		i--
	}
}

func (bc *ByteCounter) Zero() {
	counter := bc.counter
	slicehelper.ClearInteger(counter)
}

func (bc *ByteCounter) Slice() []byte {
	return bc.counter[bc.first:]
}

func (bc *ByteCounter) FullSlice() []byte {
	return bc.counter
}

// -------- Setter and getter functions --------

func (bc *ByteCounter) SetCount(value uint64) {
	counter := bc.counter
	i := bc.last
	for {
		counter[i] = byte(value)
		value >>= 8
		if value == 0 {
			break
		}

		i--
	}

	bc.first = i
}

func (bc *ByteCounter) GetCount() uint64 {
	result := uint64(0)
	first := bc.first
	counter := bc.counter
	i := bc.last
	for {
		result <<= 8
		result |= uint64(counter[i])

		if i == first {
			break
		}

		i--
	}

	return result
}
