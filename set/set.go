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

package set

import "golang.org/x/exp/maps"

// ******** Public types ********

// Set implements a set, i.e. a data structure that contains a given element exactly once.
type Set[K comparable] struct {
	m map[K]bool
}

// ******** Public creation functions ********

// New creates a new set.
func New[K comparable]() *Set[K] {
	return &Set[K]{m: make(map[K]bool)}
}

// NewWithLength creates a new set with a given length.
func NewWithLength[K comparable](length int) *Set[K] {
	return &Set[K]{m: make(map[K]bool, length)}
}

// NewWithElements creates a new set which contains the given elements.
func NewWithElements[K comparable](elements ...K) *Set[K] {
	n := make(map[K]bool, len(elements))

	for _, k := range elements {
		n[k] = true
	}

	return &Set[K]{m: n}
}

// ******** Public functions ********

// -------- Element function --------

// Elements returns the elements of the set as a slice.
func (s *Set[K]) Elements() []K {
	return maps.Keys(s.m)
}

// Len returns the number of elements in the set.
func (s *Set[K]) Len() int {
	return len(s.m)
}

// Add adds an element to the set.
func (s *Set[K]) Add(element K) {
	s.m[element] = true
}

// Remove deletes an element from the set.
func (s *Set[K]) Remove(element K) {
	delete(s.m, element)
}

// Clear removes all elements from the set.
func (s *Set[K]) Clear() {
	clear(s.m)
}

// Contains tests whether or not the element is in the set.
func (s *Set[K]) Contains(element K) bool {
	return s.m[element]
}

// Do calls function fn for each element in the set.
func (s *Set[K]) Do(fn func(K)) {
	for k := range s.m {
		fn(k)
	}
}

// -------- Set functions --------

// Difference finds the difference between two sets.
func (s *Set[K]) Difference(other *Set[K]) *Set[K] {
	n := make(map[K]bool)

	for k := range s.m {
		if !other.m[k] {
			n[k] = true
		}
	}

	return &Set[K]{m: n}
}

// Intersection finds the intersection of two sets.
func (s *Set[K]) Intersection(other *Set[K]) *Set[K] {
	n := make(map[K]bool)

	for k := range s.m {
		if other.m[k] {
			n[k] = true
		}
	}

	return &Set[K]{m: n}
}

// Union returns the union of two sets.
func (s *Set[K]) Union(other *Set[K]) *Set[K] {
	n := make(map[K]bool, len(s.m)+len(other.m))

	// 1. Add elements of this set.
	for k := range s.m {
		n[k] = true
	}

	// 2. Add elements of other set.
	for k := range other.m {
		n[k] = true
	}

	return &Set[K]{m: n}
}

// -------- Set test function ---------

// IsSubsetOf tests whether or not this set is a subset of "other".
func (s *Set[K]) IsSubsetOf(other *Set[K]) bool {
	if s.Len() > other.Len() {
		return false
	}

	for k := range s.m {
		if !other.m[k] {
			return false
		}
	}

	return true
}

// IsProperSubsetOf tests whether or not this set is a proper subset of "other".
func (s *Set[K]) IsProperSubsetOf(other *Set[K]) bool {
	return s.Len() < other.Len() && s.IsSubsetOf(other)
}
