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

// Package base32encoding implements the base32 encodings used in this program.
package base32encoding

import (
	"encoding/base32"
)

// ******** Private constants ********

const keySeparator = '-'
const keyGroupSize = 4

// ******** Private variables ********

// enc is a base32 encoder that uses the word-safe alphabet.
var enc = base32.NewEncoding(`3479BCDFGHJLMRQSTVZbcdfghjmrstvz`).WithPadding(base32.NoPadding)

// encKey is a base32 encoder which uses a custom alphabet with only upper-case letters and digits to encode keys.
var encKey = base32.NewEncoding(`B9C8D7E6F5G4H3J2K1L0MNPQRSTVWXYZ`).WithPadding(base32.NoPadding)
