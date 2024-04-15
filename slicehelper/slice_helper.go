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
//    2024-03-17: V1.1.0: Add FillToCap.
//

// Package slicehelper implements helper functions for slices.
package slicehelper

import "filesigner/constraints"

// ******** Public functions ********

// Fill fills a slice with a value in an efficient way up to its length.
func Fill[S ~[]T, T any](s S, v T) {
	sLen := len(s)

	if sLen > 0 {
		doFill(s, v, sLen)
	}
}

// FillToCap fills a slice with a value in an efficient way up to its capacity.
func FillToCap[S ~[]T, T any](s S, v T) {
	sLen := cap(s)

	if sLen > 0 {
		doFill(s[:sLen], v, sLen)
	}
}

// ClearNumber clears a number type slice.
func ClearNumber[S ~[]T, T constraints.Number](a S) {
	FillToCap(a, 0)
}

// NewReverse returns a new slice with the elements in the reverse order of the argument.
func NewReverse[S ~[]T, T any](a S) S {
	aLen := len(a)
	result := make(S, aLen)

	i := aLen
	for _, e := range a {
		i--
		result[i] = e
	}

	return result
}

// Concat returns a new slice concatenating the passed in slices.
// This is a streamlined version of the slices.Concat function of Go V1.22.
func Concat[S ~[]T, T any](slices ...S) S {
	// 1. Calculate total size.
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	// 2. Make new slice with the total size as the capacity and 0 length.
	result := make(S, 0, size)

	// 3. Append all source slices.
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}

// Copy makes a copy of a slice.
func Copy[S ~[]T, T any](a S) S {
	// This is twice as fast, as using append.
	result := make(S, len(a))
	copy(result, a)
	return result
}

// SetCap sets the capacity of a slice to be at least n.
// If n is negative or too large to allocate the memory, SetCap panics.
func SetCap[S ~[]T, T any](s S, n int) S {
	if n < 0 {
		panic(`cannot be negative`)
	}

	c := cap(s)
	if n -= c; n > 0 {
		s = append(s[:c], make(S, n)...)[:len(s)]
	}

	return s
}

// Prepend adds elements at the beginning of a slice s.
func Prepend[S ~[]T, T any](s S, e ...T) []T {
	return append(e, s...)
}

// ******** Private functions ********

// doFill fills a slice in an efficient way.
func doFill[S ~[]T, T any](s S, v T, l int) {
	// Put the value into the first slice element
	s[0] = v

	// Incrementally duplicate the value into the rest of the slice
	for j := 1; j < l; j <<= 1 {
		copy(s[j:], s[:j])
	}
}
