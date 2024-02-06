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
//    2024-02-05: V1.0.0: Created.
//

package set

import (
	"golang.org/x/text/cases"
	"runtime"
)

// ******** Public types ********

// FileSystemStringSet implements a set of strings that has compare semantics
// according to the the current platform, i.e. either case-sensitive or case-insensitive.
type FileSystemStringSet struct {
	s *Set[string]
}

// ******** Private variables ********

// isFileSystemCaseInsensitive is "true", if the file system is case insensitive,
// "false" otherwise.
var isFileSystemCaseInsensitive = runtime.GOOS == `windows`

// foldCaser is a caser that returns a folded string, useful for case-insensitive comparison.
var foldCaser = cases.Fold()

// ******** Public creation functions ********

// NewFileSystemStringSet creates a new platform-specific string set.
func NewFileSystemStringSet() *FileSystemStringSet {
	return &FileSystemStringSet{s: New[string]()}
}

// NewFileSystemStringSetWithLength creates a new platform-specific string set with a given length.
func NewFileSystemStringSetWithLength(length int) *FileSystemStringSet {
	return &FileSystemStringSet{s: NewWithLength[string](length)}
}

// NewFileSystemStringSetWithElements creates a new platform-specific string set which contains the given elements.
func NewFileSystemStringSetWithElements(elements ...string) *FileSystemStringSet {
	var n []string

	if isFileSystemCaseInsensitive {
		n = make([]string, len(elements))

		for i, k := range elements {
			n[i] = foldCaser.String(k)
		}
	} else {
		n = elements
	}

	return &FileSystemStringSet{s: NewWithElements[string](n...)}
}

// ******** Public functions ********

// -------- Element function --------

// Elements returns the elements of the set as a slice.
func (f *FileSystemStringSet) Elements() []string {
	return f.s.Elements()
}

// Len returns the number of elements in the set.
func (f *FileSystemStringSet) Len() int {
	return f.s.Len()
}

// Add adds an element to the set.
func (f *FileSystemStringSet) Add(element string) {
	f.s.Add(properCasedString(element))
}

// Remove deletes an element from the set.
func (f *FileSystemStringSet) Remove(element string) {
	f.s.Remove(properCasedString(element))
}

// Clear removes all elements from the set.
func (f *FileSystemStringSet) Clear() {
	f.s.Clear()
}

// Contains tests whether or not the element is in the set.
func (f *FileSystemStringSet) Contains(element string) bool {
	return f.s.Contains(properCasedString(element))
}

// Do calls function fn for each element in the set.
func (f *FileSystemStringSet) Do(fn func(string)) {
	f.s.Do(fn)
}

// -------- Set functions --------

// Difference finds the difference between two sets.
func (f *FileSystemStringSet) Difference(other *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Difference(other.s)}
}

// Intersection finds the intersection of two sets.
func (f *FileSystemStringSet) Intersection(other *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Intersection(other.s)}
}

// Union returns the union of two sets.
func (f *FileSystemStringSet) Union(other *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Union(other.s)}
}

// -------- Set test function ---------

// IsSubsetOf tests whether or not this set is a subset of "other".
func (f *FileSystemStringSet) IsSubsetOf(other *FileSystemStringSet) bool {
	return f.s.IsSubsetOf(other.s)
}

// IsProperSubsetOf tests whether or not this set is a proper subset of "other".
func (f *FileSystemStringSet) IsProperSubsetOf(other *FileSystemStringSet) bool {
	return f.s.IsProperSubsetOf(other.s)
}

// -------- String case handling functions ---------

// properCasedString returns either the source string, or the folded source string,
// depending on the value of isFileSystemCaseInsensitive.
func properCasedString(source string) string {
	if isFileSystemCaseInsensitive {
		return foldCaser.String(source)
	} else {
		return source
	}
}
