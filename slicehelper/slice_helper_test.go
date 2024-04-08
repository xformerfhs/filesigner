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
//    2024-03-17: V1.0.0: Created.
//

package slicehelper

import (
	"math"
	"testing"
)

func TestFill(t *testing.T) {
	var s []byte
	Fill(s, 0xaa)
	if len(s) != 0 {
		t.Fatal(`Error filling empty slice.`)
	}

	u := make([]string, 7)
	Fill(u, `Empty`)
	for _, e := range u {
		if e != `Empty` {
			t.Fatal(`Error filling string slice.`)
		}
	}

	v := make([][]int, 7)
	w := make([]int, 3)
	Fill(w, -3)
	Fill(v, w)
	for _, e := range v {
		for _, f := range e {
			if f != -3 {
				t.Fatal(`Error filling slice of int slice.`)
			}
		}
	}
}

func TestFillToCap(t *testing.T) {
	u := make([]float64, 1, 3)
	FillToCap(u, math.Pi)
	x := u[:3]

	for _, e := range x {
		if e != math.Pi {
			t.Fatal(`Error filling float64 slice.`)
		}
	}
}

func TestClearNumber(t *testing.T) {
	v := make([]complex128, 7)
	Fill(v, complex(3, -7))
	ClearNumber(v)
	n := complex(0, 0)
	for _, e := range v {
		if e != n {
			t.Fatal(`Error clearing complex slice.`)
		}
	}
}

func TestConcat(t *testing.T) {
	a := make([]uint64, 7)
	Fill(a, 1)
	b := make([]uint64, 11)
	Fill(b, 99)
	c := Concat(a, b)
	for i, e := range c {
		if i < 7 {
			if e != 1 {
				t.Fatal(`Error in 1. part of concatenated slice.`)
			}
		} else {
			if e != 99 {
				t.Fatal(`Error in 2. part of concatenated slice.`)
			}
		}
	}
}

func TestPrepend(t *testing.T) {
	a := make([]uint64, 13)
	Fill(a, 11111)
	c := Prepend(a, 4747474)
	for i, e := range c {
		if i < 1 {
			if e != 4747474 {
				t.Fatal(`Error in 1. part of prepended slice.`)
			}
		} else {
			if e != 11111 {
				t.Fatal(`Error in 2. part of prepended slice.`)
			}
		}
	}
}

func TestCopy(t *testing.T) {
	a := make([]int16, 11)
	Fill(a, -7)
	b := Copy(a)
	Fill(a, 13)

	if len(a) != len(b) {
		t.Fatal(`Copy of slice has not same length as original`)
	}

	for _, e := range b {
		if e != -7 {
			t.Fatal(`Error copying slice.`)
		}
	}
}

func TestNewReverse(t *testing.T) {
	a := make([]int, 7)
	for i := 0; i < len(a); i++ {
		a[i] = 6 - i
	}
	b := NewReverse(a)

	if len(a) != len(b) {
		t.Fatal(`Reverse has not same length as original`)
	}

	for i := 0; i < len(b); i++ {
		if b[i] != i {
			t.Fatal(`Error reversing slice`)
		}
	}
}

func TestSetCap(t *testing.T) {
	a := make([]byte, 7)
	a = SetCap(a, 2)
	if cap(a) < 7 {
		t.Fatal(`Error setting smaller capacity`)
	}
	a = SetCap(a, 7)
	if cap(a) < 7 {
		t.Fatal(`Error setting equal capacity`)
	}
	a = SetCap(a, 11)
	if cap(a) < 11 {
		t.Fatal(`Error setting greater capacity`)
	}
}

func TestSetCapPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal(`No panic with negative capacity`)
		}
	}()
	a := make([]byte, 7)
	a = SetCap(a, -1)
}
