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
	"filesigner/maphelper"
	"testing"
)

func TestFSStringSet(t *testing.T) {
	s := NewFileSystemStringSet()
	s.Add(`onlylowercase`)
	s.Add(`onlyLOWERCASE`)
	if s.IsCaseInsensitive() {
		if s.Size() != 1 {
			t.Fatal(`Case insensitive file system string set is not case-insensitive`)
		}
	} else {
		if s.Size() != 2 {
			t.Fatal(`Case sensitive file system string set is not case-sensitive`)
		}
	}
}

func TestFSStringSetUTF8(t *testing.T) {
	s := NewFileSystemStringSet()

	s.Add(`äöüßéèê`)
	s.Add(`ÄÖÜSSÉÈÊ`)
	if s.IsCaseInsensitive() {
		if s.Size() != 1 {
			t.Fatal(`Case insensitive file system string set is not case-insensitive`)
		}
	} else {
		if s.Size() != 2 {
			t.Fatal(`Case sensitive file system string set is not case-sensitive`)
		}
	}
}

func TestFSStringClear(t *testing.T) {
	s := NewFileSystemStringSet()
	s.Add(`Whatever`)
	s.Add(`ØıƉ`)
	if s.Size() != 2 {
		t.Fatal(`Not enough elements`)
	}
	s.Clear()
	if s.Size() != 0 {
		t.Fatal(`Zero does not work`)
	}
}

func TestFSStringRemove(t *testing.T) {
	s := NewFileSystemStringSet()
	s.Add(`Whereever`)
	s.Add(`ƉØı`)
	if s.Size() != 2 {
		t.Fatal(`Not enough elements`)
	}
	s.Remove(`Whereever`)
	if s.Size() != 1 {
		t.Fatal(`Remove does not work`)
	}
}
func TestFSStringUnion(t *testing.T) {
	s := NewFileSystemStringSet()
	s.Add(`pi`)
	s.Add(`euler`)
	s.Add(`gauss`)

	u := NewFileSystemStringSet()
	u.Add(`pi`)
	u.Add(`lagrange`)
	u.Add(`liouville`)

	n := s.Union(u)
	if n.Size() != 5 {
		t.Fatal(`Union has wrong number of elements`)
	}

	x := make(map[string]byte, s.Size()+u.Size())
	x[`pi`] = 0
	x[`euler`] = 0
	x[`gauss`] = 0
	x[`lagrange`] = 0
	x[`liouville`] = 0
	for _, e := range n.Elements() {
		x[e]++
	}
	for _, k := range maphelper.Keys(x) {
		if x[k] != 1 {
			t.Fatalf(`Element '%s' is not contained exactly once`, k)
		}
	}
}

func TestFSStringEqual(t *testing.T) {
	s := NewFileSystemStringSet()
	s.Add(`pi`)
	s.Add(`euler`)
	s.Add(`gauss`)

	u := NewFileSystemStringSet()
	u.Add(`pi`)
	u.Add(`gauss`)

	v := NewFileSystemStringSet()
	v.Add(`pi`)
	v.Add(`euler`)
	v.Add(`gauss`)

	w := NewFileSystemStringSet()
	w.Add(`pi`)
	w.Add(`euler`)
	w.Add(`liouville`)

	if s.Equal(u) {
		t.Fatal(`Sets of different sizes are equal`)
	}
	if s.Equal(w) {
		t.Fatal(`Sets of equal sizes but different elements are equal`)
	}
	if !s.Equal(v) {
		t.Fatal(`Equal sets are not equal`)
	}
}
