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

import (
	"filesigner/maphelper"
)

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
func NewWithLength[K comparable](l int) *Set[K] {
	return &Set[K]{m: make(map[K]bool, l)}
}

// NewWithElements creates a new set which contains the given elements.
func NewWithElements[K comparable](e ...K) *Set[K] {
	n := make(map[K]bool, len(e))

	for _, k := range e {
		n[k] = true
	}

	return &Set[K]{m: n}
}

// ******** Public functions ********

// -------- Element function --------

// Elements returns the elements of the set as a slice.
func (s *Set[K]) Elements() []K {
	return maphelper.Keys(s.m)
}

// Size returns the number of elements in the set.
func (s *Set[K]) Size() int {
	return len(s.m)
}

// Add adds an element to the set.
func (s *Set[K]) Add(e K) {
	s.m[e] = true
}

// Remove deletes an element from the set.
func (s *Set[K]) Remove(e K) {
	delete(s.m, e)
}

// Clear removes all elements from the set.
func (s *Set[K]) Clear() {
	clear(s.m)
}

// Contains tests if the element is in the set.
func (s *Set[K]) Contains(e K) bool {
	return s.m[e]
}

// Do calls function fn for each element in the set.
func (s *Set[K]) Do(fn func(K)) {
	for k := range s.m {
		fn(k)
	}
}

// -------- Set functions --------

// Difference finds the difference between two sets, i.e. the elements that are not contained in the other set.
func (s *Set[K]) Difference(o *Set[K]) *Set[K] {
	n := make(map[K]bool)

	om := o.m
	for k := range s.m {
		if !om[k] {
			n[k] = true
		}
	}

	return &Set[K]{m: n}
}

// Intersection finds the intersection of two sets, i.e. the elements that are contained in both sets.
func (s *Set[K]) Intersection(o *Set[K]) *Set[K] {
	n := make(map[K]bool)

	om := o.m
	for k := range s.m {
		if om[k] {
			n[k] = true
		}
	}

	return &Set[K]{m: n}
}

// Union returns the union of two sets.
func (s *Set[K]) Union(o *Set[K]) *Set[K] {
	n := make(map[K]bool, len(s.m)+len(o.m))

	// 1. Add elements of this set.
	for k := range s.m {
		n[k] = true
	}

	// 2. Add elements of other set.
	for k := range o.m {
		n[k] = true
	}

	return &Set[K]{m: n}
}

// -------- Set test function ---------

// Equal returns true if both sets have equal lengths and equal elements.
func (s *Set[K]) Equal(o *Set[K]) bool {
	// Equal means 1. the lengths of both sets are equal and...
	if s.Size() != o.Size() {
		return false
	}

	// ... 2. this set is a subset of the other set.
	return s.isSubsetOf(o)
}

// IsSubsetOf tests if this set is a subset of the other set.
func (s *Set[K]) IsSubsetOf(o *Set[K]) bool {
	if s.Size() > o.Size() {
		return false
	}

	return s.isSubsetOf(o)
}

// IsProperSubsetOf tests if this set is a proper subset of the other set.
func (s *Set[K]) IsProperSubsetOf(o *Set[K]) bool {
	return s.Size() < o.Size() && s.IsSubsetOf(o)
}

// -------- Management functions ---------

// Copy creates a copy of the set.
func (s *Set[K]) Copy() *Set[K] {
	m := make(map[K]bool, s.Size())

	for k := range s.m {
		m[k] = true
	}

	return &Set[K]{m: m}
}

// ******** Private functions ********

// isSubsetOf tests if this set is a subset of the other set.
// This is the subset check without the length checks.
func (s *Set[K]) isSubsetOf(o *Set[K]) bool {
	om := o.m
	for k := range s.m {
		if !om[k] {
			return false
		}
	}

	return true
}
