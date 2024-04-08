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
// according to the current platform, i.e. either case-sensitive or case-insensitive.
type FileSystemStringSet struct {
	s *Set[string]
}

// ******** Private variables ********

// isFileSystemCaseInsensitive is "true", if the file system is case-insensitive,
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
func NewFileSystemStringSetWithLength(l int) *FileSystemStringSet {
	return &FileSystemStringSet{s: NewWithLength[string](l)}
}

// NewFileSystemStringSetWithElements creates a new platform-specific string set which contains the given elements.
func NewFileSystemStringSetWithElements(e ...string) *FileSystemStringSet {
	var n []string

	if isFileSystemCaseInsensitive {
		n = make([]string, len(e))

		for i, k := range e {
			n[i] = foldCaser.String(k)
		}
	} else {
		n = e
	}

	return &FileSystemStringSet{s: NewWithElements[string](n...)}
}

// ******** Public functions ********

// -------- Case sensitivity functions --------

// IsCaseInsensitive return "true", if the file system is case-insensitive, "false", if it is case-sensitive.
func (f *FileSystemStringSet) IsCaseInsensitive() bool {
	return isFileSystemCaseInsensitive
}

// -------- Element functions --------

// Elements returns the elements of the set as a slice.
func (f *FileSystemStringSet) Elements() []string {
	return f.s.Elements()
}

// Size returns the number of elements in the set.
func (f *FileSystemStringSet) Size() int {
	return f.s.Size()
}

// Add adds an element to the set.
func (f *FileSystemStringSet) Add(e string) {
	f.s.Add(properCasedString(e))
}

// Remove deletes an element from the set.
func (f *FileSystemStringSet) Remove(e string) {
	f.s.Remove(properCasedString(e))
}

// Clear removes all elements from the set.
func (f *FileSystemStringSet) Clear() {
	f.s.Clear()
}

// Contains tests whether the element is in the set.
func (f *FileSystemStringSet) Contains(e string) bool {
	return f.s.Contains(properCasedString(e))
}

// Do calls function fn for each element in the set.
func (f *FileSystemStringSet) Do(fn func(string)) {
	f.s.Do(fn)
}

// -------- Set functions --------

// Difference finds the difference between two sets, i.e. the elements that are not contained in the other set.
func (f *FileSystemStringSet) Difference(o *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Difference(o.s)}
}

// Intersection finds the intersection of two sets, i.e. the elements that are contained in both sets.
func (f *FileSystemStringSet) Intersection(o *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Intersection(o.s)}
}

// Union returns the union of two sets.
func (f *FileSystemStringSet) Union(o *FileSystemStringSet) *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Union(o.s)}
}

// -------- Set test function ---------

// Equal returns true if both sets have equal lengths and equal elements.
func (f *FileSystemStringSet) Equal(o *FileSystemStringSet) bool {
	return f.s.Equal(o.s)
}

// IsSubsetOf tests if this set is a subset of the other set.
func (f *FileSystemStringSet) IsSubsetOf(o *FileSystemStringSet) bool {
	return f.s.IsSubsetOf(o.s)
}

// IsProperSubsetOf tests if this set is a proper subset of the other set.
func (f *FileSystemStringSet) IsProperSubsetOf(o *FileSystemStringSet) bool {
	return f.s.IsProperSubsetOf(o.s)
}

// -------- Management functions ---------

// Copy creates a copy of the set.
func (f *FileSystemStringSet) Copy() *FileSystemStringSet {
	return &FileSystemStringSet{s: f.s.Copy()}
}

// -------- String case handling functions ---------

// properCasedString returns either the source string, or the folded source string,
// depending on the value of isFileSystemCaseInsensitive.
func properCasedString(s string) string {
	if isFileSystemCaseInsensitive {
		return foldCaser.String(s)
	} else {
		return s
	}
}
