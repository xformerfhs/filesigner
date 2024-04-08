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

// ******** Public methods ********

// --------- 32 bit methods --------

// ByteCountForUint32 counts how many bytes are needed to represent the given number.
func ByteCountForUint32(number uint32) byte {
	// Using a switch statement is about 14% faster than a loop.
	switch {
	case number <= 0xff:
		return 1
	case number <= 0xffff:
		return 2
	case number <= 0xffffff:
		return 3
	default:
		return 4
	}
}

// BigEndianBytesAsUint32 returns the uint32 from the supplied byte array in big endian byte order.
func BigEndianBytesAsUint32(slice []byte) uint32 {
	sliceLen := len(slice)
	if sliceLen > 4 {
		panic(`slice is too long`)
	}

	result := uint32(0)
	for i := 0; i < sliceLen; i++ {
		result = (result << 8) | uint32(slice[i])
	}

	return result
}

// static32Buffer is the static buffer that contains a converted uint32.
var static32Buffer = [4]byte{}

// StaticIntAsShortestBigEndianBytes returns an int as the shortest possible byte slice in big endian byte order.
// It uses a static buffer, so that the content of the buffer changes with each call.
// Only use it when the returned value is used directly after the call and there is no concurrency.
func StaticIntAsShortestBigEndianBytes(value int) []byte {
	return StaticUint32AsShortestBigEndianBytes(uint32(value))
}

// StaticUint32AsShortestBigEndianBytes returns an uint32 as the shortest possible byte slice in big endian byte order.
// It uses a static buffer, so that the content of the buffer changes with each call.
// Only use it when the returned value is used directly after the call and there is no concurrency.
func StaticUint32AsShortestBigEndianBytes(value uint32) []byte {
	return Uint32AsShortestBigEndianBytesIntoBuffer(static32Buffer[:], value)
}

// IntAsShortestBigEndianBytes returns an int as the shortest possible byte slice in big endian byte order.
func IntAsShortestBigEndianBytes(value int) []byte {
	return Uint32AsShortestBigEndianBytes(uint32(value))
}

// Uint32AsShortestBigEndianBytes returns an uint32 as the shortest possible byte slice in big endian byte order.
func Uint32AsShortestBigEndianBytes(value uint32) []byte {
	buffer := make([]byte, 4)
	return Uint32AsShortestBigEndianBytesIntoBuffer(buffer, value)
}

// Uint32AsShortestBigEndianBytesIntoBuffer returns an uint32 as the shortest possible byte slice in big endian byte order.
func Uint32AsShortestBigEndianBytesIntoBuffer(buffer []byte, value uint32) []byte {
	buffer[3] = byte(value)
	value >>= 8
	if value == 0 {
		return buffer[3:]
	}

	buffer[2] = byte(value)
	value >>= 8
	if value == 0 {
		return buffer[2:]
	}

	buffer[1] = byte(value)
	value >>= 8
	if value == 0 {
		return buffer[1:]
	}

	buffer[0] = byte(value)
	return buffer
}

// --------- 64 bit methods --------

// ByteCountForUint64 counts how many bytes are needed to represent the given number.
func ByteCountForUint64(number uint64) byte {
	// Using a switch statement is about 14% faster than a loop.
	switch {
	case number <= 0xff:
		return 1
	case number <= 0xffff:
		return 2
	case number <= 0xffffff:
		return 3
	case number <= 0xffffffff:
		return 4
	case number <= 0xffffffffff:
		return 5
	case number <= 0xffffffffffff:
		return 6
	case number <= 0xffffffffffffff:
		return 7
	default:
		return 8
	}
}

// BigEndianBytesAsUint64 returns the uint64 from the supplied byte array in big endian byte order.
func BigEndianBytesAsUint64(slice []byte) uint64 {
	sliceLen := len(slice)
	if sliceLen > 8 {
		panic(`slice is too long`)
	}

	result := uint64(0)
	for i := 0; i < sliceLen; i++ {
		result = (result << 8) | uint64(slice[i])
	}

	return result
}

// static64Buffer is the static buffer that contains a converted uint32.
var static64Buffer = [8]byte{}

// StaticUint64AsShortestBigEndianBytes returns an uint64 as the shortest possible byte slice in big endian byte order.
// It uses a static buffer, so that the content of the buffer changes with each call.
// Only use it when the returned value is used directly after the call and there is no concurrency.
func StaticUint64AsShortestBigEndianBytes(value uint64) []byte {
	return Uint64AsShortestBigEndianBytesIntoBuffer(static64Buffer[:], value)
}

// Int64AsShortestBigEndianBytes returns an int64 as the shortest possible byte slice in big endian byte order.
func Int64AsShortestBigEndianBytes(value int64) []byte {
	return Uint64AsShortestBigEndianBytes(uint64(value))
}

// Uint64AsShortestBigEndianBytes returns an uint64 as the shortest possible byte slice in big endian byte order.
func Uint64AsShortestBigEndianBytes(value uint64) []byte {
	buffer := make([]byte, 8)
	return Uint64AsShortestBigEndianBytesIntoBuffer(buffer, value)
}

// Uint64AsShortestBigEndianBytesIntoBuffer returns an uint64 as the shortest possible byte slice in big endian byte order.
func Uint64AsShortestBigEndianBytesIntoBuffer(buffer []byte, value uint64) []byte {
	for i := 7; i >= 0; i-- {
		buffer[i] = byte(value)
		value >>= 8
		if value == 0 {
			return buffer[i:]
		}
	}

	return buffer
}
