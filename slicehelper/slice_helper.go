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

package slicehelper

import "golang.org/x/exp/constraints"

// ******** Public functions ********

// Fill fills a generic slice with a generic value in an efficient way.
func Fill[S ~[]T, T any](a S, value T) {
	aLen := ensureLengthIsCapacity(&a)

	if aLen > 0 {
		// Put the value into the first slice element
		a[0] = value

		// Incrementally duplicate the value into the rest of the slice
		for j := 1; j < aLen; j <<= 1 {
			copy(a[j:], a[:j])
		}
	}
}

// ClearInteger clears an integer type slice.
func ClearInteger[S ~[]T, T constraints.Integer](a S) {
	Fill(a, 0)
}

// MakeCopy makes a copy of slice
func MakeCopy[S ~[]T, T any](a S) S {
	result := make([]T, len(a))
	copy(result, a)
	return result
}

// Prepend adds an element v at the beginning of a slice s.
func Prepend[T any](v T, s []T) []T {
	return append([]T{v}, s...)
}

// ******** Private functions ********

// ensureLengthIsCapacity ensures that the length of the slice is its capacity.
// We need the address of the slice as the parameter. If the '*' would be missing
// we would get a copy of the slice and not the slice itself.
func ensureLengthIsCapacity[S ~[]T, T any](a *S) int {
	ra := *a
	aLen := len(ra)
	aCap := cap(ra)
	if aLen != aCap {
		*a = ra[:aCap]
		aLen = aCap
	}

	return aLen
}
