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
//	http://www.apache.org/licenses/LICENSE-2.0
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
//	   2024-03-23: V1.0.0: Created.
//

package stretcher

import (
	"bytes"
	"testing"
)

func TestKeyFromBytesEmpty(t *testing.T) {
	var a []byte
	k := KeyFromBytes(a)
	e := []byte{
		0x25, 0x90, 0x48, 0xbb, 0x5a, 0x2a, 0xc8, 0x45,
		0xd2, 0x92, 0xf8, 0xd9, 0x83, 0xa0, 0x4f, 0x1d,
		0xe7, 0xd7, 0x55, 0x8e, 0x1f, 0x8a, 0xed, 0x87,
		0x04, 0x99, 0x4a, 0xae, 0xad, 0xe8, 0xad, 0xa6,
		0x00, 0x6a, 0x0f, 0x80, 0x09, 0x94, 0x5a, 0x59,
		0xe1, 0x03, 0x3d, 0x26, 0xf4, 0x20, 0xb8, 0xf0,
		0x57, 0xf9, 0x1b, 0x45, 0xf2, 0x9d, 0x85, 0x7f,
		0x53, 0x83, 0x02, 0x14, 0x68, 0xcc, 0x4e, 0xa1,
		0x2a}
	if bytes.Compare(k, e) != 0 {
		t.Fatal(`Wrong result from empty context id`)
	}
}

func TestKeyFromBytes(t *testing.T) {
	k := KeyFromBytes([]byte(`WärmeØlGóðaNótt`))
	e := []byte{
		0xce, 0x8f, 0xca, 0xde, 0xc5, 0x9f, 0x65, 0xbb,
		0xcb, 0x66, 0x57, 0x8f, 0xa8, 0xee, 0xde, 0x25,
		0xa3, 0x77, 0xea, 0xcc, 0x90, 0xfc, 0xb6, 0x52,
		0xf3, 0x27, 0xa7, 0xd2, 0xc3, 0x2e, 0xcf, 0x4b,
		0x57, 0xc3, 0xa4, 0x72, 0x6d, 0x65, 0xc3, 0x98,
		0x6c, 0x47, 0xc3, 0xb3, 0xc3, 0xb0, 0x61, 0x4e,
		0xc3, 0xb3, 0x74, 0x74, 0x14, 0x11, 0xc2, 0x0f,
		0x48, 0x9a, 0x38, 0xc8, 0x91, 0x67, 0xdd, 0x59,
		0x67, 0x96, 0x71, 0xe7, 0x83, 0x04, 0xd3, 0x75,
		0xde, 0xd5, 0x50, 0xef, 0x97, 0x84, 0x19, 0x78,
		0x03, 0x1f, 0x7c, 0x49, 0x76,
	}
	if bytes.Compare(k, e) != 0 {
		t.Fatal(`Wrong result from normal context id`)
	}
}
