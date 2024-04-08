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
//    2024-02-15: V1.0.0: Created.
//    2024-03-05: V2.0.0: Added 64 bit methods.
//

package numberhelper

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestByteCountForUint32(t *testing.T) {
	bc := ByteCountForUint32(0)
	if bc != 1 {
		t.Fatalf(`Invalid byte count for 0: %d`, bc)
	}
	number := uint32(0)
	for i := byte(1); i <= 4; i++ {
		number = (number << 8) | 0xff
		bc = ByteCountForUint32(number)
		if bc != i {
			t.Fatalf(`Invalid byte count for %d: %d (should be %d)`, number, bc, i)
		}
	}
}

func TestUint32AsBytesBoundaries(t *testing.T) {
	counter := uint32(0)
	slice := Uint32AsShortestBigEndianBytes(counter)
	if len(slice) != 1 {
		t.Fatal(`Converting 0 to bytes does not yield a slice of length 1`)
	}
	if slice[0] != 0 {
		t.Fatal(`Converting 0 to bytes does not yield a byte with value 0`)
	}

	counter--
	slice = Uint32AsShortestBigEndianBytes(counter)
	if len(slice) != 4 {
		t.Fatal(`Converting MaxUint32 to bytes does not yield a slice of length 4`)
	}
	if bytes.Compare(slice, []byte{0xff, 0xff, 0xff, 0xff}) != 0 {
		t.Fatal(`Converting MaxUint32 to bytes does not yield an all 0xff byte slice`)
	}
}

const loopCount = 1_000_000

func TestUint32AsBytesRandom(t *testing.T) {
	var counter uint32
	var slice []byte
	for i := 0; i < loopCount; i++ {
		counter = rand.Uint32()
		slice = Uint32AsShortestBigEndianBytes(counter)
		if byte(len(slice)) != ByteCountForUint32(counter) {
			t.Fatalf(`Converting %d has invalid slice length %d`, counter, len(slice))
		}
		sliceValue := BigEndianBytesAsUint32(slice)
		if sliceValue != counter {
			t.Fatalf(`Converting %d to bytes and back yields %d`, counter, sliceValue)
		}
	}
}

func TestUint64AsBytesBoundaries(t *testing.T) {
	counter := uint64(0)
	slice := Uint64AsShortestBigEndianBytes(counter)
	if len(slice) != 1 {
		t.Fatal(`Converting 0 to bytes does not yield a slice of length 1`)
	}
	if slice[0] != 0 {
		t.Fatal(`Converting 0 to bytes does not yield a byte with value 0`)
	}

	counter--
	slice = Uint64AsShortestBigEndianBytes(counter)
	if len(slice) != 8 {
		t.Fatal(`Converting MaxUint64 to bytes does not yield a slice of length 4`)
	}
	if bytes.Compare(slice, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) != 0 {
		t.Fatal(`Converting MaxUint64 to bytes does not yield an all 0xff byte slice`)
	}
}

func TestUint64AsBytesRandom(t *testing.T) {
	var counter uint64
	var slice []byte
	for i := 0; i < loopCount; i++ {
		counter = rand.Uint64()
		slice = Uint64AsShortestBigEndianBytes(counter)
		if byte(len(slice)) != ByteCountForUint64(counter) {
			t.Fatalf(`Converting %d has invalid slice length %d`, counter, len(slice))
		}
		sliceValue := BigEndianBytesAsUint64(slice)
		if sliceValue != counter {
			t.Fatalf(`Converting %d to bytes and back yields %d`, counter, sliceValue)
		}
	}
}
