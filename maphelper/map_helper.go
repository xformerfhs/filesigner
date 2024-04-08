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

package maphelper

import (
	"cmp"
	"slices"
)

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
// This is a copy of golang.org/x/exp/maps.Keys to get rid of the
// golang.org/x/exp dependency which contains a lot of weird stuff.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// SortedKeys returns the keys of the map m. The keys will be sorted.
// slices.Sort needs cmp.Ordered.
func SortedKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	r := Keys(m)
	slices.Sort(r)

	return r
}
