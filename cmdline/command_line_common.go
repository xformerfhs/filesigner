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

package cmdline

import (
	"errors"
	"filesigner/filehelper"
	"strings"
)

// ******** Private constants ********

// defaultSignaturesFileNamePrefix is the default prefix name part of the signatures file.
const defaultSignaturesFileNamePrefix = `filesigner`

// signaturesFileNameSuffix is the suffix of the signatures file name.
const signaturesFileNameSuffix = `-signatures.json`

// wildCards contains the valid wild card characters.
const wildCards = `*?`

// ******** Private functions ********

// checkSignaturesFileName checks if the supplied file path is only a file name.
func checkSignaturesFileName(filePath string) error {
	if !filehelper.IsFileName(filePath) {
		return errors.New(`Prefix must not contain path separators`)
	}
	if strings.ContainsAny(filePath, wildCards) {
		return errors.New(`Prefix must not contain wild cards`)
	}

	return nil
}
